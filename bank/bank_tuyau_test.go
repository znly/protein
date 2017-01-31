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
	"github.com/znly/protein/protoscan"
	tuyau_client "github.com/znly/tuyauDB/client"
	tuyau_kv "github.com/znly/tuyauDB/kv"
	tuyau_pipe "github.com/znly/tuyauDB/pipe"
	tuyau_service "github.com/znly/tuyauDB/service"
)

// -----------------------------------------------------------------------------

func TestBank_Tuyau_RAM_PutGet(t *testing.T) {
	// fetched locally instanciated schemas
	schemas, err := protoscan.ScanSchemas()
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)

	// build the underlying TuyauDB components: Client{Pipe, KV}
	bufSize := uint(len(schemas) + 1) // cannot block that way
	cs, err := tuyau_client.NewSimple(
		tuyau_pipe.NewRAM(bufSize, bufSize),
		tuyau_kv.NewRAM(),
	)
	assert.Nil(t, err)
	assert.NotNil(t, cs)

	// build a simple TuyauDB Service to sync-up the underlying Pipe & KV
	// components (i.e. what's pushed into the pipe should en up in the kv
	// store)
	ctx, canceller := context.WithCancel(context.Background())
	s, err := tuyau_service.NewSimple(cs)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	go s.Run(ctx, 0)

	// build the actual Bank that integrates with the TuyauDB Client
	ty := NewTuyau(cs)
	go func() {
		for _, ps := range schemas {
			assert.Nil(t, ty.Put(ps)) // feed it all the local schemas
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

	var expectedUID string
	var revUIDs []string
	// `.protein.TestSchema` should be in there
	expectedUID = "PROT-b4f1216c74d15da21b72e7064e8a5ad1e023ee64e016a09884b01f9c2622da4b"
	revUIDs = ty.FQNameToUID(".protein.TestSchema")
	assert.NotEmpty(t, revUIDs)
	assert.Equal(t, 1, len(revUIDs))
	assert.Equal(t, expectedUID, revUIDs[0])
	schemas, err = ty.Get(revUIDs[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
	assert.Equal(t, 2, len(schemas)) // `.protein.TestSchema` + nested `DepsEntry`

	// `.protein.TestSchema.DepsEntry` should be in there
	expectedUID = "PROT-2eef830874406d6ccf9a9ae9ac787d4be60f105695fd10a91f73a84d43a235b4"
	revUIDs = ty.FQNameToUID(".protein.TestSchema.DepsEntry")
	assert.NotEmpty(t, revUIDs)
	assert.Equal(t, 1, len(revUIDs))
	assert.Equal(t, expectedUID, revUIDs[0])
	schemas, err = ty.Get(revUIDs[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
	assert.Equal(t, 1, len(schemas)) // `.protein.TestSchema.DepsEntry` only

	// destroy the underlying TuyauDB components
	err1, err2 := cs.Close()
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}
