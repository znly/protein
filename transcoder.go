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
	"unsafe"

	"github.com/gogo/protobuf/proto"
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
	schemas map[string]*ProtobufSchema
	// reverse-mapping of fully-qualified names to UIDs
	revmap map[string][]string

	// TODO(cmc)
	getter func(ctx context.Context, uid string) ([]byte, error)
	setter func(ctx context.Context, uid string, data []byte) error
}

// TODO(cmc)
var TranscoderGetterNoOp = func(context.Context, string) ([]byte, error) {
	return nil, nil
}
var TranscoderSetterNoOp = func(context.Context, string, []byte) error {
	return nil
}

// TODO(cmc)
func NewTranscoder(ctx context.Context,
	getter func(ctx context.Context, uid string) ([]byte, error),
	setter func(ctx context.Context, uid string, data []byte) error,
) (*Transcoder, error) {
	schemas, err := ScanSchemas()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if setter == nil || getter == nil {
		// TODO(cmc): real error
		return nil, errors.New("getter and/or setter cannot be nil")
	}

	t := &Transcoder{
		schemas: schemas,
		revmap:  map[string][]string{},

		getter: getter, setter: setter,
	}
	pss := make([]*ProtobufSchema, 0, len(schemas))
	for _, ps := range schemas {
		pss = append(pss, ps)
	}
	if err := t.set(ctx, pss...); err != nil {
		return nil, errors.WithStack(err)
	}

	return t, nil
}

// -----------------------------------------------------------------------------

// TODO(cmc)
//
// get retrieves the ProtobufSchema associated with the specified identifier,
// plus all of its direct & indirect dependencies.
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
	if s, ok := t.schemas[uid]; ok { // try the local cache first..
		schemas[uid] = s
	} else { // ..then fallback on user-defined getter
		b, err := t.getter(ctx, uid)
		if err != nil {
			return nil, errors.Wrapf(err, "`%s`: schema not found", uid)
		}
		var root ProtobufSchema
		if err := proto.Unmarshal(b, &root); err != nil {
			return nil, errors.Wrapf(err, "`%s`: invalid schema", uid)
		}
		schemas[uid] = &root
		t.schemas[uid] = &root // upsert local cache just in case
	}

	// get dependencies
	deps := schemas[uid].GetDeps()

	// try the local cache first..
	psNotFound := make(map[string]struct{}, len(deps))
	for depUID := range deps {
		if s, ok := t.schemas[depUID]; ok {
			schemas[depUID] = s
			continue
		}
		psNotFound[depUID] = struct{}{}
	}
	if len(psNotFound) <= 0 { // found everything needed in local cache!
		return schemas, nil
	}

	// ..then fallback on user-defined getter
	var err error
	for depUID := range psNotFound {
		var b []byte
		b, err = t.getter(ctx, depUID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		delete(psNotFound, depUID) // it's been found!
		var ps ProtobufSchema
		if err := proto.Unmarshal(b, &ps); err != nil {
			return nil, errors.Wrapf(err, "`%s`: invalid schema (dependency)", depUID)
		}
		schemas[depUID] = &ps
		t.schemas[depUID] = &ps // upsert local cache just in case
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

// TODO(cmc)
//
// set synchronously adds the specified ProtobufSchemas to the local local
// cache; then pushes them to the underlying tuyau client's pipe.
// Whether this push is synchronous or not depends on the implementation
// of the tuyau.Client used.
//
// Put doesn't care about pre-existing keys: if a schema with the same key
// already exist, it will be overwritten; both in the local cache as well in the
// TuyauDB store.
func (t *Transcoder) set(
	ctx context.Context, pss ...*ProtobufSchema,
) error {
	// serialization + local cache & reverse-mapping
	var b []byte
	var err error
	for _, ps := range pss {
		select {
		case <-ctx.Done():
			return errors.WithStack(ctx.Err())
		default:
		}
		b, err = proto.Marshal(ps)
		if err != nil {
			return errors.WithStack(err)
		}
		uid := ps.GetUID()
		t.schemas[uid] = ps
		t.revmap[ps.GetFQName()] = append(t.revmap[ps.GetFQName()], uid)
		if err := t.setter(ctx, uid, b); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

// FQNameToUID returns the UID associated with the given fully-qualified name.
//
// It is possible that multiple versions of a schema identified by a FQ name
// are currently available in the local cache; in which case all of the
// associated UIDs will be returned to the caller, *in random order*.
//
// The reverse-mapping is pre-computed; don't hesitate to call this method, it'll
// be real fast.
//
// It returns nil if `fqName` doesn't match any schema in the bank.
func (t *Transcoder) FQNameToUID(fqName string) []string { return t.revmap[fqName] }

// TODO(cmc)
func (t *Transcoder) UIDToFQName(uid string) string {
	if s, ok := t.schemas[uid]; ok {
		return s.GetFQName()
	}
	return ""
}

// -----------------------------------------------------------------------------

// TODO(cmc): Explain why this is necessary.
func (t *Transcoder) Encode(o proto.Message) ([]byte, error) {
	return t.EncodeWithName(o, proto.MessageName(o))
}

// EncodeWithName marshals the given protobuf message then wraps it up into a
// ProtobufPayload object that adds additional versioning metadata.
//
// Encode uses the message's fully-qualified name to reverse-lookup its UID.
// Note that a single FQ-name might point to multiple UIDs if multiple versions
// of the associated message are currently availaible in the bank.
// When this happens, the first UID from the returned list will be used (and
// since this list is randomly-ordered, effectively a random UID will be used).
//
// TODO(cmc): explain this mess.
func (t *Transcoder) EncodeWithName(
	o proto.Message, fqName string,
) ([]byte, error) {
	// marshal the actual protobuf message
	payload, err := proto.Marshal(o)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// find the first UID associated with the fully-qualified name of `o`
	uids := t.FQNameToUID("." + fqName)
	if len(uids) <= 0 {
		return nil, errors.Errorf("`%s`: FQ-name not found in bank", fqName)
	}
	// wrap the marshaled payload within a ProtobufPayload message
	pp := &ProtobufPayload{
		UID:     uids[0],
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
func (t *Transcoder) Decode(payload []byte) (*reflect.Value, error) {
	var pp ProtobufPayload
	if err := proto.Unmarshal(payload, &pp); err != nil {
		return nil, errors.WithStack(err)
	}
	var structType *reflect.Type
	var err error
	//structType, err := protostruct.CreateStructType(t.b, pp.GetUID())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if (*structType).Kind() != reflect.Struct {
		return nil, errors.Errorf("`%s`: not a struct type", *structType)
	}

	// allocate a new structure using the given type definition, the
	// returned `reflect.Value`'s underlying type is a pointer to struct
	obj := reflect.New(*structType)

	b := proto.NewBuffer(pp.GetPayload())
	unmarshalType(b,
		// the structure definition, computed at runtime
		*structType,
		// the protobuf properties of the struct, computed via its struct tags
		proto.GetProperties(*structType),
		// is_group, deprecated
		false,
		// the address we want to deserialize to
		unsafe.Pointer(obj.Elem().Addr().Pointer()),
	)

	return &obj, nil
}

// -----------------------------------------------------------------------------

// TODO(cmc): doc & test
func (t *Transcoder) DecodeAs(payload []byte, dst proto.Message) error {
	var ps ProtobufPayload
	if err := proto.Unmarshal(payload, &ps); err != nil {
		return errors.WithStack(err)
	}
	return proto.Unmarshal(ps.Payload, dst)
}
