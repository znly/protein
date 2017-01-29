// Copyright © 2016 Zenly <hello@zen.ly>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protoscan

import (
	"encoding/hex"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

type descriptorNode struct {
	// descr is either a DescriptorProto (Message) or an
	// EnumDescriptorProto (Enum).
	descr proto.Message
	// fqName is the fully-qualified name of the descriptor `descr`
	// (e.g. .google.protobuf.Timestamp).
	fqName string

	// hashSingle is the the hash (as implemented in byte_sslice.go) of
	// `this.descr`.
	// hashSingle = SHA1(d.descr)
	//
	// It is used to uniquely and deterministically identify a Message or Enum
	// type that has been instanciated by a protobuf library.
	//
	// hashSingle is necessary for the internal DescriptorTree machinery to work,
	// but it is never exposed to the end-user of the protoscan package.
	hashSingle []byte
	// hashRecursive is the the hash (as implemented in byte_sslice.go) of
	// `this.descr` + the hashes of every type it depends on (recursively).
	// Note that, to ensure determinism, the list of hashes used to compute
	// the final `hashRecursive` value is alphabetically-sorted first.
	// hashRecursive = SHA1(SORT(d.hashSingle, d.AllDependencyHashes))
	//
	// It is the unique, deterministic and versioned identifier for a Message or
	// Enum type.
	// This is the actual value that is exposed to the end-user of the protoscan
	// package, and is traditionally used as the key inside a schema registry.
	hashRecursive []byte
}

// -----------------------------------------------------------------------------

// A DescriptorTree is a dependency tree of Message/Enum descriptors.
type DescriptorTree struct {
	*descriptorNode
	deps []*DescriptorTree
}

// FQName returns the fully-qualified name of the underlying descriptor `descr`
// (e.g. .google.protobuf.Timestamp).
func (dt *DescriptorTree) FQName() string { return dt.fqName }

// UID returns a unique, deterministic, versioned identifier for this particular
// DescriptorTree.
//
// This identifier is computed from `dt`'s protobuf schema as well as its
// dependencies' schemas.
// This means that modifying any schema in the dependency tree, not just `dt`'s,
// will result in a new identifier being generated.
//
// The returned string is the hexadecimal representation of `dt`'s internal
// recursive hash.
func (dt *DescriptorTree) UID() string {
	return hex.EncodeToString(dt.hashRecursive)
}

// DependencyUIDs recursively walks through the dependencies of `dt` and
// returns a sorted list of the hexadecimal representations of their
// respective recursive hashes.
//
// This should only be called once the linking & recursive hashing computation
// have been done.
func (dt *DescriptorTree) DependencyUIDs() []string {
	// used to avoid non-linear dependencies
	alreadyMet := make(map[*DescriptorTree]struct{}, len(dt.deps))
	var recurseChildren func(dts []*DescriptorTree) []string
	recurseChildren = func(dts []*DescriptorTree) []string {
		recursiveHashes := make([]string, 0, len(dts))
		for _, dt := range dts {
			if _, ok := alreadyMet[dt]; ok {
				continue
			}
			alreadyMet[dt] = struct{}{}
			recursiveHashes = append(recursiveHashes, dt.UID())
			recursiveHashes = append(
				recursiveHashes, recurseChildren(dt.deps)...,
			)
		}
		return recursiveHashes
	}

	return recurseChildren(dt.deps)
}

// -----------------------------------------------------------------------------

func NewDescriptorTrees(
	fdps map[string]*descriptor.FileDescriptorProto,
) (map[string]*DescriptorTree, error) {
	// all of the Message/Enum/Nested types currently instanciated,
	// arranged by their fully-qualified names
	dtsByName, err := collectDescriptorTypes(fdps)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Computes the immediated dependencies of every available DescriptorTree.
	//
	// Note that we need all the dependency trees of all DescriptorTrees
	// to be already computed before we can compute the definitive recursive
	// hashes.
	for _, dt := range dtsByName {
		if err := dt.computeDependencyLinks(dtsByName); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// computes the recursive hash value for each available DescriptorTree
	for _, dt := range dtsByName {
		if err := dt.computeRecursiveHash(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// builds the final map of DescriptorTrees, arranged by UIDs
	// (i.e. recursive hashes)
	dtsByUID := make(map[string]*DescriptorTree, len(dtsByName))
	for _, dt := range dtsByName {
		dtsByUID[dt.UID()] = dt
	}

	return dtsByUID, nil
}

// newDescriptorTree returns a new DescriptorTree with its `hashSingle`
// field already computed.
//
// `fqName` is necessary to defuse single-hash collision issues.
// Have a look at the implementation for details.
//
// At this stage, dependencies have not been calculated yet, hence the
// `deps` list of the DescriptorTree is nil.
func newDescriptorTree(
	fqName string, descr proto.Message,
) (*DescriptorTree, error) {
	if err := checkDescriptorType(descr); err != nil {
		return nil, errors.WithStack(err)
	}

	descrBytes, err := proto.Marshal(descr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Two purely-identical Message or Enum types will share the same exact
	// binary protofile representation, even if they have different names.
	// That's why we must make sure that the fully-qualified name of the type
	// is part of the input of the hashing function.
	//
	// E.g., in the protein package, both `.protoscan.ProtobufSchema.DepsEntry`
	// and `.protoscan.TestSchema.DepsEntry` have the same binary encoding, but
	// different names.
	hashSingle, err := ByteSSlice{[]byte(fqName), descrBytes}.Hash()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dt := &DescriptorTree{
		descriptorNode: &descriptorNode{
			descr:      descr,
			fqName:     fqName,
			hashSingle: hashSingle,
		},
	}

	return dt, nil
}

// -----------------------------------------------------------------------------

// computeDependencyLinks adds the immediate dependencies of `dt` to its
// `deps` list.
//
// At this stage, the newly-added dependencies are not guaranteed to have
// their own dependency-links already computed (i.e. their respective
// `deps` lists are still nil), since their `computeDependencyLinks` methods
// might not have been called yet (the order in which each DescriptorTree
// gets its links computed is non-deterministic).
// For that reason, recursive hashes cannot yet be computed from here.
//
// NOTE: This must be called before calling `computeRecursiveHash()`.
func (dt *DescriptorTree) computeDependencyLinks(
	dtsByName map[string]*DescriptorTree,
) error {
	// Multiple fields of the current descriptor might rely on the same
	// Message or Enum type: we need to count those recurring dependencies
	// as one, or all hell will break loose later on in the pipeline.
	alreadyMet := make(map[*DescriptorTree]struct{}, len(dt.deps))
	switch descr := dt.descr.(type) {
	case *descriptor.DescriptorProto:
		for _, f := range descr.GetField() {
			// `typeName` is non-empty only if it references a Message or
			// Enum type
			if typeName := f.GetTypeName(); len(typeName) > 0 {
				dep, ok := dtsByName[typeName]
				if !ok {
					return errors.Errorf("`%s`: no such type", typeName)
				}
				if _, ok := alreadyMet[dep]; ok {
					continue
				}
				alreadyMet[dep] = struct{}{}
				dt.deps = append(dt.deps, dep)
			}
		}
	case *descriptor.EnumDescriptorProto:
		// nothing to do
	default:
		return errors.Errorf("`%v`: illegal type", reflect.TypeOf(descr))
	}
	return nil
}

// computeRecursiveHash recursively walks through `dt`'s dependencies
// in order to compute its recursive hash.
//
// NOTE: `computeDependencyLinks()` must have been called first for every
// available DescriptorTree.
func (dt *DescriptorTree) computeRecursiveHash() error {
	// used to avoid non-linear dependencies
	alreadyMet := make(map[*DescriptorTree]struct{}, len(dt.deps))
	var recurseChildren func(dts []*DescriptorTree) (ByteSSlice, error)
	recurseChildren = func(dts []*DescriptorTree) (ByteSSlice, error) {
		singleHashes := make(ByteSSlice, 0, len(dts))
		for _, dt := range dts {
			if _, ok := alreadyMet[dt]; ok {
				continue
			}
			alreadyMet[dt] = struct{}{}
			singleHashes = append(singleHashes, dt.hashSingle)
			singleHashesRec, err := recurseChildren(dt.deps)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			singleHashes = append(singleHashes, singleHashesRec...)
		}
		return singleHashes, nil
	}

	singleHashes, err := recurseChildren(dt.deps)
	if err != nil {
		return errors.WithStack(err)
	}

	// /!\ Do not forget to add `dt`'s own single hash!
	// Otherwise, any DescriptorTree with the exact same set of dependencies
	// would end up with the same recursive hash, thus overwriting each other.
	singleHashes = append(singleHashes, dt.hashSingle)

	hashRecursive, err := singleHashes.Hash()
	if err != nil {
		return errors.WithStack(err)
	}
	dt.hashRecursive = hashRecursive

	return nil
}

// -----------------------------------------------------------------------------

// collectDescriptorTypes recursively goes through all the FileDescriptorProtos
// that are passed to it in order to collect all of the Message, Enum & Nested
// types that can be found in these descriptors.
//
// The result is a flattened map of single-level (i.e. `deps` is nil)
// DescriptorTrees that are arranged by the types' respective fully-qualified
// names (e.g. `.google.protobuf.Timestamp`).
//
// This map is a necessary intermediate representation for computing the final
// dependency trees.
func collectDescriptorTypes(
	fdps map[string]*descriptor.FileDescriptorProto,
) (map[string]*DescriptorTree, error) {
	// Pre-allocating len(fdps)*5, assuming an average of 5 Message/Enum
	// types per file.
	// This mapping is arranged by the descriptors' respective FQ-names and is
	// a necessary intermediate representation for the final computation of
	// dependency trees.
	dtsByName := make(map[string]*DescriptorTree, len(fdps)*5)

	// recurse through a protoFile and pushes every Message/Enum/Nested
	// type to the `dtsByName` map
	var recurseMessages func(mt *descriptor.DescriptorProto, parent string) error
	recurseMessages = func(mt *descriptor.DescriptorProto, parent string) error {
		parent += ("." + mt.GetName())
		dt, err := newDescriptorTree(parent, mt)
		if err != nil {
			return errors.WithStack(err)
		}
		dtsByName[dt.fqName] = dt
		for _, et := range mt.GetEnumType() {
			dt, err := newDescriptorTree(parent+"."+et.GetName(), et)
			if err != nil {
				return errors.WithStack(err)
			}
			dtsByName[dt.fqName] = dt
		}
		for _, nmt := range mt.GetNestedType() {
			if err := recurseMessages(nmt, parent); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}

	// walk through the protoFiles
	for _, fdp := range fdps {
		parent := "." + fdp.GetPackage()
		for _, et := range fdp.GetEnumType() {
			dt, err := newDescriptorTree(parent+"."+et.GetName(), et)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			dtsByName[dt.fqName] = dt
		}
		for _, mt := range fdp.GetMessageType() {
			if err := recurseMessages(mt, parent); err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}

	return dtsByName, nil
}

// -----------------------------------------------------------------------------

func checkDescriptorType(descr proto.Message) error {
	switch descr.(type) {
	case *descriptor.DescriptorProto:
	case *descriptor.EnumDescriptorProto:
	default:
		return errors.Errorf("`%v`: illegal type", reflect.TypeOf(descr))
	}
	return nil
}
