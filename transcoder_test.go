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
	"testing"

	"github.com/gogo/protobuf/proto"
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
//func TestTranscoder_localCache(t *testing.T) {
//var expectedUID string
//var revUIDs []string
//// `.test.TestSchema` should be in there
//expectedUID = "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
//revUIDs = trc.FQNameToUID(".test.TestSchema")
//assert.NotEmpty(t, revUIDs)
//assert.Equal(t, 1, len(revUIDs))
//assert.Equal(t, expectedUID, revUIDs[0])
//schemas, err := trc.get(context.Background(), revUIDs[0])
//assert.Nil(t, err)
//assert.NotEmpty(t, schemas)
//assert.Equal(t, 2, len(schemas)) // `.test.TestSchema` + nested `DepsEntry`

//// `.test.TestSchema.DepsEntry` should be in there
//expectedUID = "PROT-d278f5561f05e68f6e68fcbc6b801d29a69b4bf6044bf3e6242ea8fe388ebd6e"
//revUIDs = trc.FQNameToUID(".test.TestSchema.DepsEntry")
//assert.NotEmpty(t, revUIDs)
//assert.Equal(t, 1, len(revUIDs))
//assert.Equal(t, expectedUID, revUIDs[0])
//schemas, err = trc.get(context.Background(), revUIDs[0])
//assert.Nil(t, err)
//assert.NotEmpty(t, schemas)
//assert.Equal(t, 1, len(schemas)) // `.test.TestSchema.DepsEntry` only
//}

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

// TODO(cmc)
func TestTranscoder_Decode(t *testing.T)   {}
func TestTranscoder_DecodeAs(t *testing.T) {}
