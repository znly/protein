// Code generated by protoc-gen-gogo.
// source: protobuf_schema.proto
// DO NOT EDIT!

/*
Package protoscan is a generated protocol buffer package.

It is generated from these files:
	protobuf_schema.proto
	test_schema.proto

It has these top-level messages:
	ProtobufSchema
	TestSchema
*/
package protoscan

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
import _ "github.com/gogo/protobuf/gogoproto"

import strings "strings"
import reflect "reflect"
import github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// ProtobufSchema is a versioned protobuf Message or Enum descriptor that can
// be used to decode protobuf payloads at runtime.
type ProtobufSchema struct {
	// UID is the unique, deterministic & versioned identifier for this schema.
	UID string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	// FQName is the fully-qualified name for this schema,
	// e.g. `.google.protobuf.Timestamp`.
	FQName string `protobuf:"bytes,2,opt,name=fq_name,json=fqName,proto3" json:"fq_name,omitempty"`
	// Descriptor is either a Message or an Enum protobuf descriptor.
	//
	// Types that are valid to be assigned to Descr:
	//	*ProtobufSchema_Message
	//	*ProtobufSchema_Enum
	Descr isProtobufSchema_Descr `protobuf_oneof:"descr"`
	// Deps contains every direct and indirect dependencies required by this
	// schema.
	//
	// Key: the dependency's UID
	// Value: the dependency's fully-qualified name
	Deps map[string]string `protobuf:"bytes,4,rep,name=deps" json:"deps" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *ProtobufSchema) Reset()                    { *m = ProtobufSchema{} }
func (m *ProtobufSchema) String() string            { return proto.CompactTextString(m) }
func (*ProtobufSchema) ProtoMessage()               {}
func (*ProtobufSchema) Descriptor() ([]byte, []int) { return fileDescriptorProtobufSchema, []int{0} }

type isProtobufSchema_Descr interface {
	isProtobufSchema_Descr()
}

type ProtobufSchema_Message struct {
	Message *google_protobuf.DescriptorProto `protobuf:"bytes,30,opt,name=msg,oneof"`
}
type ProtobufSchema_Enum struct {
	Enum *google_protobuf.EnumDescriptorProto `protobuf:"bytes,31,opt,name=enm,oneof"`
}

func (*ProtobufSchema_Message) isProtobufSchema_Descr() {}
func (*ProtobufSchema_Enum) isProtobufSchema_Descr()    {}

func (m *ProtobufSchema) GetDescr() isProtobufSchema_Descr {
	if m != nil {
		return m.Descr
	}
	return nil
}

func (m *ProtobufSchema) GetUID() string {
	if m != nil {
		return m.UID
	}
	return ""
}

func (m *ProtobufSchema) GetFQName() string {
	if m != nil {
		return m.FQName
	}
	return ""
}

func (m *ProtobufSchema) GetMessage() *google_protobuf.DescriptorProto {
	if x, ok := m.GetDescr().(*ProtobufSchema_Message); ok {
		return x.Message
	}
	return nil
}

func (m *ProtobufSchema) GetEnum() *google_protobuf.EnumDescriptorProto {
	if x, ok := m.GetDescr().(*ProtobufSchema_Enum); ok {
		return x.Enum
	}
	return nil
}

func (m *ProtobufSchema) GetDeps() map[string]string {
	if m != nil {
		return m.Deps
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ProtobufSchema) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ProtobufSchema_OneofMarshaler, _ProtobufSchema_OneofUnmarshaler, _ProtobufSchema_OneofSizer, []interface{}{
		(*ProtobufSchema_Message)(nil),
		(*ProtobufSchema_Enum)(nil),
	}
}

func _ProtobufSchema_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ProtobufSchema)
	// descr
	switch x := m.Descr.(type) {
	case *ProtobufSchema_Message:
		_ = b.EncodeVarint(30<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Message); err != nil {
			return err
		}
	case *ProtobufSchema_Enum:
		_ = b.EncodeVarint(31<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Enum); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ProtobufSchema.Descr has unexpected type %T", x)
	}
	return nil
}

func _ProtobufSchema_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ProtobufSchema)
	switch tag {
	case 30: // descr.msg
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.DescriptorProto)
		err := b.DecodeMessage(msg)
		m.Descr = &ProtobufSchema_Message{msg}
		return true, err
	case 31: // descr.enm
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.EnumDescriptorProto)
		err := b.DecodeMessage(msg)
		m.Descr = &ProtobufSchema_Enum{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ProtobufSchema_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ProtobufSchema)
	// descr
	switch x := m.Descr.(type) {
	case *ProtobufSchema_Message:
		s := proto.Size(x.Message)
		n += proto.SizeVarint(30<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ProtobufSchema_Enum:
		s := proto.Size(x.Enum)
		n += proto.SizeVarint(31<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ProtobufSchema)(nil), "protoscan.ProtobufSchema")
}
func (this *ProtobufSchema) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&protoscan.ProtobufSchema{")
	s = append(s, "UID: "+fmt.Sprintf("%#v", this.UID)+",\n")
	s = append(s, "FQName: "+fmt.Sprintf("%#v", this.FQName)+",\n")
	if this.Descr != nil {
		s = append(s, "Descr: "+fmt.Sprintf("%#v", this.Descr)+",\n")
	}
	keysForDeps := make([]string, 0, len(this.Deps))
	for k, _ := range this.Deps {
		keysForDeps = append(keysForDeps, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForDeps)
	mapStringForDeps := "map[string]string{"
	for _, k := range keysForDeps {
		mapStringForDeps += fmt.Sprintf("%#v: %#v,", k, this.Deps[k])
	}
	mapStringForDeps += "}"
	if this.Deps != nil {
		s = append(s, "Deps: "+mapStringForDeps+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtobufSchema_Message) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&protoscan.ProtobufSchema_Message{` +
		`Message:` + fmt.Sprintf("%#v", this.Message) + `}`}, ", ")
	return s
}
func (this *ProtobufSchema_Enum) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&protoscan.ProtobufSchema_Enum{` +
		`Enum:` + fmt.Sprintf("%#v", this.Enum) + `}`}, ", ")
	return s
}
func valueToGoStringProtobufSchema(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

func init() { proto.RegisterFile("protobuf_schema.proto", fileDescriptorProtobufSchema) }

var fileDescriptorProtobufSchema = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x51, 0xc1, 0x4e, 0xea, 0x40,
	0x14, 0x9d, 0xd2, 0x42, 0x1f, 0x43, 0xf2, 0xf2, 0x32, 0x79, 0x2f, 0xe9, 0x23, 0x71, 0xda, 0x88,
	0x0b, 0x56, 0xc5, 0xe0, 0x42, 0x63, 0xe2, 0xc2, 0x06, 0x8c, 0x2e, 0x34, 0x5a, 0xe3, 0x9a, 0x14,
	0x18, 0x2a, 0x91, 0xe9, 0x94, 0x0e, 0x35, 0x61, 0xc7, 0xd2, 0xa5, 0x9f, 0xe0, 0x27, 0xf8, 0x19,
	0x2c, 0x5d, 0xb2, 0x6a, 0x74, 0xfa, 0x03, 0x2e, 0x5d, 0x9a, 0x4e, 0xa1, 0x89, 0xba, 0xba, 0xf7,
	0xf6, 0x9c, 0x7b, 0xee, 0xe9, 0x19, 0xf8, 0x2f, 0x8c, 0xd8, 0x8c, 0xf5, 0xe3, 0x51, 0x8f, 0x0f,
	0x6e, 0x09, 0xf5, 0x6c, 0x39, 0xa3, 0xaa, 0x2c, 0x7c, 0xe0, 0x05, 0x75, 0xcb, 0x67, 0xcc, 0x9f,
	0x90, 0xd6, 0x86, 0xd8, 0x1a, 0x12, 0x3e, 0x88, 0xc6, 0xe1, 0x8c, 0x45, 0x39, 0xb9, 0xbe, 0x55,
	0x40, 0x3e, 0xf3, 0x99, 0x1c, 0x64, 0x97, 0xc3, 0xdb, 0xab, 0x12, 0xfc, 0x7d, 0xb9, 0x66, 0x5c,
	0xcb, 0x23, 0xe8, 0x3f, 0x54, 0xe3, 0xf1, 0xd0, 0x50, 0x2c, 0xa5, 0x59, 0x75, 0x74, 0x91, 0x98,
	0xea, 0xcd, 0x59, 0xc7, 0xcd, 0xbe, 0xa1, 0x06, 0xd4, 0x47, 0xd3, 0x5e, 0xe0, 0x51, 0x62, 0x94,
	0x24, 0x0c, 0x45, 0x62, 0x56, 0x4e, 0xae, 0x2e, 0x3c, 0x4a, 0xdc, 0xca, 0x68, 0x9a, 0x55, 0x74,
	0x0c, 0x55, 0xca, 0x7d, 0x03, 0x5b, 0x4a, 0xb3, 0xd6, 0xb6, 0xec, 0xdc, 0xa1, 0xbd, 0xb1, 0x61,
	0x77, 0x0a, 0x87, 0xf2, 0xae, 0x53, 0x13, 0x89, 0xa9, 0x9f, 0x13, 0xce, 0x3d, 0x9f, 0x9c, 0x02,
	0x37, 0xdb, 0x45, 0x0e, 0x54, 0x49, 0x40, 0x0d, 0x53, 0x4a, 0xec, 0xfc, 0x90, 0xe8, 0x06, 0x31,
	0xfd, 0x2e, 0xf3, 0x4b, 0x24, 0xa6, 0x96, 0x01, 0x99, 0x06, 0x09, 0x28, 0x3a, 0x82, 0xda, 0x90,
	0x84, 0xdc, 0xd0, 0x2c, 0xb5, 0x59, 0x6b, 0x37, 0xec, 0x22, 0x34, 0xfb, 0xeb, 0xff, 0xda, 0x1d,
	0x12, 0xf2, 0x6e, 0x30, 0x8b, 0xe6, 0x8e, 0xb6, 0x4c, 0x4c, 0xe0, 0xca, 0xb5, 0xfa, 0x3e, 0xac,
	0x16, 0x00, 0xfa, 0x03, 0xd5, 0x3b, 0x32, 0xcf, 0x23, 0x71, 0xb3, 0x16, 0xfd, 0x85, 0xe5, 0x7b,
	0x6f, 0x12, 0xaf, 0x73, 0x70, 0xf3, 0xe1, 0xb0, 0x74, 0xa0, 0x38, 0x3a, 0x2c, 0xcb, 0x47, 0x70,
	0x76, 0xdf, 0xdf, 0xb0, 0xb2, 0x10, 0x18, 0x3c, 0x08, 0x0c, 0x96, 0x02, 0x83, 0x17, 0x81, 0xc1,
	0x4a, 0x60, 0xf0, 0x2a, 0x30, 0xf8, 0x10, 0x18, 0x2c, 0x52, 0x0c, 0x1e, 0x53, 0x0c, 0x9e, 0x52,
	0x0c, 0x9e, 0x53, 0xac, 0xf4, 0x2b, 0xd2, 0xe3, 0xde, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0xba,
	0xa6, 0x67, 0x7c, 0xf8, 0x01, 0x00, 0x00,
}
