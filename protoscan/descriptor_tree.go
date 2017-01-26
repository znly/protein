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
	"crypto/sha1"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

type Descriptor struct {
	// descr is either a DescriptorProto (Message) or an
	// EnumDescriptorProto (Enum).
	descr proto.Message

	// hashSingle is the the SHA-1 hash of `this.descr`.
	// hashSingle = SHA1(d.descr)
	//
	// It is used to uniquely and deterministically identify a Message or Enum
	// type that has been instanciated by a protobuf library.
	// NOTE: You cannot rely on a type's GetName() method; you will end up with
	//       discrepancies between the name returned by this method and the
	//       name used by other types which import this type as a dependency.
	//       Really, you don't wanna go down that road.
	//
	// hashSingle is necessary for the internal DescriptorTree machinery to,
	// but it is never exposed to the end-user of the protoscan package.
	hashSingle string
	// hashRecursive is the the SHA-1 hash of `this.descr` + the SHA-1 hashes
	// of every type it depends of.
	// To be deterministic, this list of hashes is sorted alphabetically first.
	// hashRecursive = SHA1(SORT(d.hashSingle, d.AllDependencyHashes))
	//
	// It is the unique, deterministic and versioned identifier for a Message or
	// Enum type.
	// This is the actual value that is exposed to the end-user of the protoscan
	// package, and is traditionally used as key inside a schema registry.
	hashRecursive string
}

// Descr returns the Message/Enum descriptor.
func (d Descriptor) Descr() proto.Message { return d.descr }

// UID returns the unique, deterministic, versioned identifier of a Descriptor.
func (d Descriptor) UID() string { return d.hashRecursive }

// -----------------------------------------------------------------------------

// A DescriptorTree is a dependency tree of Message/Enum descriptors.
type DescriptorTree struct {
	*Descriptor
	children []*DescriptorTree
}

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

	// Pre-allocating len(fdps)*2, assuming an average of 3 Message/Enum
	// types per file.
	dtsSingle := make(map[string]*DescriptorTree, len(fdps)*3)
	// This intermediate mapping is arranged by the descriptors' `hashSingle`
	// and is used for the final computation of dependency trees.
	for _, fdp := range fdps {
		for _, et := range fdp.GetEnumType() {
			dt, err := NewDescriptorTree(et)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			dtsSingle[dt.Descriptor.hashSingle] = dt
		}
		for _, mt := range fdp.GetMessageType() {
			dt, err := NewDescriptorTree(mt)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			dtsSingle[dt.Descriptor.hashSingle] = dt
		}
	}

	return nil, nil
}

// NewDescriptorTree returns a new DescriptorTree with its `hashSingle`
// field already computed.
// At this stage, dependencies have not been calculated, hence the `children`
// slice of the Descriptor is nil.
func NewDescriptorTree(descr proto.Message) (*DescriptorTree, error) {
	dt := &DescriptorTree{Descriptor: &Descriptor{descr: descr}}

	switch descr.(type) {
	case *descriptor.DescriptorProto:
	case *descriptor.EnumDescriptorProto:
	default:
		return nil, errors.Errorf("")
	}

	descrBytes, err := proto.Marshal(descr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	h := sha1.New()
	_, err = h.Write(descrBytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dt.Descriptor.hashSingle = fmt.Sprintf("%x", h.Sum(nil))

	return dt, nil
}

func (dt *DescriptorTree) ComputeDependencies(
	dtsSingle map[string]*DescriptorTree,
) {

}
