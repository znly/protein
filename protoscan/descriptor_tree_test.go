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
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

func _collectProtobufSchemaTrees(t *testing.T) map[string]*DescriptorTree {
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

	// should at least find 2 messages types here: `.protoscan.ProtobufSchema`
	// and its nested `DepsEntry` message
	assert.True(t, len(dtsByName) >= 2)
	assert.NotNil(t, dtsByName[".protoscan.ProtobufSchema"])
	assert.NotNil(t, dtsByName[".protoscan.ProtobufSchema.DepsEntry"])

	return dtsByName
}

func TestProtoscan_collectDescriptorTypes(t *testing.T) {
	dtsByName := _collectProtobufSchemaTrees(t)

	psDT := dtsByName[".protoscan.ProtobufSchema"]
	assert.Nil(t, psDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, ".protoscan.ProtobufSchema", psDT.fqName)
	assert.NotNil(t, psDT.descr)
	assert.Equal(t, "ProtobufSchema", psDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err := proto.Marshal(psDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	psExpectedHash, err := ByteSSlice{b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, psExpectedHash, psDT.hashSingle)
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[".protoscan.ProtobufSchema.DepsEntry"]
	assert.Nil(t, deDT.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, ".protoscan.ProtobufSchema.DepsEntry", deDT.fqName)
	assert.NotNil(t, deDT.descr)
	assert.Equal(t, "DepsEntry", deDT.descr.(*descriptor.DescriptorProto).GetName())
	b, err = proto.Marshal(deDT.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	deExpectedHash, err := ByteSSlice{b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, deExpectedHash, deDT.hashSingle)
	assert.Nil(t, deDT.hashRecursive)
}

// -----------------------------------------------------------------------------

func TestProtoscan_DescriptorTree_computeDependencyLinks(t *testing.T) {
	dtsByName := _collectProtobufSchemaTrees(t)

	psDT := dtsByName[".protoscan.ProtobufSchema"]
	assert.Nil(t, psDT.computeDependencyLinks(dtsByName))
	assert.NotEmpty(t, psDT.deps)
	depsMap := make(map[string]*DescriptorTree, len(psDT.deps))
	for _, dep := range psDT.deps {
		depsMap[dep.FQName()] = dep
	}
	// should at least find a dependency to `DepsEntry` in here
	assert.Contains(t, depsMap, ".protoscan.ProtobufSchema.DepsEntry")
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, psDT.hashRecursive)

	deDT := dtsByName[".protoscan.ProtobufSchema.DepsEntry"]
	assert.Nil(t, deDT.computeDependencyLinks(dtsByName))
	assert.Empty(t, deDT.deps) // DepsEntry has no dependency
	// recursive hash still shouldn't have been computed at this point
	assert.Nil(t, deDT.hashRecursive)
}
