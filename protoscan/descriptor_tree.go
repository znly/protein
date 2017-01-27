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
	"fmt"
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
	children []*DescriptorTree
}

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
	return fmt.Sprintf("%x", dt.hashRecursive)
}

// -----------------------------------------------------------------------------

func NewDescriptorTrees(
	fdps map[string]*descriptor.FileDescriptorProto,
) (map[string]*DescriptorTree, error) {
	/* DEBUG */
	for file, fdp := range fdps {
		fmt.Println("file:", file)
		for _, dep := range fdp.GetDependency() {
			fmt.Println("\tdependency:", dep)
		}
		for _, et := range fdp.GetEnumType() {
			fmt.Println("\tenum type:", et.GetName())
		}
		for _, mt := range fdp.GetMessageType() {
			fmt.Println("\tmessage type:", mt.GetName())
		}
	}

	// all of the Message/Enum/Nested types currently instanciated,
	// arranged by their fully-qualified names
	dtsByName, err := collectDescriptorTypes(fdps)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Computes the `children` lists for every available DescriptorTree.
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

	return nil, nil
}

// newDescriptorTree returns a new DescriptorTree with its `hashSingle`
// field already computed.
// At this stage, dependencies have not been calculated yet, hence the
// `children` slice of the DescriptorTree is nil.
func newDescriptorTree(descr proto.Message) (*DescriptorTree, error) {
	dt := &DescriptorTree{descriptorNode: &descriptorNode{descr: descr}}

	if err := checkDescriptorType(descr); err != nil {
		return nil, errors.WithStack(err)
	}

	descrBytes, err := proto.Marshal(descr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	hashSingle, err := ByteSSlice{descrBytes}.Hash()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dt.hashSingle = hashSingle

	return dt, nil
}

// -----------------------------------------------------------------------------

// computeDependencyLinks appends the DependencyTrees of `dt`'s dependencies
// to its list of children.
//
// At this stage, the newly-appended dependencies are not guaranteed to have
// their own dependency-links already computed (i.e. their respective
// `children` lists are still nil), since their `computeDependencyLinks`
// might not have been called yet (the order in which each DescriptorTree
// gets its links computed is random).
// For that reason, recursive hashes cannot yet be computed here.
//
// This must be called before calling `computeRecursiveHash()`.
func (dt *DescriptorTree) computeDependencyLinks(
	dtsByName map[string]*DescriptorTree,
) error {
	switch descr := dt.descr.(type) {
	case *descriptor.DescriptorProto:
		for _, f := range descr.GetField() {
			// `typeName` is non-empty only if it references a Message or
			// Enum type
			if typeName := f.GetTypeName(); len(typeName) > 0 {
				child, ok := dtsByName[typeName]
				if !ok {
					return errors.Errorf("`%s`: no such type")
				}
				dt.children = append(dt.children, child)
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
// `computeDependencyLinks()` must be called first for every available
// DescriptorTree.
func (dt *DescriptorTree) computeRecursiveHash() error {
	// used to avoid cyclic dependencies
	alreadyMet := make(map[*DescriptorTree]struct{}, len(dt.children))

	var recurseChildren func(dts []*DescriptorTree) (ByteSSlice, error)
	recurseChildren = func(dts []*DescriptorTree) (ByteSSlice, error) {
		singleHashes := make(ByteSSlice, 0, len(dts))
		for _, dt := range dts {
			if _, ok := alreadyMet[dt]; ok {
				continue
			}
			alreadyMet[dt] = struct{}{}
			singleHashesRec, err := recurseChildren(dt.children)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			singleHashes = append(singleHashes, dt.hashSingle)
			singleHashes = append(singleHashes, singleHashesRec...)
		}
		return singleHashes, nil
	}

	singleHashes, err := recurseChildren(dt.children)
	if err != nil {
		return errors.WithStack(err)
	}
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
// The result is a flattened map of single-level (i.e. `children` is nil)
// DescriptorTrees arranged by the types' absolute names
// (e.g. `.google.protobuf.Timestamp`).
//
// This map is a necessary intermediate representation for computing the final
// dependency trees.
func collectDescriptorTypes(
	fdps map[string]*descriptor.FileDescriptorProto,
) (map[string]*DescriptorTree, error) {
	// Pre-allocating len(fdps)*3, assuming an average of 3 Message/Enum
	// types per file.
	// This mapping is arranged by the descriptors' respective name and is
	// a necessary intermediate representation for the final computation of
	// dependency trees.
	dtsByName := make(map[string]*DescriptorTree, len(fdps)*3)

	// recurse through a protoFile and pushes every Message/Enum/Nested
	// type to the `dtsByName` map
	var recurseMessages func(mt *descriptor.DescriptorProto, parent string) error
	recurseMessages = func(mt *descriptor.DescriptorProto, parent string) error {
		dt, err := newDescriptorTree(mt)
		if err != nil {
			return errors.WithStack(err)
		}
		parent += ("." + mt.GetName())
		dtsByName[parent] = dt
		for _, et := range mt.GetEnumType() {
			dt, err := newDescriptorTree(et)
			if err != nil {
				return errors.WithStack(err)
			}
			dtsByName[parent+"."+et.GetName()] = dt
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
			dt, err := newDescriptorTree(et)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			dtsByName[parent+"."+et.GetName()] = dt
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
