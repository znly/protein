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

	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

func TestProtoscan_BindProtofileSymbols(t *testing.T) {
	// protein/protoscan only imports a non-vendored gogo/protobuf package,
	// so the only symbol found should be the following one:
	symbol := "github.com/gogo/protobuf/proto.protoFiles"
	protoFilesBindings, err := BindProtofileSymbols()
	assert.Nil(t, err)
	assert.NotNil(t, protoFilesBindings)
	assert.Equal(t, 1, len(protoFilesBindings))
	assert.NotNil(t, protoFilesBindings[symbol])

	// at least `descriptor.proto`, `gogo.proto` & `protobuf_schema.proto` are
	// expected to have been registered for the gogo/protobuf symbol
	registered := []string{
		"descriptor.proto",
		"gogo.proto",
		"protobuf_schema.proto",
	}
	protoFiles := *protoFilesBindings[symbol]
	assert.NotNil(t, protoFiles)
	assert.True(t, len(protoFiles) >= len(registered))
	for _, fileName := range registered {
		assert.NotNil(t, protoFiles[fileName])
	}

	// `protobuf_schema.proto`'s file descriptor should at least contain
	// the ProtobufSchema message type, as well as the nested
	// ProtobufSchema.DepsEntry type
	descr, err := UnzipAndUnmarshal(protoFiles["protobuf_schema.proto"])
	assert.Nil(t, err)
	assert.NotNil(t, descr)
	assert.Equal(t, "protobuf_schema.proto", descr.GetName())
	schemaMessage := descr.GetMessage("ProtobufSchema")
	assert.NotNil(t, schemaMessage)
	assert.NotNil(t, descr.GetNestedMessage(schemaMessage, "DepsEntry"))
}
