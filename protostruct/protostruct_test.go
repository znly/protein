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

package protostruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/znly/protein"
	"github.com/znly/protein/protobuf/test"
)

// -----------------------------------------------------------------------------

func TestProtostruct_CreateStructType(t *testing.T) {
	// fetched locally instanciated schemas
	sm, err := protein.ScanSchemas()
	assert.Nil(t, err)
	assert.NotEmpty(t, sm)

	expectedType := reflect.TypeOf(test.TestSchemaXXX{})
	assert.True(t, expectedType.Kind() == reflect.Struct)
	expectedFields := make(map[string]reflect.StructField, expectedType.NumField())
	for i := 0; i < expectedType.NumField(); i++ {
		f := expectedType.Field(i)
		expectedFields[f.Name] = f
	}
	assert.NotEmpty(t, expectedFields)

	uid := sm.GetByFQName(".test.TestSchemaXXX").UID
	assert.NotEmpty(t, uid)
	actualType, err := CreateStructType(uid, sm)
	assert.Nil(t, err)
	assert.NotNil(t, actualType)
	assert.True(t, actualType.Kind() == reflect.Struct)
	actualFields := make(map[string]reflect.StructField, actualType.NumField())
	for i := 0; i < actualType.NumField(); i++ {
		f := actualType.Field(i)
		actualFields[f.Name] = f
	}
	assert.NotEmpty(t, actualFields)

	var expectedField, actualField reflect.StructField

	// Attributes:
	//   - name: "uid"
	//   - custom_name: "SchemaUID"
	//   - type: string
	//   - custom_type: ""
	//   - id: 1
	expectedField = expectedFields["SchemaUID"]
	actualField = actualFields["SchemaUID"]
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, expectedField.Name, actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	assert.Equal(t, expectedField.Type, actualField.Type)

	// Attributes:
	//   - name: "fq_name"
	//   - custom_name: "FQNames"
	//   - type: repeated string
	//   - custom_type: ""
	//   - id: 2
	expectedField = expectedFields["FQNames"]
	actualField = actualFields["FQNames"]
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, expectedField.Name, actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	assert.Equal(t, expectedField.Type, actualField.Type)

	// Attributes:
	//   - name: "deps"
	//   - custom_name: ""
	//   - type: map<string, .test.TestSchemaXXX.NestedEntry>
	//   - custom_type: ""
	//   - id: 4
	expectedField = expectedFields["Deps"]
	actualField = actualFields["Deps"]
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, expectedField.Name, actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	// assert.Equal(t, expectedField.Type, actualField.Type) TODO

	// Attributes:
	//   - name: "ids"
	//   - custom_name: ""
	//   - type: map<int32, string>
	//   - custom_type: ""
	//   - id: 12
	expectedField = expectedFields["Ids"]
	actualField = actualFields["IDs"] // protein pretty-cases acronyms
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, "IDs", actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	assert.Equal(t, expectedField.Type, actualField.Type)

	// Attributes:
	//   - name: "ts"
	//   - custom_name: ""
	//   - type: .google.protobuf.Timestamp
	//   - custom_type: ""
	//   - id: 7
	expectedField = expectedFields["Ts"]
	actualField = actualFields["TS"] // protein pretty-cases acronyms
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, "TS", actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	// assert.Equal(t, expectedField.Type, actualField.Type) TODO

	// Attributes:
	//   - name: "ots"
	//   - custom_name: ""
	//   - type: .test.OtherTestSchemaXXX
	//   - custom_type: ""
	//   - id: 9
	expectedField = expectedFields["Ots"]
	actualField = actualFields["Ots"]
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, expectedField.Name, actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	// assert.Equal(t, expectedField.Type, actualField.Type) TODO

	// Attributes:
	//   - name: "nss"
	//   - custom_name: ""
	//   - type: repeated .test.TestSchemaXXX.NestedEntry
	//   - custom_type: ""
	//   - id: 8
	expectedField = expectedFields["Nss"]
	actualField = actualFields["Nss"]
	assert.Equal(t, expectedField.Anonymous, actualField.Anonymous)
	assert.Equal(t, expectedField.Name, actualField.Name)
	assert.Equal(t, expectedField.PkgPath, actualField.PkgPath)
	assert.Equal(t, expectedField.Tag, actualField.Tag)
	// assert.Equal(t, expectedField.Type, actualField.Type) TODO
}
