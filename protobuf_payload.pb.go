// Code generated by protoc-gen-gogo.
// source: protobuf_payload.proto
// DO NOT EDIT!

package protein

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strings "strings"
import reflect "reflect"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// ProtobufPayload is a protobuf payload annotated with the unique versioning
// identifier of its schema.
//
// This allows a `ProtobufPayload` to be decoded at runtime using Protein's
// `Transcoder`.
//
// See `ScanSchemas`'s documentation for more information.
type ProtobufPayload struct {
	// SchemaUID is the unique, deterministic & versioned identifier of the
	// `payload`'s schema.
	SchemaUID string `protobuf:"bytes,1,opt,name=schema_uid,json=schemaUid,proto3" json:"schema_uid,omitempty"`
	// Payload is the actual, marshaled protobuf payload.
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *ProtobufPayload) Reset()                    { *m = ProtobufPayload{} }
func (m *ProtobufPayload) String() string            { return proto.CompactTextString(m) }
func (*ProtobufPayload) ProtoMessage()               {}
func (*ProtobufPayload) Descriptor() ([]byte, []int) { return fileDescriptorProtobufPayload, []int{0} }

func (m *ProtobufPayload) GetSchemaUID() string {
	if m != nil {
		return m.SchemaUID
	}
	return ""
}

func (m *ProtobufPayload) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*ProtobufPayload)(nil), "protein.ProtobufPayload")
}
func (this *ProtobufPayload) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&protein.ProtobufPayload{")
	s = append(s, "SchemaUID: "+fmt.Sprintf("%#v", this.SchemaUID)+",\n")
	s = append(s, "Payload: "+fmt.Sprintf("%#v", this.Payload)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringProtobufPayload(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

func init() { proto.RegisterFile("protobuf_payload.proto", fileDescriptorProtobufPayload) }

var fileDescriptorProtobufPayload = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0x2a, 0x4d, 0x8b, 0x2f, 0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0xd1, 0x03, 0x0b, 0x08,
	0xb1, 0x83, 0xa8, 0xd4, 0xcc, 0x3c, 0x29, 0x59, 0x98, 0x02, 0xfd, 0xf4, 0xfc, 0xf4, 0x7c, 0x30,
	0x07, 0xcc, 0x82, 0xa8, 0x53, 0x4a, 0xe3, 0xe2, 0x0f, 0x80, 0x2a, 0x08, 0x80, 0x18, 0x20, 0xa4,
	0xc3, 0xc5, 0x55, 0x9c, 0x9c, 0x91, 0x9a, 0x9b, 0x18, 0x5f, 0x9a, 0x99, 0x22, 0xc1, 0xa8, 0xc0,
	0xa8, 0xc1, 0xe9, 0xc4, 0xfb, 0xe8, 0x9e, 0x3c, 0x67, 0x30, 0x58, 0x34, 0xd4, 0xd3, 0x25, 0x88,
	0x13, 0xa2, 0x20, 0x34, 0x33, 0x45, 0x48, 0x95, 0x8b, 0x1d, 0x6a, 0xb3, 0x04, 0x93, 0x02, 0xa3,
	0x06, 0x8f, 0x13, 0xf7, 0xa3, 0x7b, 0xf2, 0xec, 0x50, 0xb3, 0x82, 0x60, 0x72, 0x4e, 0x06, 0x1f,
	0x1e, 0xca, 0x31, 0x36, 0x3c, 0x92, 0x63, 0xe8, 0x78, 0x24, 0xc7, 0x70, 0xe2, 0x91, 0x1c, 0xc3,
	0x85, 0x47, 0x72, 0x0c, 0x37, 0x1e, 0xc9, 0x31, 0x3c, 0x78, 0x24, 0xc7, 0xf0, 0xe3, 0x91, 0x1c,
	0x43, 0xc3, 0x63, 0x39, 0x86, 0x09, 0x8f, 0xe5, 0x18, 0x16, 0x3c, 0x96, 0x63, 0xd8, 0xf0, 0x58,
	0x8e, 0x31, 0x89, 0x0d, 0xec, 0x40, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x37, 0xbb, 0x8d,
	0x62, 0xe2, 0x00, 0x00, 0x00,
}
