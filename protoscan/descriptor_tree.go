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

import "github.com/gogo/protobuf/proto"

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
	// package, and is traditionally used a key inside a schema registry.
	hashRecursive string
}

// Descr returns the Message/Enum descriptor.
func (d Descriptor) Descr() proto.Message { return d.descr }

// UID returns the unique, deterministic, versioned identifier of a Descriptor.
func (d Descriptor) UID() string { return d.hashRecursive }
