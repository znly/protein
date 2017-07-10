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
	"reflect"
	"sync"
	"unsafe"

	"github.com/gogo/protobuf/proto"
	proto_vanilla "github.com/golang/protobuf/proto"
	"github.com/znly/protein/failure"
	"github.com/znly/protein/protoscan"

	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

// A TranscoderGetter is called by the `Transcoder` when it cannot find a
// specific `schemaUID` in its local `SchemaMap`.
//
// The function returns a byte-slice that will be deserialized into a
// `ProtobufSchema` by a `TranscoderDeserializer` (see below).
//
// A `TranscoderGetter` is typically used to fetch `ProtobufSchema`s from a
// remote data-store.
// To that end, several ready-to-use implementations are provided by this
// package for different protocols: memcached, redis & CQL (i.e. cassandra).
// See `transcoder_helpers.go` for more information.
//
// The default `TranscoderGetter` always returns a not-found error.
type TranscoderGetter func(ctx context.Context, schemaUID string) ([]byte, error)

// A TranscoderSetter is called by the `Transcoder` for every schema that it
// can find in memory.
//
// The function receives a byte-slice that corresponds to a `ProtobufSchema`
// which has been previously serialized by a `TranscoderSerializer` (see below).
//
// A `TranscoderSetter` is typically used to push the local `ProtobufSchema`s
// sniffed from memory into a remote data-store.
// To that end, several ready-to-use implementations are provided by this
// package for different protocols: memcached, redis & CQL (i.e. cassandra).
// See `transcoder_helpers.go` for more information.
//
// The default `TranscoderSetter` is a no-op.
type TranscoderSetter func(ctx context.Context, schemaUID string, payload []byte) error

// A TranscoderSerializer is used to serialize `ProtobufSchema`s before passing
// them to a `TranscoderSetter`.
// See `TranscoderSetter` documentation for more information.
//
// The default `TranscoderSerializer` wraps the schema within a `ProtobufPayload`;
// i.e. it uses Protein's encoding to encode the schema.
type TranscoderSerializer func(ps *ProtobufSchema) ([]byte, error)

// A TranscoderDeserializer is used to deserialize the payloads returned by
// a `TranscoderGetter` into a `ProtobufSchema`.
// See `TranscoderGetter` documentation for more information.
//
// The default `TranscoderDeserializer` unwraps the schema from its
// `ProtobufPayload` wrapper; i.e. it uses Protein's decoding to decode the
// schema.
type TranscoderDeserializer func(payload []byte, ps *ProtobufSchema) error

// -----------------------------------------------------------------------------

// A TranscoderOpt is passed to the `Transcoder` constructor to configure
// various options.
type TranscoderOpt func(trc *Transcoder)

var (
	// TranscoderOptGetter is used to configure the `TranscoderGetter` used by
	// the `Transcoder`.
	// See `TranscoderGetter` documentation for more information.
	TranscoderOptGetter = func(getter TranscoderGetter) TranscoderOpt {
		return func(trc *Transcoder) { trc.getter = getter }
	}
	// TranscoderOptSetter is used to configure the `TranscoderSetter` used by
	// the `Transcoder`.
	// See `TranscoderSetter` documentation for more information.
	TranscoderOptSetter = func(setter TranscoderSetter) TranscoderOpt {
		return func(trc *Transcoder) { trc.setter = setter }
	}
	// TranscoderOptSerializer is used to configure the `TranscoderSerializer`
	// used by the `Transcoder`.
	// See `TranscoderSerializer` documentation for more information.
	TranscoderOptSerializer = func(serializer TranscoderSerializer) TranscoderOpt {
		return func(trc *Transcoder) { trc.serializer = serializer }
	}
	// TranscoderOptDeserializer is used to configure the `TranscoderDeserializer`
	// used by the `Transcoder`.
	// See `TranscoderDeserializer` documentation for more information.
	TranscoderOptDeserializer = func(deserializer TranscoderDeserializer) TranscoderOpt {
		return func(trc *Transcoder) { trc.deserializer = deserializer }
	}
)

// -----------------------------------------------------------------------------

// A Transcoder is a protobuf encoder/decoder with schema-versioning as well as
// runtime-decoding capabilities.
type Transcoder struct {
	sm *SchemaMap

	getter       TranscoderGetter
	setter       TranscoderSetter
	serializer   TranscoderSerializer
	deserializer TranscoderDeserializer

	typeCacheLock *sync.RWMutex
	typeCache     map[string]reflect.Type
}

// NewTranscoder returns a new `Transcoder`.
//
// See `ScanSchemas`'s documentation for more information regarding the use
// of `hasher` and `hashPrefix`.
//
// See `TranscoderOpt`'s documentation for the list of available options.
//
// The given context is passed to the user-specified `TranscoderSetter`, if any.
func NewTranscoder(ctx context.Context,
	hasher protoscan.Hasher, hashPrefix string, opts ...TranscoderOpt,
) (*Transcoder, error) {
	sm, err := ScanSchemas(hasher, hashPrefix)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewTranscoderFromSchemaMap(ctx, sm, opts...)
}

// NewTranscoderFromSchemaMap returns a new `Transcoder` backed by the
// user-specified `sm` schema-map.
// This is reserved for advanced usages.
//
// When using the vanilla `NewTranscoder` constructor, the schema-map is
// internally computed using the `ScanSchemas` function of the protoscan API.
// This constructor allows the developer to build this map themselves when
// needed; more often than not, this is achieved by using the `LoadSchemas`
// function from the protoscan API.
func NewTranscoderFromSchemaMap(ctx context.Context,
	sm *SchemaMap, opts ...TranscoderOpt,
) (*Transcoder, error) {
	t := &Transcoder{}
	t.sm = sm
	t.typeCacheLock = &sync.RWMutex{}
	t.typeCache = map[string]reflect.Type{}

	opts = append([]TranscoderOpt{
		/* default getter: always returns `ErrSchemaNotFound` */
		TranscoderOptGetter(func(
			ctx context.Context, schemaUID string,
		) ([]byte, error) {
			/* err not found */
			return nil, errors.Wrapf(failure.ErrSchemaNotFound,
				"`%s`: no schema with this schemaUID", schemaUID,
			)
		}),
		/* default setter: no-op */
		TranscoderOptSetter(func(context.Context, string, []byte) error {
			return nil
		}),
		/* default serializer: wraps the `ProtobufSchema` within a `ProtobufPayload` */
		TranscoderOptSerializer(func(ps *ProtobufSchema) ([]byte, error) {
			return t.Encode(ps)
		}),
		/* default deserializer: unwraps a `ProtobufPayload` */
		TranscoderOptDeserializer(func(payload []byte, ps *ProtobufSchema) error {
			return t.DecodeAs(payload, ps)
		}),
	}, opts...)
	for _, opt := range opts {
		opt(t)
	}

	if err := sm.ForEach(func(ps *ProtobufSchema) error {
		select {
		case <-ctx.Done():
			return errors.WithStack(ctx.Err())
		default:
		}
		b, err := t.serializer(ps)
		if err != nil {
			return errors.WithStack(err)
		}
		return t.setter(ctx, ps.GetSchemaUID(), b)
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return t, nil
}

// -----------------------------------------------------------------------------

// GetAndUpsert retrieves the `ProtobufSchema` associated with the specified
// `schemaUID`, plus all of its direct & indirect dependencies.
//
// The retrieval process is done in two steps:
//
// - First, the root schema, as identified by `schemaUID`, is fetched from the
// local `SchemaMap`; if it cannot be found in there, it'll try to retrieve
// it via the user-defined `TranscoderGetter`, as passed to the constructor
// of the `Transcoder`.
// If it cannot be found in there either, then a schema-not-found error is
// returned.
//
// - Second, this exact same process is applied for every direct & indirect
// dependency of the root schema.
// Once again, a schema-not-found error is returned if one or more dependency
// couldn't be found (the returned error does indicate which of them).
//
// The `ProtobufSchema`s found during this process are both:
// - added to the local `SchemaMap` so that they don't need to be searched for
// ever again during the lifetime of this `Transcoder`, and
// - returned to the caller as flattened map.
func (t *Transcoder) GetAndUpsert(
	ctx context.Context, schemaUID string,
) (map[string]*ProtobufSchema, error) {
	schemas := map[string]*ProtobufSchema{}

	// get root schema
	if ps := t.sm.GetByUID(schemaUID); ps != nil { // try the local `SchemaMap` first..
		schemas[schemaUID] = ps
	} else { // ..then fallback on user-defined getter
		b, err := t.getter(ctx, schemaUID)
		if err != nil {
			return nil, errors.Wrapf(err, "`%s`: schema not found", schemaUID)
		}
		var ps ProtobufSchema
		if err := t.deserializer(b, &ps); err != nil {
			return nil, errors.WithStack(err)
		}
		schemas[schemaUID] = &ps
		t.sm.Add(map[string]*ProtobufSchema{schemaUID: &ps}) // upsert `SchemaMap`
	}

	// get dependencies
	deps := schemas[schemaUID].GetDeps()

	// try the local `SchemaMap` first..
	psNotFound := make(map[string]struct{}, len(deps))
	for depUID := range deps {
		if ps := t.sm.GetByUID(depUID); ps != nil {
			schemas[depUID] = ps
			t.sm.Add(map[string]*ProtobufSchema{depUID: ps}) // upsert `SchemaMap`
		} else {
			psNotFound[depUID] = struct{}{}
		}
	}
	if len(psNotFound) <= 0 { // found everything needed in local `SchemaMap`!
		return schemas, nil
	}

	// ..then fallback on user-defined getter
	var b []byte
	var err error
	for depUID := range psNotFound {
		b, err = t.getter(ctx, depUID)
		if err != nil {
			return nil, errors.Wrapf(err, "`%s`: schema not found", depUID)
		}
		var ps ProtobufSchema
		if err := t.deserializer(b, &ps); err != nil {
			return nil, errors.WithStack(err)
		}
		delete(psNotFound, depUID) // it's been found!
		schemas[depUID] = &ps
		t.sm.Add(map[string]*ProtobufSchema{depUID: &ps}) // upsert `SchemaMap`
	}
	if len(psNotFound) > 0 {
		err := errors.Errorf("one or more dependencies couldn't be found")
		for depUID := range psNotFound {
			err = errors.Wrapf(err, "`%s`: dependency not found", depUID)
		}
		return nil, err
	}

	return schemas, nil
}

// FQName returns the fully-qualified name of the protobuf schema associated
// with `schemaUID`.
//
// Iff this schema cannot be found in the local cache, it'll try and fetch it
// from the remote registry via a call to `GetAndUpsert`.
//
// An empty string is returned if the schema is found neither locally nor
// remotely.
func (t *Transcoder) FQName(ctx context.Context, schemaUID string) string {
	if ps := t.sm.GetByUID(schemaUID); ps != nil {
		return ps.FQName
	}
	pss, err := t.GetAndUpsert(ctx, schemaUID)
	if err != nil {
		return ""
	}
	return pss[schemaUID].FQName
}

// -----------------------------------------------------------------------------

// Encode bundles the given protobuf `Message` and its associated versioning
// metadata together within a `ProtobufPayload`, marshals it all together in a
// byte-slice then returns the result.
//
// `Encode` needs the message's fully-qualified name in order to reverse-lookup
// its schemaUID (i.e. its versioning hash).
//
// In order to find this name, it will look at different places until either one
// of those does return a result or none of them does, in which case the
// encoding will fail. In search order, those places are:
// 1. first, the `fqName` parameter is checked; if it isn't set, then
// 2. the `golang/protobuf` package is queried for the FQN; if it isn't
// available there then
// 3. finally, the `gogo/protobuf` package is queried too, as a last resort.
//
// Note that a single fully-qualified name might point to multiple schemaUIDs
// if multiple versions of that schema are currently available in the `SchemaMap`.
// When this happens, the first schemaUID from the list will be used, which
// corresponds to the first version of the schema to have ever been added to
// the local `SchemaMap` (i.e. the oldest one).
func (t *Transcoder) Encode(msg proto.Message, fqName ...string) ([]byte, error) {
	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// find the fully-qualified name of `msg`'s schema
	var fqn string
	if len(fqName) > 0 {
		fqn = fqName[0]
	} else if fqn = proto_vanilla.MessageName(msg); len(fqn) > 0 {
	} else if fqn = proto.MessageName(msg); len(fqn) > 0 {
	} else {
		return nil, errors.Errorf("cannot encode, unknown protobuf schema")
	}
	fqn = "." + fqn

	// fetch the first-registered schema associated with the FQN of `msg`
	ps := t.sm.GetByFQName(fqn)
	if ps == nil {
		return nil, errors.Errorf("`%s`: fully-qualified name not found", fqn)
	}
	// wrap the marshaled payload within a `ProtobufPayload` message
	pp := &ProtobufPayload{SchemaUID: ps.SchemaUID, Payload: payload}

	// marshal the `ProtobufPayload` and return the result
	return proto.Marshal(pp)
}

// -----------------------------------------------------------------------------

// We need to access protobuf's internal decoding machinery: the `go:linkname`
// directive instructs the compiler to declare a local symbol as an alias
// for an external one, even if it's private.
// This allows us to bind to the private `unmarshalType` method of the
// `proto.Buffer` class, which does the actual work of decoding the payload
// based on the structure-tags of the receiving object.
//
// `unmarshalType` is actually a method of the `proto.Buffer` class, hence the
// `b` given as first parameter will be used as 'this'.
//
// Due to the way Go mangles symbol names when using vendoring, the `go:linkname`
// clause is automatically generated via *linkname-gen*[1].
// [1] https://github.com/znly/linkname-gen.
//
//go:generate linkname-gen -symbol "github.com/gogo/protobuf/proto.(*Buffer).unmarshalType" -def "func unmarshalType(*proto.Buffer, reflect.Type, *proto.StructProperties, bool, unsafe.Pointer) error"

// Decode decodes the given protein-encoded `payload` into a dynamically
// generated structure-type.
//
// It is used when you need to work with protein-encoded data in a completely
// agnostic way (e.g. when you merely know the respective names of the fields
// you're interested in, such as a generic data-enricher for example).
//
// When decoding a specific version of a schema for the first-time in the
// lifetime of a `Transcoder`, a structure-type must be created from the
// dependency tree of this schema.
// This is a costly operation that involves a lot of reflection, see
// `CreateStructType` documentation for more information.
// Fortunately, the resulting structure-type is cached so that it can be freely
// re-used by later calls to `Decode`; i.e. you pay the price only once.
//
// Also, when trying to decode a specific schema for the first-time, `Decode`
// might not have all of the dependencies directly available in its local
// `SchemaMap`, in which case it will call the user-defined `TranscoderGetter`
// in the hope that it might return these missing dependencies.
// This user-defined function may or may not do some kind of I/O; the given
// context will be passed to it.
//
// Once again, this price is paid only once.
func (t *Transcoder) Decode(
	ctx context.Context, payload []byte,
) (reflect.Value, error) {
	var pp ProtobufPayload
	if err := proto.Unmarshal(payload, &pp); err != nil {
		return reflect.ValueOf(nil), errors.WithStack(err)
	}

	// fetch structure-type from cache, or create it if it doesn't exist
	var structType reflect.Type
	var ok bool
	schemaUID := pp.GetSchemaUID()
	t.typeCacheLock.RLock()
	structType, ok = t.typeCache[schemaUID]
	t.typeCacheLock.RUnlock()
	if !ok {
		if _, err := t.GetAndUpsert(ctx, schemaUID); err != nil {
			return reflect.ValueOf(nil), errors.WithStack(err)
		}
		st, err := CreateStructType(schemaUID, t.sm)
		if err != nil {
			return reflect.ValueOf(nil), errors.WithStack(err)
		}
		if st.Kind() != reflect.Struct {
			return reflect.ValueOf(nil), errors.Errorf(
				"`%s`: not a struct type", structType,
			)
		}
		structType = st
		t.typeCacheLock.Lock()
		t.typeCache[schemaUID] = st // upsert type-cache
		t.typeCacheLock.Unlock()
	}

	// allocate a new structure using the given structure-type definition, the
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

// DecodeAs decodes the given protein-encoded `payload` into the specified
// protobuf `Message` using the standard protobuf methods, thus bypassing all
// of the runtime-decoding and schema versioning machinery.
//
// It is very often used when you need to work with protein-encoded data in
// a non-agnostic way (i.e. when you know beforehand how you want to decode
// and interpret the data).
//
// `DecodeAs` basically adds zero overhead compared to a straightforward
// `proto.Unmarshal` call.
//
// `DecodeAs` never does any kind of I/O.
func (t *Transcoder) DecodeAs(payload []byte, msg proto.Message) error {
	var ps ProtobufPayload
	if err := proto.Unmarshal(payload, &ps); err != nil {
		return errors.WithStack(err)
	}
	return proto.Unmarshal(ps.Payload, msg)
}
