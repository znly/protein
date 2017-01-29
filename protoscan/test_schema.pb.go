// Code generated by protoc-gen-gogo.
// source: test_schema.proto
// DO NOT EDIT!

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

type TestSchema struct {
	UID    string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	FQName string `protobuf:"bytes,2,opt,name=fq_name,json=fqName,proto3" json:"fq_name,omitempty"`
	// Types that are valid to be assigned to Descr:
	//	*TestSchema_Message
	//	*TestSchema_Enum
	Descr isTestSchema_Descr `protobuf_oneof:"descr"`
	Deps  map[string]string  `protobuf:"bytes,4,rep,name=deps" json:"deps" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *TestSchema) Reset()                    { *m = TestSchema{} }
func (m *TestSchema) String() string            { return proto.CompactTextString(m) }
func (*TestSchema) ProtoMessage()               {}
func (*TestSchema) Descriptor() ([]byte, []int) { return fileDescriptorTestSchema, []int{0} }

type isTestSchema_Descr interface {
	isTestSchema_Descr()
}

type TestSchema_Message struct {
	Message *google_protobuf.DescriptorProto `protobuf:"bytes,30,opt,name=msg,oneof"`
}
type TestSchema_Enum struct {
	Enum *google_protobuf.EnumDescriptorProto `protobuf:"bytes,31,opt,name=enm,oneof"`
}

func (*TestSchema_Message) isTestSchema_Descr() {}
func (*TestSchema_Enum) isTestSchema_Descr()    {}

func (m *TestSchema) GetDescr() isTestSchema_Descr {
	if m != nil {
		return m.Descr
	}
	return nil
}

func (m *TestSchema) GetUID() string {
	if m != nil {
		return m.UID
	}
	return ""
}

func (m *TestSchema) GetFQName() string {
	if m != nil {
		return m.FQName
	}
	return ""
}

func (m *TestSchema) GetMessage() *google_protobuf.DescriptorProto {
	if x, ok := m.GetDescr().(*TestSchema_Message); ok {
		return x.Message
	}
	return nil
}

func (m *TestSchema) GetEnum() *google_protobuf.EnumDescriptorProto {
	if x, ok := m.GetDescr().(*TestSchema_Enum); ok {
		return x.Enum
	}
	return nil
}

func (m *TestSchema) GetDeps() map[string]string {
	if m != nil {
		return m.Deps
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TestSchema) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TestSchema_OneofMarshaler, _TestSchema_OneofUnmarshaler, _TestSchema_OneofSizer, []interface{}{
		(*TestSchema_Message)(nil),
		(*TestSchema_Enum)(nil),
	}
}

func _TestSchema_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TestSchema)
	// descr
	switch x := m.Descr.(type) {
	case *TestSchema_Message:
		_ = b.EncodeVarint(30<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Message); err != nil {
			return err
		}
	case *TestSchema_Enum:
		_ = b.EncodeVarint(31<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Enum); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TestSchema.Descr has unexpected type %T", x)
	}
	return nil
}

func _TestSchema_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TestSchema)
	switch tag {
	case 30: // descr.msg
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.DescriptorProto)
		err := b.DecodeMessage(msg)
		m.Descr = &TestSchema_Message{msg}
		return true, err
	case 31: // descr.enm
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.EnumDescriptorProto)
		err := b.DecodeMessage(msg)
		m.Descr = &TestSchema_Enum{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TestSchema_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TestSchema)
	// descr
	switch x := m.Descr.(type) {
	case *TestSchema_Message:
		s := proto.Size(x.Message)
		n += proto.SizeVarint(30<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TestSchema_Enum:
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
	proto.RegisterType((*TestSchema)(nil), "protoscan.TestSchema")
}
func (this *TestSchema) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&protoscan.TestSchema{")
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
func (this *TestSchema_Message) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&protoscan.TestSchema_Message{` +
		`Message:` + fmt.Sprintf("%#v", this.Message) + `}`}, ", ")
	return s
}
func (this *TestSchema_Enum) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&protoscan.TestSchema_Enum{` +
		`Enum:` + fmt.Sprintf("%#v", this.Enum) + `}`}, ", ")
	return s
}
func valueToGoStringTestSchema(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

func init() { proto.RegisterFile("test_schema.proto", fileDescriptorTestSchema) }

var fileDescriptorTestSchema = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x50, 0xbb, 0x4e, 0xeb, 0x40,
	0x10, 0x5d, 0xc7, 0x4e, 0x7c, 0xb3, 0x69, 0xee, 0x5d, 0xdd, 0xc2, 0x44, 0x62, 0x6d, 0x01, 0x45,
	0x2a, 0x07, 0x85, 0x82, 0x47, 0x87, 0x95, 0x20, 0x28, 0x40, 0x60, 0xa0, 0x8e, 0x9c, 0x64, 0x62,
	0x22, 0xe2, 0x47, 0xbc, 0x36, 0x52, 0xba, 0x94, 0x94, 0x7c, 0x02, 0x9f, 0xc0, 0x67, 0xa4, 0x4c,
	0x49, 0x65, 0xc1, 0xfa, 0x07, 0x28, 0x29, 0x91, 0xd7, 0x60, 0x24, 0xa8, 0x66, 0xe6, 0x9c, 0x9d,
	0xb3, 0x67, 0x0e, 0xfe, 0x17, 0x03, 0x8b, 0xfb, 0x6c, 0x78, 0x03, 0x9e, 0x63, 0x86, 0x51, 0x10,
	0x07, 0xa4, 0x2e, 0x0a, 0x1b, 0x3a, 0x7e, 0xd3, 0x70, 0x83, 0xc0, 0x9d, 0x42, 0x5b, 0x20, 0x83,
	0x64, 0xdc, 0x1e, 0x01, 0x1b, 0x46, 0x93, 0x30, 0x0e, 0xa2, 0xe2, 0x71, 0x73, 0xbd, 0xa4, 0xdc,
	0xc0, 0x0d, 0xc4, 0x20, 0xba, 0x82, 0xde, 0x58, 0x55, 0x30, 0xbe, 0x02, 0x16, 0x5f, 0x8a, 0x0f,
	0xc8, 0x1a, 0x96, 0x93, 0xc9, 0x48, 0x93, 0x0c, 0xa9, 0x55, 0xb7, 0x54, 0x9e, 0xea, 0xf2, 0xf5,
	0x49, 0xd7, 0xce, 0x31, 0xb2, 0x89, 0xd5, 0xf1, 0xac, 0xef, 0x3b, 0x1e, 0x68, 0x15, 0x41, 0x63,
	0x9e, 0xea, 0xb5, 0xa3, 0x8b, 0x33, 0xc7, 0x03, 0xbb, 0x36, 0x9e, 0xe5, 0x95, 0x1c, 0x62, 0xd9,
	0x63, 0xae, 0x46, 0x0d, 0xa9, 0xd5, 0xe8, 0x18, 0x66, 0xe1, 0xce, 0xfc, 0xb2, 0x60, 0x76, 0x4b,
	0x77, 0xe7, 0x39, 0x64, 0x35, 0x78, 0xaa, 0xab, 0xa7, 0xc0, 0x98, 0xe3, 0xc2, 0x31, 0xb2, 0xf3,
	0x5d, 0x62, 0x61, 0x19, 0x7c, 0x4f, 0xd3, 0x85, 0xc4, 0xd6, 0x2f, 0x89, 0x9e, 0x9f, 0x78, 0x3f,
	0x65, 0xfe, 0xf0, 0x54, 0x57, 0x72, 0x22, 0xd7, 0x00, 0xdf, 0x23, 0xfb, 0x58, 0x19, 0x41, 0xc8,
	0x34, 0xc5, 0x90, 0x5b, 0x8d, 0x8e, 0x6e, 0x96, 0x81, 0x99, 0xdf, 0xb7, 0x9a, 0x5d, 0x08, 0x59,
	0xcf, 0x8f, 0xa3, 0xb9, 0xa5, 0x2c, 0x53, 0x1d, 0xd9, 0x62, 0xa5, 0xb9, 0x8b, 0xeb, 0x25, 0x41,
	0xfe, 0x62, 0xf9, 0x16, 0xe6, 0x45, 0x1c, 0x76, 0xde, 0x92, 0xff, 0xb8, 0x7a, 0xe7, 0x4c, 0x93,
	0xcf, 0x0c, 0xec, 0x62, 0x38, 0xa8, 0xec, 0x49, 0x96, 0x8a, 0xab, 0x22, 0x7c, 0x6b, 0xfb, 0xed,
	0x95, 0x4a, 0x0b, 0x4e, 0xd1, 0x3d, 0xa7, 0x68, 0xc9, 0x29, 0x5a, 0x71, 0x8a, 0x9e, 0x39, 0x45,
	0x2f, 0x9c, 0xa2, 0x77, 0x4e, 0xd1, 0x22, 0xa3, 0xe8, 0x21, 0xa3, 0xe8, 0x31, 0xa3, 0xe8, 0x29,
	0xa3, 0xd2, 0xa0, 0x26, 0xfc, 0xed, 0x7c, 0x04, 0x00, 0x00, 0xff, 0xff, 0x8c, 0xa1, 0x5f, 0xf7,
	0xec, 0x01, 0x00, 0x00,
}
