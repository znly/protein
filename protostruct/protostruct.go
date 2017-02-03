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

package protostruct

import (
	"fmt"
	"reflect"
	"strings"

	_ "unsafe" // go:linkname

	"github.com/fatih/camelcase"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/pkg/errors"
	"github.com/znly/protein"
	protein_bank "github.com/znly/protein/bank"
)

// -----------------------------------------------------------------------------

func CreateStructType(
	b protein_bank.Bank, schemaUID string,
) (*reflect.Type, error) {
	pss, err := b.Get(schemaUID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// map[FQName]schemaUID
	pssRevMap := make(map[string]string, len(pss))
	for uid, ps := range pss {
		pssRevMap[ps.GetFQName()] = uid
	}

	// pre-allocations are mucho importante here
	structFields := make(map[string][]reflect.StructField, len(pss))
	structTypes := make(map[string]reflect.Type, len(pss))

	if err := buildScalarTypes(pss, structFields); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := buildCustomTypes(
		pss[schemaUID], pss, pssRevMap, structFields, structTypes,
	); err != nil {
		return nil, errors.WithStack(err)
	}

	structType := structTypes[schemaUID]
	return &structType, nil
}

func buildScalarTypes(
	pss map[string]*protein.ProtobufSchema,
	structFields map[string][]reflect.StructField,
) error {
	for uid, ps := range pss {
		msg, ok := ps.GetDescr().(*protein.ProtobufSchema_Message)
		if !ok {
			continue
		}

		fields := make([]reflect.StructField, 0, len(msg.Message.GetField()))
		for _, f := range msg.Message.GetField() {
			if len(f.GetTypeName()) > 0 { // message or enum
				continue
			}

			fmt.Println(f)
			// name
			fName := fieldName(f)
			// type
			var fType reflect.Type
			t, err := fieldType(f)
			if err != nil {
				return errors.WithStack(err)
			}
			fType = t
			// tag
			fTag, err := fieldTag(msg.Message, f)
			if err != nil {
				return errors.WithStack(err)
			}

			fields = append(fields, reflect.StructField{
				Name: fName,
				Type: fType,
				Tag:  fTag,
			})
		}
		structFields[uid] = fields
	}

	return nil
}

func buildCustomTypes(
	ps *protein.ProtobufSchema,
	pss map[string]*protein.ProtobufSchema,
	pssRevMap map[string]string,
	structFields map[string][]reflect.StructField,
	structTypes map[string]reflect.Type,
) error {
	for uid := range ps.GetDeps() {
		if err := buildCustomTypes(
			pss[uid], pss, pssRevMap, structFields, structTypes,
		); err != nil {
			return errors.WithStack(err)
		}
	}
	msg, ok := ps.GetDescr().(*protein.ProtobufSchema_Message)
	if !ok {
		return nil
	}

	fields := structFields[ps.GetUID()]
	for _, f := range msg.Message.GetField() {
		if len(f.GetTypeName()) <= 0 { // neither message nor enum
			continue
		}

		fmt.Println(f)
		// name
		fName := fieldName(f)
		// type
		fType := structTypes[pssRevMap[f.GetTypeName()]]
		// tag
		fTag, err := fieldTag(msg.Message, f)
		if err != nil {
			return errors.WithStack(err)
		}

		fields = append(fields, reflect.StructField{
			Name: fName,
			Type: fType,
			Tag:  fTag,
		})
	}

	structTypes[ps.GetUID()] = reflect.StructOf(fields)
	return nil
}

// -----------------------------------------------------------------------------

func findSchemaUID(fqName string, pss map[string]*protein.ProtobufSchema) string {
	for _, ps := range pss {
		if fqName == ps.GetFQName() {
			return ps.GetUID()
		}
	}
	return ""
}

func fieldName(f *descriptor.FieldDescriptorProto) string {
	name := generator.CamelCase(f.GetName())
	parts := camelcase.Split(name)
	for i, p := range parts {
		if p = strings.ToUpper(p); _commonAcronyms[p] {
			parts[i] = p
		}
	}
	return strings.Join(parts, "")
}

func fieldType(f *descriptor.FieldDescriptorProto) (t reflect.Type, err error) {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		t = reflect.TypeOf(float64(0))
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		t = reflect.TypeOf(float32(0))
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		t = reflect.TypeOf(int64(0))
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		t = reflect.TypeOf(uint64(0))
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		t = reflect.TypeOf(int32(0))
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		t = reflect.TypeOf(uint32(0))
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		t = reflect.TypeOf(uint64(0))
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		t = reflect.TypeOf(uint32(0))
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		t = reflect.TypeOf(bool(false))
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		t = reflect.TypeOf(string(""))
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		err = errors.Wrapf(
			ErrFieldTypeNotSupported,
			"`%s`: field type not supported", f.GetType(),
		)
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		err = errors.Wrapf(
			ErrFieldTypeNotSupported,
			"`%s`: field type not supported", f.GetType(),
		)
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		t = reflect.TypeOf([]byte{})
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		err = errors.Wrapf(
			ErrFieldTypeNotSupported,
			"`%s`: field type not supported", f.GetType(),
		)
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		t = reflect.TypeOf(int32(0))
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		t = reflect.TypeOf(int64(0))
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		t = reflect.TypeOf(int32(0))
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		t = reflect.TypeOf(int64(0))
	}

	switch f.GetLabel() {
	case descriptor.FieldDescriptorProto_LABEL_OPTIONAL:
		// do nothing?
	case descriptor.FieldDescriptorProto_LABEL_REQUIRED:
		// do nothing?
	case descriptor.FieldDescriptorProto_LABEL_REPEATED:
		t = reflect.SliceOf(t)
	default:
		err = errors.Wrapf(
			ErrFieldLabelNotSupported,
			"`%s`: field label not supported", f.GetLabel(),
		)
	}

	return t, err
}

// -----------------------------------------------------------------------------

// We need to access protobuf's internal decoding machinery: the `go:linkname`
// directive instructs the compiler to declare a local symbol as an alias
// for an external one, even if it's private.
// This allows us to bind to the private `goTag` method of the
// `generator.Generator` class, which does the actual work of computing the
// necessary struct tags for a given protobuf field.
//
// `goTag` is actually a method of the `generator.Generator` class, hence the
// `g` given as first parameter will be used as "this".
//
//go:linkname goTag github.com/gogo/protobuf/protoc-gen-gogo/generator.(*Generator).goTag
func goTag(b *generator.Generator,
	d *generator.Descriptor, f *descriptor.FieldDescriptorProto, wt string,
) string

func fieldTag(
	d *descriptor.DescriptorProto, f *descriptor.FieldDescriptorProto,
) (reflect.StructTag, error) {
	var wt string
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		wt = "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		wt = "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		wt = "varint"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		wt = "varint"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		wt = "varint"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		wt = "varint"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		wt = "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		wt = "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		wt = "varint"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		wt = "bytes"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		wt = "group"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		wt = "bytes"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		wt = "bytes"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		wt = "varint"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		wt = "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		wt = "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		wt = "zigzag32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		wt = "zigzag64"
	default:
		return "", errors.Wrapf(
			ErrFieldTypeNotSupported,
			"`%s`: field type not supported", f.GetType(),
		)
	}

	g := &generator.Generator{}
	tag := goTag(g, &generator.Descriptor{DescriptorProto: d}, f, wt)

	return reflect.StructTag(`protobuf:"` + tag + `"`), nil
}
