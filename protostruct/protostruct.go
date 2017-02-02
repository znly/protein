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
	"reflect"
	"strings"

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
		return nil, err
	}

	structTypes := make(map[string]reflect.Type, len(pss))
	// map[schemaUID][]schemaUID
	//delayedTypes := map[string][]string{}

	// for every dep, build the corresponding structure type
	for uid, ps := range pss {
		msg, ok := ps.GetDescr().(*protein.ProtobufSchema_Message)
		if !ok {
			continue
		}

		fields := msg.Message.GetField()
		structFields := make([]reflect.StructField, 0, len(fields))
		for _, f := range fields {
			// name
			fName := fieldName(f)

			// type
			var fType reflect.Type
			if depName := f.GetTypeName(); len(depName) > 0 { // message or enum
				_ = uid
			} else { // everything else
				t, err := fieldType(f)
				if err != nil {
					return nil, err
				}
				fType = t
			}

			//tag

			if f.GetTypeName() == "" {
				structFields = append(structFields, reflect.StructField{
					Name: fName,
					Type: fType,
				})
			}
		}

		structTypes[uid] = reflect.StructOf(structFields)
	}

	structType := structTypes[schemaUID]
	return &structType, nil
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
	// TODO(cmc): stars?
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
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		t = reflect.TypeOf(float64(0))
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		t = reflect.TypeOf(float32(0))
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
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		t = reflect.TypeOf(uint32(0))
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		err = errors.Wrapf(
			ErrFieldTypeNotSupported,
			"`%s`: field type not supported", f.GetType(),
		)
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		t = reflect.TypeOf(float32(0))
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		t = reflect.TypeOf(float64(0))
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		t = reflect.TypeOf(int32(0))
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		t = reflect.TypeOf(int64(0))
	}

	return t, err
}

// From https://github.com/golang/lint/blob/master/lint.go
var _commonAcronyms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"FQ":    true, // fully-qualified
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}
