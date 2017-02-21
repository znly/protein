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

	_ "github.com/znly/protein/protobuf/test"
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
	tsKnownName = ".test.TestSchema"
	deKnownName = ".test.TestSchema.DepsEntry"

	tsKnownHashSingle = "PROT-d12b5690ec5877c9c5f6c84b57523b14291b46de50ffaf487a8a3f04d8198bff"
	deKnownHashSingle = "PROT-39b441710129f5f7e4c92863b02173db7aa71658799031d6a1dedca260842dfa"

	tsKnownHashRecurse = "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	deKnownHashRecurse = "PROT-d278f5561f05e68f6e68fcbc6b801d29a69b4bf6044bf3e6242ea8fe388ebd6e"
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

	// should at least find 2 messages types here: `.test.TestSchema`
	// and its nested `DepsEntry` message
	assert.True(t, len(dtsByName) >= 2)
	assert.NotNil(t, dtsByName[tsKnownName])
	assert.NotNil(t, dtsByName[deKnownName])

	return dtsByName
}

func TestProtoscan_collectDescriptorTypes(t *testing.T) {
	dtsByName := _collectTestSchemaTrees(t)

	psDT := dtsByName[tsKnownName]
	assert.Nil(t, psDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, tsKnownName, psDT.fqName)
	assert.NotNil(t, psDT.descr)
	assert.Equal(t, "TestSchema", psDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err := proto.Marshal(psDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	psExpectedHash, err := ByteSSlice{[]byte(psDT.FQName()), b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, tsKnownHashSingle, "PROT-"+hex.EncodeToString(psExpectedHash))
	assert.Equal(t, psExpectedHash, psDT.hashSingle)
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[deKnownName]
	assert.Nil(t, deDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, deKnownName, deDT.fqName)
	assert.NotNil(t, deDT.descr)
	assert.Equal(t, "DepsEntry", deDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err = proto.Marshal(deDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	deExpectedHash, err := ByteSSlice{[]byte(deDT.FQName()), b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, deKnownHashSingle, "PROT-"+hex.EncodeToString(deExpectedHash))
	assert.Equal(t, deExpectedHash, deDT.hashSingle)
	assert.Nil(t, deDT.hashRecursive)
}

// -----------------------------------------------------------------------------

func TestProtoscan_DescriptorTree_computeDependencyLinks(t *testing.T) {
	dtsByName := _collectTestSchemaTrees(t)
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeDependencyLinks(dtsByName))
	}

	psDT := dtsByName[tsKnownName]
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should only find a dependency to `DepsEntry` in here
	assert.Len(t, depsMap, 1)
	assert.Contains(t, depsMap, deKnownName)
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[deKnownName]
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

	psDT := dtsByName[tsKnownName]
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should only find a dependency to `DepsEntry` in here
	assert.Len(t, depsMap, 1)
	assert.Contains(t, depsMap, deKnownName)
	// since `TestSchema` has no dependency, its recursive hash should just
	// be a re-hash of its single hash
	psExpectedHash, err := ByteSSlice{
		psDT.hashSingle,
		depsMap[deKnownName].hashSingle,
	}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, tsKnownHashRecurse, "PROT-"+hex.EncodeToString(psExpectedHash))
	assert.Equal(t, psDT.hashRecursive, psExpectedHash)

	deDT := dtsByName[deKnownName]
	assert.Empty(t, deDT.deps) // DepsEntry has no dependency
	// since `DepsEntry` has no dependency, its recursive hash should just
	// be a re-hash of its single hash
	deExpectedHash, err := ByteSSlice{deDT.hashSingle}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, deKnownHashRecurse, "PROT-"+hex.EncodeToString(deExpectedHash))
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

	// should at least find the `.test.TestSchema` and its nested
	// `DepsEntry` message in the DescriptorTrees
	assert.NotEmpty(t, dtsByUID)
	assert.NotNil(t, dtsByUID[tsKnownHashRecurse])
	assert.Equal(t, tsKnownName,
		dtsByUID[tsKnownHashRecurse].FQName(),
	)
	assert.Equal(t, deKnownName,
		dtsByUID[deKnownHashRecurse].FQName(),
	)
}
