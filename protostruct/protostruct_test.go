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
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/znly/protein/bank"
	"github.com/znly/protein/protoscan"
	tuyau_client "github.com/znly/tuyauDB/client"
	tuyau_kv "github.com/znly/tuyauDB/kv"
	tuyau_pipe "github.com/znly/tuyauDB/pipe"
	tuyau_service "github.com/znly/tuyauDB/service"

	"github.com/znly/protein/protobuf/schemas/test"
)

// -----------------------------------------------------------------------------

func TestProtostruct_CreateStructType(t *testing.T) {
	// fetched locally instanciated schemas
	schemas, err := protoscan.ScanSchemas()
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)

	// build the underlying TuyauDB components: Client{Pipe, KV}
	bufSize := uint(len(schemas) + 1) // cannot block that way
	cs, err := tuyau_client.New(
		tuyau_pipe.NewRAMConstructor(bufSize),
		tuyau_kv.NewRAMConstructor(),
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, cs)

	// build a simple TuyauDB Service to sync-up the underlying Pipe & KV
	// components (i.e. what's pushed into the pipe should en up in the kv
	// store)
	ctx, canceller := context.WithCancel(context.Background())
	s, err := tuyau_service.New(cs, 10)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	go s.Run(ctx)

	// build the actual Bank that integrates with the TuyauDB Client
	ty := bank.NewTuyau(cs)
	go func() {
		for _, ps := range schemas {
			assert.Nil(t, ty.Put(context.Background(), ps)) // feed it all the local schemas
		}
		time.Sleep(time.Millisecond * 20)
		canceller() // we're done
	}()

	<-ctx.Done()
	// At this point, all the locally-instanciated protobuf schemas should
	// have been Put() into the Bank, which Push()ed them all to its underlying
	// Tuyau Client and, hence, into the RAM-based Tuyau Pipe.
	//
	// Since a Simple Tuyau Service had been running all along, making sure the
	// underlying RAM-based Tuyau KV store was kept in synchronization with
	// the RAM-based Pipe, our Bank should now be able to retrieve any schema
	// directly from its underlying KV store.

	// ------- ^^^^^ need helpers for all this boilerplate

	expectedType := reflect.TypeOf(test.TestSchemaXXX{})
	assert.True(t, expectedType.Kind() == reflect.Struct)
	expectedFields := make(map[string]reflect.StructField, expectedType.NumField())
	for i := 0; i < expectedType.NumField(); i++ {
		f := expectedType.Field(i)
		expectedFields[f.Name] = f
	}
	assert.NotEmpty(t, expectedFields)

	uids := ty.FQNameToUID(".test.TestSchemaXXX")
	assert.NotEmpty(t, uids)
	actualType, err := CreateStructType(ty, uids[0])
	assert.Nil(t, err)
	assert.NotNil(t, actualType)
	assert.True(t, (*actualType).Kind() == reflect.Struct)
	actualFields := make(map[string]reflect.StructField, (*actualType).NumField())
	for i := 0; i < (*actualType).NumField(); i++ {
		f := (*actualType).Field(i)
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
