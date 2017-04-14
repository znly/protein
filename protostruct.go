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
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/pkg/errors"
	"github.com/znly/protein/failure"
)

// -----------------------------------------------------------------------------

// CreateStructType constructs a new structure-type definition at runtime from
// a `ProtobufSchema` tree.
//
// This newly-created structure embeds all the necessary tags & hints for the
// protobuf SDK to deserialize payloads into it.
// I.e., it allows for runtime-decoding of protobuf payloads.
//
// This is a complex and costly operation, it is strongly recommended to cache
// the result like the `Transcoder` does.
//
// It requires the new `reflect` APIs provided by Go 1.7+.
func CreateStructType(schemaUID string, sm *SchemaMap) (reflect.Type, error) {
	ps := sm.GetByUID(schemaUID)
	if ps == nil {
		return reflect.TypeOf(nil), errors.Wrapf(failure.ErrSchemaNotFound,
			"`%s`: no schema with this UID", schemaUID,
		)
	}
	deps := ps.GetDeps()
	pss := make(map[string]*ProtobufSchema, len(deps))
	for depUID := range deps {
		if dep := sm.GetByUID(depUID); dep == nil {
			return reflect.TypeOf(nil), errors.Wrapf(failure.ErrDependencyNotFound,
				"`%s`: no dependency with this UID", depUID,
			)
		} else {
			pss[depUID] = dep
		}
	}
	pss[schemaUID] = ps
	// map[FQName]schemaUID
	pssRevMap := make(map[string]string, len(pss))
	for uid, ps := range pss {
		pssRevMap[ps.GetFQName()] = uid
	}

	// pre-allocations are mucho importante here
	structFields := make(map[string][]reflect.StructField, len(pss))
	structTypes := make(map[string]reflect.Type, len(pss))

	mapEntryTags := make(map[string][2]reflect.StructTag)

	if err := buildScalarTypes(pss, structFields); err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO(cmc)[3]: handle gogo's custom/std types

	if err := buildCompoundTypes(
		ps, pss, pssRevMap,
		structFields, structTypes, mapEntryTags,
	); err != nil {
		return nil, errors.WithStack(err)
	}

	return structTypes[schemaUID], nil
}

func buildScalarTypes(
	pss map[string]*ProtobufSchema,
	structFields map[string][]reflect.StructField,
) error {
	for uid, ps := range pss {
		msg, ok := ps.GetDescr().(*ProtobufSchema_Message)
		if !ok {
			continue
		}

		fields := make([]reflect.StructField, 0, len(msg.Message.GetField()))
		for _, f := range msg.Message.GetField() {
			if len(f.GetTypeName()) > 0 { // message or enum
				continue
			}

			// name
			fName := fieldName(f)
			// type
			var fType reflect.Type
			t, err := scalarType(f)
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

func buildCompoundTypes(
	ps *ProtobufSchema,
	pss map[string]*ProtobufSchema,
	pssRevMap map[string]string,
	structFields map[string][]reflect.StructField,
	structTypes map[string]reflect.Type,
	mapEntryTags map[string][2]reflect.StructTag,
) error {
	for uid := range ps.GetDeps() {
		if err := buildCompoundTypes(
			pss[uid], pss, pssRevMap,
			structFields, structTypes, mapEntryTags,
		); err != nil {
			return errors.WithStack(err)
		}
	}
	msg, ok := ps.GetDescr().(*ProtobufSchema_Message)
	if !ok {
		return nil
	}

	fields := structFields[ps.GetUID()]
	for _, f := range msg.Message.GetField() {
		if len(f.GetTypeName()) <= 0 { // neither message nor enum
			continue
		}

		// name
		fName := fieldName(f)
		// type
		fType := structTypes[pssRevMap[f.GetTypeName()]]
		// tag
		fTag, err := fieldTag(msg.Message, f)
		if err != nil {
			return errors.WithStack(err)
		}
		if fType.Kind() == reflect.Map { // maps are special creatures
			fTag = fieldTagEntries(fTag, mapEntryTags[pssRevMap[f.GetTypeName()]])
		} else {
			switch f.GetLabel() {
			case descriptor.FieldDescriptorProto_LABEL_OPTIONAL:
				if gogoproto.IsNullable(f) {
					fType = reflect.PtrTo(fType)
				}
			case descriptor.FieldDescriptorProto_LABEL_REQUIRED:
				// do nothing
			case descriptor.FieldDescriptorProto_LABEL_REPEATED:
				fType = reflect.SliceOf(fType)
			default:
				err = errors.Wrapf(failure.ErrFieldLabelNotSupported,
					"`%s`: field label not supported", f.GetLabel(),
				)
			}
		}

		fields = append(fields, reflect.StructField{
			Name: fName,
			Type: fType,
			Tag:  fTag,
		})
	}

	if msg.Message.GetOptions().GetMapEntry() { // map type
		structTypes[ps.GetUID()] = reflect.MapOf(
			reflect.Type(fields[0].Type), reflect.Type(fields[1].Type),
		)
		mapEntryTags[ps.GetUID()] = [2]reflect.StructTag{
			fields[0].Tag, fields[1].Tag,
		}
	} else { // everything else
		structTypes[ps.GetUID()] = reflect.StructOf(fields)
	}

	return nil
}

// -----------------------------------------------------------------------------

func fieldName(f *descriptor.FieldDescriptorProto) string {
	if gogoproto.IsCustomName(f) { // custom names bypass everything
		return gogoproto.GetCustomName(f)
	}

	name := generator.CamelCase(f.GetName())
	parts := camelcase.Split(name)
	for i, p := range parts {
		if p = strings.ToUpper(p); _commonAcronyms[p] {
			parts[i] = p
		} else if p[len(p)-1] == 'S' && _commonAcronyms[p[:len(p)-1]] {
			// IDS -> IDs
			parts[i] = p[:len(p)-1] + "s"
		}
	}
	return strings.Join(parts, "")
}

func scalarType(f *descriptor.FieldDescriptorProto) (t reflect.Type, err error) {
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
		err = errors.Wrapf(failure.ErrFieldTypeNotSupported,
			"`%s`: not a scalar type", f.GetType(),
		)
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		err = errors.Wrapf(failure.ErrFieldTypeNotSupported,
			"`%s`: not a scalar type", f.GetType(),
		)
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		t = reflect.TypeOf([]byte{})
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		err = errors.Wrapf(failure.ErrFieldTypeNotSupported,
			"`%s`: not a scalar type", f.GetType(),
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
		// do nothing (deprecated in proto3)
	case descriptor.FieldDescriptorProto_LABEL_REQUIRED:
		// do nothing (deprecated in proto3)
	case descriptor.FieldDescriptorProto_LABEL_REPEATED:
		t = reflect.SliceOf(t)
	default:
		err = errors.Wrapf(failure.ErrFieldLabelNotSupported,
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
// `g` given as first parameter will be used as 'this'.
//
// Due to the way Go mangles symbol names when using vendoring, the `go:linkname`
// clause is automatically generated via *linkname-gen*[1].
// [1] https://github.com/znly/linkname-gen.
//
//go:generate linkname-gen -symbol "github.com/gogo/protobuf/protoc-gen-gogo/generator.(*Generator).goTag" -def "func goTag(*generator.Generator, *generator.Descriptor, *descriptor.FieldDescriptorProto, string) string"

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
		return "", errors.Wrapf(failure.ErrFieldTypeNotSupported,
			"`%s`: field type not supported", f.GetType(),
		)
	}

	g := &generator.Generator{}
	tagProto := goTag(g, &generator.Descriptor{DescriptorProto: d}, f, wt)
	tagProto = `protobuf:` + tagProto
	// the following condition has been imported from gogoprotobuf's
	// `protoc-gen-go/generator/generator.go`
	if f.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE &&
		f.GetType() != descriptor.FieldDescriptorProto_TYPE_GROUP &&
		!f.IsRepeated() {
		tagProto = tagProto[:len(tagProto)-1] + `,proto3"`
	}

	tagJSON := `json:"`
	if jn := f.GetJsonName(); len(jn) > 0 {
		jnParts := camelcase.Split(jn)
		for i, jnp := range jnParts {
			jnParts[i] = strings.ToLower(jnp)
		}
		tagJSON += strings.Join(jnParts, "_")
	}
	if gogoproto.IsNullable(f) {
		tagJSON += ",omitempty"
	}
	tagJSON += `"`

	tags := strings.Join([]string{tagProto, tagJSON}, " ")
	return reflect.StructTag(tags), nil
}

// fieldTagEntries computes the special tags required by protobuf maps.
func fieldTagEntries(
	tag reflect.StructTag, entryTags [2]reflect.StructTag,
) reflect.StructTag {
	tag = reflect.StructTag(fmt.Sprintf("%s %s %s", tag,
		strings.Replace(string(entryTags[0]), "protobuf", "protobuf_key", -1),
		strings.Replace(string(entryTags[1]), "protobuf", "protobuf_val", -1),
	))

	jsonTagFound := false
	tagParts := strings.Split(string(tag), " ")
	tagPartsClean := make([]string, 0, len(tagParts))
	for _, t := range tagParts {
		if strings.HasPrefix(t, "json:") {
			if jsonTagFound {
				continue
			}
			jsonTagFound = true
		}
		tagPartsClean = append(tagPartsClean, t)
	}

	return reflect.StructTag(strings.Join(tagPartsClean, " "))
}
