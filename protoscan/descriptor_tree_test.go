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
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

// `TestSchema` and `DepsEntry` are both immutable, for-testing-purposes
// only schemas and, as such, both their respective single & recursive
// hashes can be known in advance and shouldn't ever change unless the
// `ByteSSlice.Hash` method is ever modified.
//
// If such a modification of `ByteSSlice.Hash`'s behavior were to happen,
// you'd have to modify the following expected values in order to fix the
// tests.. That is, if you're sure about what you're doing.
var (
	psKnownHashSingle = "a9fbbb51c6d1e9aabee6465bba6331cdfad337f420a18d0e3e6f4d8e2cb6035a"
	deKnownHashSingle = "b0c930d41e1f44abe96a624a1dc9b7b0473c1f5adb8270abc3c67f01244329af"

	psKnownHashRecurse = "bb06a3864164fb5af718c6a2f4f0b70032be9e04a2a756ef45d530d4b6e2a570"
	deKnownHashRecurse = "d1f1090b178845e096dc0569b7ad93bd0eb711b1fa2d2172a53a7940d5d9024f"
)

// -----------------------------------------------------------------------------

func _collectTestSchemaTrees(t *testing.T) map[string]*DescriptorTree {
	// retrieve instanciated gogo/protobuf protofiles
	symbol := "github.com/gogo/protobuf/proto.protoFiles"
	protoFilesBindings, err := BindProtofileSymbols()
	assert.Nil(t, err)
	assert.NotEmpty(t, protoFilesBindings)
	assert.NotEmpty(t, protoFilesBindings[symbol])
	protoFiles := *protoFilesBindings[symbol]
	assert.NotEmpty(t, protoFiles)

	fdps := map[string]*descriptor.FileDescriptorProto{}
	for file, descr := range protoFiles {
		fdp, err := UnzipAndUnmarshal(descr)
		assert.Nil(t, err)
		fdps[file] = fdp
	}

	dtsByName, err := collectDescriptorTypes(fdps)
	assert.Nil(t, err)
	assert.NotEmpty(t, dtsByName)

	// should at least find 2 messages types here: `.protoscan.TestSchema`
	// and its nested `DepsEntry` message
	assert.True(t, len(dtsByName) >= 2)
	assert.NotNil(t, dtsByName[".protoscan.TestSchema"])
	assert.NotNil(t, dtsByName[".protoscan.TestSchema.DepsEntry"])

	return dtsByName
}

func TestProtoscan_collectDescriptorTypes(t *testing.T) {
	dtsByName := _collectTestSchemaTrees(t)

	psDT := dtsByName[".protoscan.TestSchema"]
	assert.Nil(t, psDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, ".protoscan.TestSchema", psDT.fqName)
	assert.NotNil(t, psDT.descr)
	assert.Equal(t, "TestSchema", psDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err := proto.Marshal(psDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	psExpectedHash, err := ByteSSlice{b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, psKnownHashSingle, hex.EncodeToString(psExpectedHash))
	assert.Equal(t, psExpectedHash, psDT.hashSingle)
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[".protoscan.TestSchema.DepsEntry"]
	assert.Nil(t, deDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, ".protoscan.TestSchema.DepsEntry", deDT.fqName)
	assert.NotNil(t, deDT.descr)
	assert.Equal(t, "DepsEntry", deDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err = proto.Marshal(deDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	deExpectedHash, err := ByteSSlice{b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, deKnownHashSingle, hex.EncodeToString(deExpectedHash))
	assert.Equal(t, deExpectedHash, deDT.hashSingle)
	assert.Nil(t, deDT.hashRecursive)
}

// -----------------------------------------------------------------------------

func TestProtoscan_DescriptorTree_computeDependencyLinks(t *testing.T) {
	dtsByName := _collectTestSchemaTrees(t)
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeDependencyLinks(dtsByName))
	}

	psDT := dtsByName[".protoscan.TestSchema"]
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should only find a dependency to `DepsEntry` in here
	assert.Len(t, depsMap, 1)
	assert.Contains(t, depsMap, ".protoscan.TestSchema.DepsEntry")
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[".protoscan.TestSchema.DepsEntry"]
	assert.Empty(t, deDT.deps) // DepsEntry has no dependency
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, deDT.hashRecursive)
}

func TestProtoscan_DescriptorTree_computeRecursiveHash(t *testing.T) {
	dtsByName := _collectTestSchemaTrees(t)
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeDependencyLinks(dtsByName))
	}
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeRecursiveHash())
	}

	psDT := dtsByName[".protoscan.TestSchema"]
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should only find a dependency to `DepsEntry` in here
	assert.Len(t, depsMap, 1)
	assert.Contains(t, depsMap, ".protoscan.TestSchema.DepsEntry")
	// since `TestSchema` has no dependency, its recursive hash should just
	// be a re-hash of its single hash
	psExpectedHash, err := ByteSSlice{
		psDT.hashSingle,
		depsMap[".protoscan.TestSchema.DepsEntry"].hashSingle,
	}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, psKnownHashRecurse, hex.EncodeToString(psExpectedHash))
	assert.Equal(t, psDT.hashRecursive, psExpectedHash)

	deDT := dtsByName[".protoscan.TestSchema.DepsEntry"]
	assert.Empty(t, deDT.deps) // DepsEntry has no dependency
	// since `DepsEntry` has no dependency, its recursive hash should just
	// be a re-hash of its single hash
	deExpectedHash, err := ByteSSlice{deDT.hashSingle}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, deKnownHashRecurse, hex.EncodeToString(deExpectedHash))
	assert.Equal(t, deDT.hashRecursive, deExpectedHash)
}
