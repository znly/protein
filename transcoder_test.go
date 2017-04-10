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

package protein

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"github.com/znly/protein/protobuf/test"
)

// -----------------------------------------------------------------------------

// TODO(cmc): add benchmarks

var trc *Transcoder

func TestMain(m *testing.M) {
	var err error
	trc, err = NewTranscoder(context.Background(),
		TranscoderGetterNoOp, TranscoderSetterNoOp,
	)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

// -----------------------------------------------------------------------------

// TODO(cmc): complete this
func TestTranscoder_localCache(t *testing.T) {
	var expectedUID string
	var revUID string
	// `.test.TestSchema` should be in there
	expectedUID = "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	revUID = trc.sm.GetByFQName(".test.TestSchema").UID
	assert.NotEmpty(t, revUID)
	assert.Equal(t, expectedUID, revUID)
	schemas, err := trc.get(context.Background(), revUID)
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
	assert.Equal(t, 2, len(schemas)) // `.test.TestSchema` + nested `DepsEntry`

	// `.test.TestSchema.DepsEntry` should be in there
	expectedUID = "PROT-d278f5561f05e68f6e68fcbc6b801d29a69b4bf6044bf3e6242ea8fe388ebd6e"
	revUID = trc.sm.GetByFQName(".test.TestSchema.DepsEntry").UID
	assert.NotEmpty(t, revUID)
	assert.Equal(t, expectedUID, revUID)
	schemas, err = trc.get(context.Background(), revUID)
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
	assert.Equal(t, 1, len(schemas)) // `.test.TestSchema.DepsEntry` only
}

// -----------------------------------------------------------------------------

func TestTranscoder_Encode(t *testing.T) {
	tsExpected := &test.TestSchema{
		Uid:    "test-uuid",
		FqName: "test-schema",
		Deps: map[string]string{
			"test": "schema",
		},
	}
	payload, err := trc.Encode(tsExpected)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	var pp ProtobufPayload
	var ts test.TestSchema
	assert.Nil(t, proto.Unmarshal(payload, &pp))
	uidExpected := "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	assert.Equal(t, uidExpected, pp.GetUID())
	assert.NotEmpty(t, pp.GetPayload())
	assert.Nil(t, proto.Unmarshal(pp.GetPayload(), &ts))
	assert.Equal(t, tsExpected, &ts)
}

// -----------------------------------------------------------------------------

var _transcoderTestSchemaXXX = &test.TestSchemaXXX{
	SchemaUID: "uuid-A",
	FQNames:   []string{"fqname-A", "fqname-B"},
	Deps: map[string]*test.TestSchemaXXX_NestedEntry{
		"dep-A": {Key: "dep-A-A", Value: "dep-A-A"},
		"dep-B": {Key: "dep-A-B", Value: "dep-A-B"},
	},
	Ids: map[int32]string{
		1: "id-A",
		2: "id-B",
	},
	Ts: types.Timestamp{
		Seconds: 42,
		Nanos:   43,
	},
	Ots: &test.OtherTestSchemaXXX{
		Ts: &types.Timestamp{
			Seconds: 942,
			Nanos:   943,
		},
	},
	Nss: []test.TestSchemaXXX_NestedEntry{
		{Key: "nss-key-A", Value: "nss-value-A"},
		{Key: "nss-key-B", Value: "nss-value-B"},
	},
}

func TestTranscoder_DecodeAs(t *testing.T) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	var ts test.TestSchemaXXX
	assert.Nil(t, trc.DecodeAs(payload, &ts))
	assert.Equal(t, _transcoderTestSchemaXXX, &ts)
}

func TestTranscoder_Decode(t *testing.T) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	v, err := trc.Decode(payload)
	assert.Nil(t, err)
	assertFieldValues(t, reflect.ValueOf(_transcoderTestSchemaXXX), v)
}

func assertFieldValues(t *testing.T, expected, actual reflect.Value) {
	switch expected.Kind() {
	case reflect.Ptr:
		assertFieldValues(t, expected.Elem(), actual.Elem())
	case reflect.Map:
		for _, ek := range expected.MapKeys() {
			assertFieldValues(t, expected.MapIndex(ek), actual.MapIndex(ek))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < expected.Len(); i++ {
			assertFieldValues(t, expected.Index(i), actual.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < expected.NumField(); i++ {
			assertFieldValues(t, expected.Field(i), actual.Field(i))
		}
	default:
		assert.True(t,
			reflect.DeepEqual(expected.Interface(), actual.Interface()),
		)
	}
}
