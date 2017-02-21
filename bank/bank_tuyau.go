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

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/znly/protein"
	tuyau "github.com/znly/tuyauDB"
)

// -----------------------------------------------------------------------------

// Tuyau implements a Bank that integrates with znly/tuyauDB in order to keep
// its local in-memory cache in sync with a TuyauDB store.
type Tuyau struct {
	c *tuyau.Client

	schemas map[string]*protein.ProtobufSchema
	// reverse-mapping of fully-qualified names to UIDs
	revmap map[string][]string
}

// NewTuyau returns a new Tuyau that uses `c` as its underlying client for
// accessing a TuyauDB store.
//
// It is the caller's responsibility to close the client once he's done with it.
func NewTuyau(c *tuyau.Client) *Tuyau {
	return &Tuyau{
		c:       c,
		schemas: map[string]*protein.ProtobufSchema{},
		revmap:  map[string][]string{},
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
func (t *Tuyau) Get(ctx context.Context, uid string) (map[string]*protein.ProtobufSchema, error) {
	schemas := map[string]*protein.ProtobufSchema{}

	// get root schema
	if s, ok := schemas[uid]; ok { // try the in-memory cache first..
		schemas[uid] = s
	} else { // ..then fallback on the remote tuyauDB store
		b, err := t.c.Get(ctx, uid)
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
	if len(psNotFound) <= 0 { // found everything needed in local cache!
		return schemas, nil
	}

	// ..then fallback on the remote tuyauDB store
	psToFetch := make([]string, 0, len(psNotFound))
	for depUID := range psNotFound {
		psToFetch = append(psToFetch, depUID)
	}
	blobs, err := t.c.GetMulti(ctx, psToFetch)
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

// FQNameToUID returns the UID associated with the given fully-qualified name.
//
// It is possible that multiple versions of a schema identified by a FQ name
// are currently available in the bank; in which case all of the associated UIDs
// will be returned to the caller, *in random order*.
//
// The reverse-mapping is pre-computed; don't hesitate to call this method, it'll
// be real fast.
//
// It returns nil if `fqName` doesn't match any schema in the bank.
func (t *Tuyau) FQNameToUID(fqName string) []string { return t.revmap[fqName] }

// Put synchronously adds the specified ProtobufSchemas to the local in-memory
// cache; then pushes them to the underlying tuyau client's pipe.
// Whether this push is synchronous or not depends on the implementation
// of the tuyau.Client used.
//
// Put doesn't care about pre-existing keys: if a schema with the same key
// already exist, it will be overwritten; both in the local cache as well in the
// TuyauDB store.
//
// TODO(cmc): note about CAS that doesn't matter here
func (t *Tuyau) Put(ctx context.Context, pss ...*protein.ProtobufSchema) error {
	blobs := make([]*tuyau.Blob, 0, len(pss))
	var b []byte
	var err error
	for _, ps := range pss {
		b, err = proto.Marshal(ps)
		if err != nil {
			return errors.WithStack(err)
		}
		uid := ps.GetUID()
		blobs = append(blobs, &tuyau.Blob{
			Key: uid, Data: b, TTL: 0, Flags: 0,
		})
		t.schemas[uid] = ps
		t.revmap[ps.GetFQName()] = append(t.revmap[ps.GetFQName()], uid)
	}
	if err := t.c.Push(ctx, blobs...); err != nil {
		return errors.WithStack(err)
	}
	if err := t.c.SetMulti(ctx, blobs); err != nil {
		if errors.Cause(err) != tuyau.KVErrOpNotSupported {
			return errors.WithStack(err)
		}
		for _, b := range blobs {
			if err := t.c.Set(ctx, b); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}
