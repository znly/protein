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

func TestProtoscan_ScanSchemas(t *testing.T) {
	schemas, err := ScanSchemas()
	assert.Nil(t, err)

	// should at least find the `.protoscan.TestSchema` and its nested
	// `DepsEntry` message in the returned protobuf schemas

	ps := schemas[tsKnownHashRecurse]
	assert.NotNil(t, ps)
	assert.Equal(t, tsKnownName, ps.GetFQName())
	assert.Equal(t, tsKnownHashRecurse, ps.GetUID())
	assert.NotNil(t, ps.GetDescr())
	assert.NotEmpty(t, ps.GetDeps())
	assert.NotNil(t, ps.GetDeps()[deKnownHashRecurse])

	de := schemas[deKnownHashRecurse]
	assert.NotNil(t, de)
	assert.Equal(t, deKnownName, de.GetFQName())
	assert.Equal(t, deKnownHashRecurse, de.GetUID())
	assert.NotNil(t, de.GetDescr())
	assert.Empty(t, de.GetDeps())
}

func TestProtoscan_BindProtofileSymbols(t *testing.T) {
	// protein/protoscan only imports a non-vendored gogo/protobuf package,
	// so the only symbol found should be the following one:
	symbol := "github.com/gogo/protobuf/proto.protoFiles"
	protoFilesBindings, err := BindProtofileSymbols()
	assert.Nil(t, err)
	assert.NotEmpty(t, protoFilesBindings)
	assert.Equal(t, 1, len(protoFilesBindings))
	assert.NotEmpty(t, protoFilesBindings[symbol])

	// at least `descriptor.proto`, `gogo.proto` & `test_schema.proto` are
	// expected to have been registered for the gogo/protobuf symbol
	registered := []string{
		"descriptor.proto",
		"gogo.proto",
		"test_schema.proto",
	}
	protoFiles := *protoFilesBindings[symbol]
	assert.NotEmpty(t, protoFiles)
	assert.True(t, len(protoFiles) >= len(registered))
	for _, fileName := range registered {
		assert.NotEmpty(t, protoFiles[fileName])
	}

	// `test_schema.proto`'s file descriptor should at least contain
	// the TestSchema message type, as well as the nested
	// TestSchema.DepsEntry type
	descr, err := UnzipAndUnmarshal(protoFiles["test_schema.proto"])
	assert.Nil(t, err)
	assert.NotNil(t, descr)
	assert.Equal(t, "test_schema.proto", descr.GetName())
	schemaMessage := descr.GetMessage("TestSchema")
	assert.NotNil(t, schemaMessage)
	assert.NotNil(t, descr.GetNestedMessage(schemaMessage, "DepsEntry"))
}
