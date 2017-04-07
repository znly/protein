// +build ignore

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

package bank

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/znly/protein"
	"github.com/znly/protein/protoscan"
	tuyaudb "github.com/znly/tuyauDB"
	tuyaudb_kv "github.com/znly/tuyauDB/kv"
	tuyaudb_pipe "github.com/znly/tuyauDB/pipe"
	tuyaudb_service "github.com/znly/tuyauDB/service"

	_ "github.com/znly/protein/protobuf/test"
)

// -----------------------------------------------------------------------------

func TestBank_Tuyau_RAM_PutGet(t *testing.T) {
	// fetched locally instanciated schemas
	schems, err := protoscan.ScanSchemas()
	assert.Nil(t, err)
	assert.NotEmpty(t, schems)

	// build the underlying TuyauDB components: Client{Pipe, KV}
	bufSize := uint(len(schems) + 1) // cannot block that way
	cs, err := tuyaudb.NewClient(
		tuyaudb_pipe.NewRAMConstructor(bufSize),
		tuyaudb_kv.NewRAMConstructor(),
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, cs)

	// build a simple TuyauDB Service to sync-up the underlying Pipe & KV
	// components (i.e. what's pushed into the pipe should end up in the kv
	// store)
	ctx, canceller := context.WithCancel(context.Background())
	s, err := tuyaudb_service.New(cs, 10, false)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	go s.Run(ctx, 10)

	// build the actual Bank that integrates with the TuyauDB Client
	ty := NewTuyau(cs)
	go func() {
		schemsL := make([]*protein.ProtobufSchema, 0, len(schems))
		for _, ps := range schems {
			schemsL = append(schemsL, ps)
		}
		assert.Nil(t, ty.Put(context.Background(), schemsL...)) // feed it all the local schemas
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

	var expectedUID string
	var revUIDs []string
	// `.test.TestSchema` should be in there
	expectedUID = "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	revUIDs = ty.FQNameToUID(".test.TestSchema")
	assert.NotEmpty(t, revUIDs)
	assert.Equal(t, 1, len(revUIDs))
	assert.Equal(t, expectedUID, revUIDs[0])
	schems, err = ty.Get(context.Background(), revUIDs[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, schems)
	assert.Equal(t, 2, len(schems)) // `.test.TestSchema` + nested `DepsEntry`

	// `.test.TestSchema.DepsEntry` should be in there
	expectedUID = "PROT-d278f5561f05e68f6e68fcbc6b801d29a69b4bf6044bf3e6242ea8fe388ebd6e"
	revUIDs = ty.FQNameToUID(".test.TestSchema.DepsEntry")
	assert.NotEmpty(t, revUIDs)
	assert.Equal(t, 1, len(revUIDs))
	assert.Equal(t, expectedUID, revUIDs[0])
	schems, err = ty.Get(context.Background(), revUIDs[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, schems)
	assert.Equal(t, 1, len(schems)) // `.test.TestSchema.DepsEntry` only
}
