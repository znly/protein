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

package wirer

import (
	"reflect"
	"unsafe"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/znly/protein"
	"github.com/znly/protein/protostruct"
)

// -----------------------------------------------------------------------------

// Versioned implements a Wirer that integrates with a Bank in order to augment
// the protobuf payloads that it encodes with additional versioning metadata.
//
// These metadata are then used by the internal decoder of the versioned Wirer
// to determinate how to decode an incoming payload on the wire.
type Versioned struct {
	b protein.Bank
}

func NewVersioned(b protein.Bank) *Versioned { return &Versioned{b: b} }

// -----------------------------------------------------------------------------

// TODO(cmc): Explain why this is necessary.
func (v *Versioned) Encode(o proto.Message) ([]byte, error) {
	return v.EncodeWithName(o, proto.MessageName(o))
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
func (v *Versioned) EncodeWithName(
	o proto.Message, fqName string,
) ([]byte, error) {
	// marshal the actual protobuf message
	payload, err := proto.Marshal(o)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// find the first UID associated with the fully-qualified name of `o`
	uids := v.b.FQNameToUID("." + fqName)
	if len(uids) <= 0 {
		return nil, errors.Errorf("`%s`: FQ-name not found in bank", fqName)
	}
	// wrap the marshaled payload within a ProtobufPayload message
	pp := &protein.ProtobufPayload{
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
// `proto.Buffer` class, which does the actual work of decoding a binary payload
// into an instance of the specified `reflect.Type`.
//
// `unmarshalType` is actually a method of the `proto.Buffer` class, hence the
// `b` given as first parameter will be used as "this".
//
//go:linkname unmarshalType github.com/gogo/protobuf/proto.(*Buffer).unmarshalType
func unmarshalType(b *proto.Buffer,
	st reflect.Type, prop *proto.StructProperties, is_group bool,
	base unsafe.Pointer, // implicitly casted to a `proto.structPointer`
) error

// DecodeStruct decodes the `payload` into a dynamically-defined structure
// type.
func (v *Versioned) DecodeStruct(payload []byte) (*reflect.Value, error) {
	var pp protein.ProtobufPayload
	if err := proto.Unmarshal(payload, &pp); err != nil {
		return nil, errors.WithStack(err)
	}
	structType, err := protostruct.CreateStructType(v.b, pp.GetUID())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if (*structType).Kind() != reflect.Struct {
		return nil, errors.Errorf("`%s`: not a struct type", structType)
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
func (v *Versioned) DecodeMessage(payload []byte, dst proto.Message) error {
	var ps protein.ProtobufPayload
	if err := proto.Unmarshal(payload, &ps); err != nil {
		return errors.WithStack(err)
	}
	return proto.Unmarshal(ps.Payload, dst)
}
