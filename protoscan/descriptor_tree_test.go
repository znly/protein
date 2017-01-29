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
	psKnownHashSingle = "ffbe6433ce68e1a2ef3947cb52471b44479bc7bc63631d854cd57389e776a5e3"
	deKnownHashSingle = "f691d70fe2b9e740f259fb7a634762e501f46c05a2ea84214f29f382395b4a7e"

	psKnownHashRecurse = "103b218c29e5ea98fc2d9a46d964b468cffa76c76dadfdf31e9f8552273af0a7"
	deKnownHashRecurse = "51fb869733d7b9926a2365cb8306125709cf53a05e97a8de00e98709cde91bab"
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
	psExpectedHash, err := ByteSSlice{[]byte(psDT.FQName()), b}.Hash()
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
	deExpectedHash, err := ByteSSlice{[]byte(deDT.FQName()), b}.Hash()
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

// -----------------------------------------------------------------------------

func TestProtoscan_NewDescriptorTrees(t *testing.T) {
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

	dtsByUID, err := NewDescriptorTrees(fdps)
	assert.Nil(t, err)

	// should at least find the `.protoscan.TestSchema` and its nested
	// `DepsEntry` message in the DescriptorTrees
	assert.NotEmpty(t, dtsByUID)
	assert.NotNil(t, dtsByUID[psKnownHashRecurse])
	assert.Equal(t, ".protoscan.TestSchema",
		dtsByUID[psKnownHashRecurse].FQName(),
	)
	assert.Equal(t, ".protoscan.TestSchema.DepsEntry",
		dtsByUID[deKnownHashRecurse].FQName(),
	)
}
