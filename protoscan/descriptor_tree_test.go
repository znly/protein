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

func _collectTestSchemaTrees(t *testing.T) map[string]*DescriptorTree {
	// retrieve instanciated gogo/protobuf protofiles
	protoFilesBindings, err := BindProtofileSymbols()
	assert.Nil(t, err)
	assert.NotEmpty(t, protoFilesBindings)
	symbol := "github.com/gogo/protobuf/proto.protoFiles"
	if protoFilesBindings[symbol] == nil { // travis-ci runs tests through vendoring
		symbol = "github.com/znly/protein/vendor/" + symbol
	}
	assert.NotEmpty(t, protoFilesBindings[symbol])
	protoFiles := *protoFilesBindings[symbol]
	assert.NotEmpty(t, protoFiles)

	fdps := map[string]*descriptor.FileDescriptorProto{}
	for file, descr := range protoFiles {
		fdp, err := UnzipAndUnmarshal(descr)
		assert.Nil(t, err)
		fdps[file] = fdp
	}

	dtsByName, err := collectDescriptorTypes(MD5, fdps)
	assert.Nil(t, err)
	assert.NotEmpty(t, dtsByName)

	// should at least find 2 messages types here: `.test.TestSchema`
	// and its nested `DepsEntry` message
	assert.True(t, len(dtsByName) >= 2)
	assert.NotNil(t, dtsByName[TEST_TSKnownName])
	assert.NotNil(t, dtsByName[TEST_DEKnownName])

	return dtsByName
}

func TestProtoscan_collectDescriptorTypes(t *testing.T) {
	var bss ByteSSlice
	dtsByName := _collectTestSchemaTrees(t)

	psDT := dtsByName[TEST_TSKnownName]
	assert.Nil(t, psDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, TEST_TSKnownName, psDT.fqName)
	assert.NotNil(t, psDT.descr)
	assert.Equal(t, "TestSchema", psDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err := proto.Marshal(psDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	bss = ByteSSlice{[]byte(psDT.FQName()), b}
	bss.Sort()
	psExpectedHash, err := MD5(bss)
	assert.Nil(t, err)
	assert.Equal(t, TEST_TSKnownHashSingle, "PROT-"+hex.EncodeToString(psExpectedHash))
	assert.Equal(t, psExpectedHash, psDT.hashSingle)
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[TEST_DEKnownName]
	assert.Nil(t, deDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, TEST_DEKnownName, deDT.fqName)
	assert.NotNil(t, deDT.descr)
	assert.Equal(t, "DepsEntry", deDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err = proto.Marshal(deDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	bss = ByteSSlice{[]byte(deDT.FQName()), b}
	bss.Sort()
	deExpectedHash, err := MD5(bss)
	assert.Nil(t, err)
	assert.Equal(t, TEST_DEKnownHashSingle, "PROT-"+hex.EncodeToString(deExpectedHash))
	assert.Equal(t, deExpectedHash, deDT.hashSingle)
	assert.Nil(t, deDT.hashRecursive)
}

// -----------------------------------------------------------------------------

func TestDescriptorTree_computeDependencyLinks(t *testing.T) {
	dtsByName := _collectTestSchemaTrees(t)
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeDependencyLinks(dtsByName))
	}

	psDT := dtsByName[TEST_TSKnownName]
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should only find a dependency to `DepsEntry` in here
	assert.Len(t, depsMap, 1)
	assert.Contains(t, depsMap, TEST_DEKnownName)
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[TEST_DEKnownName]
	assert.Empty(t, deDT.deps) // DepsEntry has no dependency
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, deDT.hashRecursive)
}

func TestDescriptorTree_computeRecursiveHash(t *testing.T) {
	var bss ByteSSlice
	dtsByName := _collectTestSchemaTrees(t)
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeDependencyLinks(dtsByName))
	}
	for _, dt := range dtsByName {
		assert.Nil(t, dt.computeRecursiveHash(MD5))
	}

	psDT := dtsByName[TEST_TSKnownName]
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should only find a dependency to `DepsEntry` in here
	assert.Len(t, depsMap, 1)
	assert.Contains(t, depsMap, TEST_DEKnownName)
	// since `TestSchema` has no dependency, its recursive hash should just
	// be a re-hash of its single hash
	bss = ByteSSlice{psDT.hashSingle, depsMap[TEST_DEKnownName].hashSingle}
	bss.Sort()
	psExpectedHash, err := MD5(bss)
	assert.Nil(t, err)
	assert.Equal(t, TEST_TSKnownHashRecurse, "PROT-"+hex.EncodeToString(psExpectedHash))
	assert.Equal(t, psDT.hashRecursive, psExpectedHash)

	deDT := dtsByName[TEST_DEKnownName]
	assert.Empty(t, deDT.deps) // DepsEntry has no dependency
	// since `DepsEntry` has no dependency, its recursive hash should just
	// be a re-hash of its single hash
	bss = ByteSSlice{deDT.hashSingle}
	bss.Sort()
	deExpectedHash, err := MD5(bss)
	assert.Nil(t, err)
	assert.Equal(t, TEST_DEKnownHashRecurse, "PROT-"+hex.EncodeToString(deExpectedHash))
	assert.Equal(t, deDT.hashRecursive, deExpectedHash)
}

// -----------------------------------------------------------------------------

func TestProtoscan_NewDescriptorTrees(t *testing.T) {
	// retrieve instanciated gogo/protobuf protofiles
	protoFilesBindings, err := BindProtofileSymbols()
	assert.Nil(t, err)
	assert.NotEmpty(t, protoFilesBindings)
	symbol := "github.com/gogo/protobuf/proto.protoFiles"
	if protoFilesBindings[symbol] == nil { // travis-ci runs tests through vendoring
		symbol = "github.com/znly/protein/vendor/" + symbol
	}
	assert.NotEmpty(t, protoFilesBindings[symbol])
	protoFiles := *protoFilesBindings[symbol]
	assert.NotEmpty(t, protoFiles)

	fdps := map[string]*descriptor.FileDescriptorProto{}
	for file, descr := range protoFiles {
		fdp, err := UnzipAndUnmarshal(descr)
		assert.Nil(t, err)
		fdps[file] = fdp
	}

	dtsByUID, err := NewDescriptorTrees(MD5, "PROT-", fdps)
	assert.Nil(t, err)

	// should at least find the `.test.TestSchema` and its nested
	// `DepsEntry` message in the DescriptorTrees
	assert.NotEmpty(t, dtsByUID)
	assert.NotNil(t, dtsByUID[TEST_TSKnownHashRecurse])
	assert.Equal(t, TEST_TSKnownName,
		dtsByUID[TEST_TSKnownHashRecurse].FQName(),
	)
	assert.Equal(t, TEST_DEKnownName,
		dtsByUID[TEST_DEKnownHashRecurse].FQName(),
	)
}
