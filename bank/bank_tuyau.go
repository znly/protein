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
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/znly/protein"
	tuyaudb "github.com/znly/tuyauDB"
	tuyau "github.com/znly/tuyauDB/client"
)

// -----------------------------------------------------------------------------

// Tuyau implements a Bank that integrates with znly/tuyauDB in order to keep
// its local in-memory cache in sync with a TuyauDB store.
type Tuyau struct {
	tc tuyau.Client

	schemas map[string]*protein.ProtobufSchema
}

// NewTuyau returns a new Tuyau that uses `tc` as its underlying client for
// accessing a TuyauDB store.
//
// It is the caller's responsibility to close the client once he's done with it.
func NewTuyau(tc tuyau.Client) *Tuyau {
	return &Tuyau{
		tc:      tc,
		schemas: map[string]*protein.ProtobufSchema{},
	}
}

// -----------------------------------------------------------------------------

// Get retrieves the ProtobufSchema associated with the specified identifier,
// plus all of its direct & indirect dependencies.
//
// The retrieval process is done in two steps:
// - First, the root schema, as identified by `uid`, is fetched from the local
//   in-memory cache; if it cannot be found in there, it'll be retrieved from
//   the backing TuyauDB store.
//   If it cannot be found in the TuyauDB store, then a "schema not found"
//   error is returned.
// - Second, the same process is applied for every direct & indirect dependency
//   of the root schema.
//   The only difference is that all the dependencies missing from the local
//   cache will be bulk-fetched from the TuyauDB store to avoid unnecessary
//   round-trips.
//   A "schemas not found" error is returned if one or more dependencies couldn't
//   be found.
func (t *Tuyau) Get(uid string) (map[string]*protein.ProtobufSchema, error) {
	schemas := map[string]*protein.ProtobufSchema{}

	// get root schema
	if s, ok := schemas[uid]; ok { // try the in-memory cache first..
		schemas[uid] = s
	} else { // ..then fallback on the remote tuyauDB store
		b, err := t.tc.Get(uid)
		if err != nil {
			return nil, errors.Wrapf(err, "`%s`: schema not found", uid)
		}
		var root protein.ProtobufSchema
		if err := proto.Unmarshal(b.Data, &root); err != nil {
			return nil, errors.Wrapf(err, "`%s`: invalid schema", uid)
		}
		schemas[uid] = &root
	}

	// get dependency schemas
	deps := schemas[uid].GetDeps()

	// try the in-memory cache first..
	psNotFound := make(map[string]struct{}, len(deps))
	for depUID := range deps {
		if s, ok := schemas[depUID]; ok {
			schemas[depUID] = s
			continue
		}
		psNotFound[depUID] = struct{}{}
	}
	if len(psNotFound) <= 0 { // found everything need in local cache!
		return schemas, nil
	}

	// ..then fallback on the remote tuyauDB store
	psToFetch := make([]string, 0, len(psNotFound))
	for depUID := range psNotFound {
		psToFetch = append(psToFetch, depUID)
	}
	blobs, err := t.tc.GetMulti(psToFetch)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, b := range blobs {
		delete(psNotFound, b.Key) // it's been found!
	}
	if len(psNotFound) > 0 {
		err := errors.Errorf("one or more dependencies couldn't be found")
		for depUID := range psNotFound {
			err = errors.Wrapf(err, "`%s`: dependency not found", depUID)
		}
		return nil, err
	}
	for _, b := range blobs {
		var ps protein.ProtobufSchema
		if err := proto.Unmarshal(b.Data, &ps); err != nil {
			return nil, errors.Wrapf(err, "`%s`: invalid schema (dependency)", b.Key)
		}
		schemas[b.Key] = &ps
	}

	return schemas, nil
}

// Put synchronously adds the specified ProtobufSchemas to the local in-memory
// cache; then pushes them to the underlying tuyau client's pipe.
// Whether this push is synchronous or not depends on the implementation
// of the tuyau.Client used.
//
// Put doesn't care about pre-existing keys: if a schema with the same key
// already exist, it will be overwritten; both in the local cache as well in the
// TuyauDB store.
func (t *Tuyau) Put(pss ...*protein.ProtobufSchema) error {
	blobs := make([]*tuyaudb.Blob, 0, len(pss))
	var b []byte
	var err error
	for _, ps := range pss {
		b, err = proto.Marshal(ps)
		if err != nil {
			return errors.WithStack(err)
		}
		blobs = append(blobs, &tuyaudb.Blob{
			Key: ps.GetUID(), Data: b, TTL: 0, Flags: 0,
		})
		t.schemas[ps.GetUID()] = ps
	}
	return t.tc.Push(blobs...)
}
