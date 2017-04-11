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
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gogo/protobuf/proto"
	proto_vanilla "github.com/golang/protobuf/proto"

	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

// TODO(cmc)
//
// Transcoder implements a Wirer that integrates with a Bank in order to augment
// the protobuf payloads that it encodes with additional versioning metadata.
//
// These metadata are then used by the internal decoder of the versioned Wirer
// to determinate how to decode an incoming payload on the wire.
type Transcoder struct {
	sm *SchemaMap
	// TODO(cmc)
	getter func(ctx context.Context, uid string) (*ProtobufSchema, error)
	setter func(ctx context.Context, ps *ProtobufSchema) error

	typeCache map[string]reflect.Type
}

// TODO(cmc)
var (
	TranscoderGetterNoOp = func(ctx context.Context, uid string) (*ProtobufSchema, error) {
		// TODO(cmc): real error
		return nil, errors.New(fmt.Sprintf("`%s`: schema not found", uid))
	}
	TranscoderSetterNoOp = func(context.Context, *ProtobufSchema) error {
		return nil
	}
	// TODO(cmc): add helpers
)

// TODO(cmc)
func NewTranscoder(ctx context.Context,
	getter func(ctx context.Context, uid string) (*ProtobufSchema, error),
	setter func(ctx context.Context, ps *ProtobufSchema) error,
) (*Transcoder, error) {
	if setter == nil || getter == nil {
		// TODO(cmc): real error
		return nil, errors.New("getter and/or setter cannot be nil")
	}

	sm, err := ScanSchemas()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sm.ForEach(func(ps *ProtobufSchema) error {
		select {
		case <-ctx.Done():
			return errors.WithStack(ctx.Err())
		default:
		}
		return setter(ctx, ps)
	})

	return &Transcoder{
		sm:        sm,
		getter:    getter,
		setter:    setter,
		typeCache: map[string]reflect.Type{},
	}, nil
}

// -----------------------------------------------------------------------------

// TODO(cmc)
//
// get retrieves the ProtobufSchema associated with the specified identifier,
// plus all of its direct & indirect dependencies flattened in a map.
//
// The retrieval process is done in two steps:
// - First, the root schema, as identified by `uid`, is fetched from the local
//   local cache; if it cannot be found in there, it'll be retrieved from
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
func (t *Transcoder) get(
	ctx context.Context, uid string,
) (map[string]*ProtobufSchema, error) {
	schemas := map[string]*ProtobufSchema{}

	// get root schema
	if ps := t.sm.GetByUID(uid); ps != nil { // try the local cache first..
		schemas[uid] = ps
	} else { // ..then fallback on user-defined getter
		ps, err := t.getter(ctx, uid)
		if err != nil {
			// TODO(cmc): real error
			return nil, errors.Wrapf(err, "`%s`: schema not found", uid)
		}
		t.sm.Add(map[string]*ProtobufSchema{uid: ps}) // upsert local-cache
	}

	// get dependencies
	deps := schemas[uid].GetDeps()

	// try the local cache first..
	psNotFound := make(map[string]struct{}, len(deps))
	for depUID := range deps {
		if ps := t.sm.GetByUID(depUID); ps != nil {
			schemas[depUID] = ps
		} else {
			psNotFound[depUID] = struct{}{}
		}
	}
	if len(psNotFound) <= 0 { // found everything needed in local cache!
		return schemas, nil
	}

	// ..then fallback on user-defined getter
	var err error
	for depUID := range psNotFound {
		var ps *ProtobufSchema
		ps, err = t.getter(ctx, depUID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		delete(psNotFound, depUID)                    // it's been found!
		t.sm.Add(map[string]*ProtobufSchema{uid: ps}) // upsert local-cache
	}
	if len(psNotFound) > 0 {
		// TODO(cmc): real error
		err := errors.Errorf("one or more dependencies couldn't be found")
		for depUID := range psNotFound {
			err = errors.Wrapf(err, "`%s`: dependency not found", depUID)
		}
		return nil, err
	}

	return schemas, nil
}

// -----------------------------------------------------------------------------

// TODO(cmc)
//
// Encode marshals the given protobuf message then wraps it up into a
// ProtobufPayload object that adds additional versioning metadata.
//
// Encode uses the message's fully-qualified name to reverse-lookup its UID.
// Note that a single FQ-name might point to multiple UIDs if multiple versions
// of the associated message are currently availaible in the bank.
// When this happens, the first UID from the returned list will be used (and
// since this list is randomly-ordered, effectively a random UID will be used).
// ---> this won't do remote calls
func (t *Transcoder) Encode(msg proto.Message, fqName ...string) ([]byte, error) {
	// marshal the actual protobuf message
	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO(cmc): explain this mess
	var fqn string
	if len(fqName) > 0 {
		fqn = fqName[0]
	} else if fqn = proto_vanilla.MessageName(msg); len(fqn) > 0 {
	} else if fqn = proto.MessageName(msg); len(fqn) > 0 {
	} else {
		// TODO(cmc): real error
		return nil, errors.Errorf("cannot encode, unknown protobuf schema")
	}

	// find the first UID associated with the fully-qualified name of `msg`
	ps := t.sm.GetByFQName("." + fqn)
	if ps == nil {
		// TODO(cmc): real error
		return nil, errors.Errorf("`%s`: FQ-name not found", fqn)
	}
	// wrap the marshaled payload within a ProtobufPayload message
	pp := &ProtobufPayload{
		UID:     ps.UID,
		Payload: payload,
	}
	// marshal the ProtobufPayload
	return proto.Marshal(pp)
}

// -----------------------------------------------------------------------------

// We need to access protobuf's internal decoding machinery: the `go:linkname`
// directive instructs the compiler to declare a local symbol as an alias
// for an external one, even if it's private.
// This allows us to bind to the private `unmarshalType` method of the
// `proto.Buffer` class, which does the actual work of computing the
// necessary struct tags for a given protobuf field.
//
// `unmarshalType` is actually a method of the `proto.Buffer` class, hence the
// `b` given as first parameter will be used as "this".
//
// Due to the way Go mangles symbol names when using vendoring, the go:linkname
// clause is automatically generated via linkname-gen[1].
// [1] https://github.com/znly/linkname-gen.
//
//go:generate linkname-gen -symbol "github.com/gogo/protobuf/proto.(*Buffer).unmarshalType" -def "func unmarshalType(*proto.Buffer, reflect.Type, *proto.StructProperties, bool, unsafe.Pointer) error"

// Decode decodes the `payload` into a dynamically-defined structure type.
func (t *Transcoder) Decode(payload []byte) (reflect.Value, error) {
	var pp ProtobufPayload
	if err := proto.Unmarshal(payload, &pp); err != nil {
		return reflect.ValueOf(nil), errors.WithStack(err)
	}

	// fetch structure-type from cache, or create it if it doesn't exist
	var structType reflect.Type
	var ok bool
	uid := pp.GetUID()
	if structType, ok = t.typeCache[uid]; !ok {
		st, err := CreateStructType(uid, t.sm)
		if err != nil {
			return reflect.ValueOf(nil), errors.WithStack(err)
		}
		if st.Kind() != reflect.Struct {
			// TODO(cmc): real error
			return reflect.ValueOf(nil), errors.Errorf(
				"`%s`: not a struct type", structType,
			)
		}
		structType = st
		t.typeCache[uid] = st // upsert type-cache
	}

	// allocate a new structure using the given type definition, the
	// returned `reflect.Value`'s underlying type is a pointer to struct
	obj := reflect.New(structType)

	b := proto.NewBuffer(pp.GetPayload())
	unmarshalType(b,
		// the structure definition, computed at runtime
		structType,
		// the protobuf properties of the struct, computed via its struct tags
		proto.GetProperties(structType),
		// is_group, deprecated
		false,
		// the address we want to deserialize to
		unsafe.Pointer(obj.Elem().Addr().Pointer()),
	)

	return obj, nil
}

// TODO(cmc): doc & test
func (t *Transcoder) DecodeAs(payload []byte, dst proto.Message) error {
	var ps ProtobufPayload
	if err := proto.Unmarshal(payload, &ps); err != nil {
		return errors.WithStack(err)
	}
	return proto.Unmarshal(ps.Payload, dst)
}
