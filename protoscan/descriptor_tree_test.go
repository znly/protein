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

func TestProtoscan_collectDescriptorTypes(t *testing.T) {
	// retrieve instanciated gogo/protobuf protofiles
	symbol := "github.com/gogo/protobuf/proto.protoFiles"
	protoFilesBindings, err := BindProtofileSymbols()
	assert.Nil(t, err)
	assert.NotEmpty(t, protoFilesBindings)
	assert.NotEmpty(t, protoFilesBindings[symbol])
	protoFiles := *protoFilesBindings[symbol]
	assert.NotEmpty(t, protoFiles)

	// keep only `protobuf_schema.proto`
	fdpRaw := protoFiles["protobuf_schema.proto"]
	assert.NotEmpty(t, fdpRaw)
	fdp, err := UnzipAndUnmarshal(fdpRaw)
	assert.Nil(t, err)
	assert.NotNil(t, fdp)
	fdps := map[string]*descriptor.FileDescriptorProto{
		"protobuf_schema.proto": fdp,
	}

	// collect DescriptorTrees for `protobuf_schema.proto`
	dtsByName, err := collectDescriptorTypes(fdps)
	assert.Nil(t, err)
	assert.NotEmpty(t, dtsByName)

	// should at least find 2 messages types here: `.protoscan.ProtobufSchema`
	// and its nested `DepsEntry` message
	assert.True(t, len(dtsByName) >= 2)

	psDescr := dtsByName[".protoscan.ProtobufSchema"]
	assert.NotNil(t, psDescr)
	assert.Nil(t, psDescr.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, ".protoscan.ProtobufSchema", psDescr.fqName)
	assert.NotNil(t, psDescr.descr)
	assert.Equal(t, "ProtobufSchema", psDescr.descr.(*descriptor.DescriptorProto).GetName())
	b, err := proto.Marshal(psDescr.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	psExpectedHash, err := ByteSSlice{b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, psExpectedHash, psDescr.hashSingle)
	assert.Nil(t, psDescr.hashRecursive)

	assert.NotNil(t, dtsByName[".protoscan.ProtobufSchema.DepsEntry"])
	deDescr := dtsByName[".protoscan.ProtobufSchema.DepsEntry"]
	assert.NotNil(t, deDescr)
	assert.Nil(t, deDescr.deps) // shouldn't have dependencies linked yet
	assert.Equal(t, ".protoscan.ProtobufSchema.DepsEntry", deDescr.fqName)
	assert.NotNil(t, deDescr.descr)
	assert.Equal(t, "DepsEntry", deDescr.descr.(*descriptor.DescriptorProto).GetName())
	b, err = proto.Marshal(deDescr.descr)
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	deExpectedHash, err := ByteSSlice{b}.Hash()
	assert.Nil(t, err)
	assert.Equal(t, deExpectedHash, deDescr.hashSingle)
	assert.Nil(t, deDescr.hashRecursive)
}
