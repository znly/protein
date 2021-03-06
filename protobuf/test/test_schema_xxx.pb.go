// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test_schema_xxx.proto

package test

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/types"

// skipping weak import google_protobuf1 "google/protobuf"
import _ "github.com/gogo/protobuf/gogoproto"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

type TestSchemaXXX_WeatherType int32

const (
	TestSchemaXXX_SUN   TestSchemaXXX_WeatherType = 0
	TestSchemaXXX_CLOUD TestSchemaXXX_WeatherType = 1
	TestSchemaXXX_RAIN  TestSchemaXXX_WeatherType = 2
	TestSchemaXXX_SNOW  TestSchemaXXX_WeatherType = 3
)

var TestSchemaXXX_WeatherType_name = map[int32]string{
	0: "SUN",
	1: "CLOUD",
	2: "RAIN",
	3: "SNOW",
}
var TestSchemaXXX_WeatherType_value = map[string]int32{
	"SUN":   0,
	"CLOUD": 1,
	"RAIN":  2,
	"SNOW":  3,
}

func (x TestSchemaXXX_WeatherType) String() string {
	return proto.EnumName(TestSchemaXXX_WeatherType_name, int32(x))
}
func (TestSchemaXXX_WeatherType) EnumDescriptor() ([]byte, []int) {
	return FileDescriptorTestSchemaXxx, []int{1, 0}
}

type OtherTestSchemaXXX struct {
	Ts *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=ts" json:"ts,omitempty"`
}

func (m *OtherTestSchemaXXX) Reset()                    { *m = OtherTestSchemaXXX{} }
func (m *OtherTestSchemaXXX) String() string            { return proto.CompactTextString(m) }
func (*OtherTestSchemaXXX) ProtoMessage()               {}
func (*OtherTestSchemaXXX) Descriptor() ([]byte, []int) { return FileDescriptorTestSchemaXxx, []int{0} }

func (m *OtherTestSchemaXXX) GetTs() *google_protobuf.Timestamp {
	if m != nil {
		return m.Ts
	}
	return nil
}

type TestSchemaXXX struct {
	SchemaUID string                                `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	FQNames   []string                              `protobuf:"bytes,2,rep,name=fq_name,json=fqName" json:"fq_name,omitempty"`
	Deps      map[string]*TestSchemaXXX_NestedEntry `protobuf:"bytes,4,rep,name=deps" json:"deps,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	Ids       map[int32]string                      `protobuf:"bytes,12,rep,name=ids" json:"ids,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Ts        google_protobuf.Timestamp             `protobuf:"bytes,7,opt,name=ts" json:"ts"`
	Ots       *OtherTestSchemaXXX                   `protobuf:"bytes,9,opt,name=ots" json:"ots,omitempty"`
	Nss       []TestSchemaXXX_NestedEntry           `protobuf:"bytes,8,rep,name=nss" json:"nss"`
	Weathers  []TestSchemaXXX_WeatherType           `protobuf:"varint,13,rep,packed,name=weathers,enum=test.TestSchemaXXX_WeatherType" json:"weathers,omitempty"`
	TSStd     time.Time                             `protobuf:"bytes,100,opt,name=ts_std,json=tsStd,stdtime" json:"ts_std"`
	DurStd    []*time.Duration                      `protobuf:"bytes,101,rep,name=dur_std,json=durStd,stdduration" json:"dur_std,omitempty"`
}

func (m *TestSchemaXXX) Reset()                    { *m = TestSchemaXXX{} }
func (m *TestSchemaXXX) String() string            { return proto.CompactTextString(m) }
func (*TestSchemaXXX) ProtoMessage()               {}
func (*TestSchemaXXX) Descriptor() ([]byte, []int) { return FileDescriptorTestSchemaXxx, []int{1} }

func (m *TestSchemaXXX) GetSchemaUID() string {
	if m != nil {
		return m.SchemaUID
	}
	return ""
}

func (m *TestSchemaXXX) GetFQNames() []string {
	if m != nil {
		return m.FQNames
	}
	return nil
}

func (m *TestSchemaXXX) GetDeps() map[string]*TestSchemaXXX_NestedEntry {
	if m != nil {
		return m.Deps
	}
	return nil
}

func (m *TestSchemaXXX) GetIds() map[int32]string {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *TestSchemaXXX) GetTs() google_protobuf.Timestamp {
	if m != nil {
		return m.Ts
	}
	return google_protobuf.Timestamp{}
}

func (m *TestSchemaXXX) GetOts() *OtherTestSchemaXXX {
	if m != nil {
		return m.Ots
	}
	return nil
}

func (m *TestSchemaXXX) GetNss() []TestSchemaXXX_NestedEntry {
	if m != nil {
		return m.Nss
	}
	return nil
}

func (m *TestSchemaXXX) GetWeathers() []TestSchemaXXX_WeatherType {
	if m != nil {
		return m.Weathers
	}
	return nil
}

func (m *TestSchemaXXX) GetTSStd() time.Time {
	if m != nil {
		return m.TSStd
	}
	return time.Time{}
}

func (m *TestSchemaXXX) GetDurStd() []*time.Duration {
	if m != nil {
		return m.DurStd
	}
	return nil
}

type TestSchemaXXX_NestedEntry struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *TestSchemaXXX_NestedEntry) Reset()         { *m = TestSchemaXXX_NestedEntry{} }
func (m *TestSchemaXXX_NestedEntry) String() string { return proto.CompactTextString(m) }
func (*TestSchemaXXX_NestedEntry) ProtoMessage()    {}
func (*TestSchemaXXX_NestedEntry) Descriptor() ([]byte, []int) {
	return FileDescriptorTestSchemaXxx, []int{1, 0}
}

func (m *TestSchemaXXX_NestedEntry) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *TestSchemaXXX_NestedEntry) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*OtherTestSchemaXXX)(nil), "test.OtherTestSchemaXXX")
	proto.RegisterType((*TestSchemaXXX)(nil), "test.TestSchemaXXX")
	proto.RegisterType((*TestSchemaXXX_NestedEntry)(nil), "test.TestSchemaXXX.NestedEntry")
	proto.RegisterEnum("test.TestSchemaXXX_WeatherType", TestSchemaXXX_WeatherType_name, TestSchemaXXX_WeatherType_value)
}

func init() { proto.RegisterFile("test_schema_xxx.proto", FileDescriptorTestSchemaXxx) }

var FileDescriptorTestSchemaXxx = []byte{
	// 578 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcf, 0x6f, 0x12, 0x41,
	0x14, 0xc7, 0x67, 0x59, 0x7e, 0xed, 0xc3, 0x1a, 0x9c, 0x68, 0xb2, 0x6e, 0xec, 0x2c, 0x69, 0x3c,
	0x90, 0x1e, 0xb6, 0xb5, 0xa6, 0x6a, 0xd4, 0x83, 0x22, 0x9a, 0x34, 0x31, 0x34, 0xee, 0x42, 0xca,
	0x8d, 0x6c, 0x9d, 0x81, 0x12, 0x0b, 0x4b, 0x77, 0x66, 0x15, 0x6e, 0x3d, 0x7a, 0xac, 0x37, 0x8f,
	0xfe, 0x09, 0xfe, 0x19, 0x1c, 0x3d, 0x7a, 0x42, 0x5d, 0xfe, 0x01, 0x8f, 0x1e, 0xcd, 0xcc, 0x02,
	0xa1, 0x96, 0xc8, 0xed, 0xcd, 0xbc, 0xf7, 0xf9, 0xbe, 0xef, 0x7b, 0x33, 0x70, 0x4b, 0x30, 0x2e,
	0x5a, 0xfc, 0xed, 0x09, 0xeb, 0xf9, 0xad, 0xe1, 0x70, 0xe8, 0x0c, 0xc2, 0x40, 0x04, 0x38, 0x2d,
	0xaf, 0x2d, 0xbb, 0x13, 0x04, 0x9d, 0x53, 0xb6, 0xa3, 0xee, 0x8e, 0xa3, 0xf6, 0x8e, 0xe8, 0xf6,
	0x18, 0x17, 0x7e, 0x6f, 0x90, 0x94, 0x59, 0xe4, 0xdf, 0x02, 0x1a, 0x85, 0xbe, 0xe8, 0x06, 0xfd,
	0x59, 0x7e, 0x73, 0x91, 0xe8, 0x04, 0x9d, 0x40, 0x1d, 0x54, 0x94, 0xa4, 0xb7, 0x9e, 0x01, 0x3e,
	0x14, 0x27, 0x2c, 0xac, 0x33, 0x2e, 0x3c, 0x65, 0xa1, 0xd9, 0x6c, 0xe2, 0x6d, 0x48, 0x09, 0x6e,
	0x6a, 0x25, 0xad, 0x5c, 0xd8, 0xb3, 0x9c, 0xa4, 0x83, 0x33, 0x17, 0x72, 0xea, 0x73, 0x0b, 0x6e,
	0x4a, 0xf0, 0xad, 0x4f, 0x59, 0xd8, 0xb8, 0x4c, 0xdb, 0xa0, 0x47, 0x5d, 0xaa, 0x70, 0xa3, 0xb2,
	0x11, 0x4f, 0x6c, 0x23, 0xc9, 0x35, 0x0e, 0xaa, 0xae, 0xcc, 0xe0, 0xbb, 0x90, 0x6b, 0x9f, 0xb5,
	0xfa, 0x7e, 0x8f, 0x99, 0xa9, 0x92, 0x5e, 0x36, 0x2a, 0x85, 0x78, 0x62, 0xe7, 0x5e, 0xbd, 0xa9,
	0xf9, 0x3d, 0xc6, 0xdd, 0x6c, 0xfb, 0x4c, 0x06, 0xf8, 0x1e, 0xa4, 0x29, 0x1b, 0x70, 0x33, 0x5d,
	0xd2, 0xcb, 0x85, 0xbd, 0x4d, 0x47, 0xee, 0xc3, 0xb9, 0xd4, 0xc9, 0xa9, 0xb2, 0x01, 0x7f, 0xd9,
	0x17, 0xe1, 0xc8, 0x55, 0xa5, 0xd8, 0x01, 0xbd, 0x4b, 0xb9, 0x79, 0x4d, 0x11, 0x77, 0x56, 0x11,
	0x07, 0x74, 0x06, 0xc8, 0x42, 0xbc, 0xab, 0xe6, 0xcc, 0xad, 0x9b, 0xb3, 0x92, 0x1e, 0x4f, 0x6c,
	0x24, 0xa7, 0xc5, 0xdb, 0xa0, 0x07, 0x82, 0x9b, 0x86, 0x42, 0xcc, 0xa4, 0xc3, 0xd5, 0x05, 0xba,
	0xb2, 0x08, 0x3f, 0x04, 0xbd, 0xcf, 0xb9, 0x99, 0x57, 0x6e, 0xec, 0x55, 0x6e, 0x6a, 0x8c, 0x0b,
	0x46, 0x95, 0xa1, 0x59, 0x0f, 0x49, 0xe0, 0x27, 0x90, 0xff, 0xc0, 0x7c, 0xa9, 0xca, 0xcd, 0x8d,
	0x92, 0x5e, 0xbe, 0xbe, 0x9a, 0x3e, 0x4a, 0x6a, 0xea, 0xa3, 0x01, 0x73, 0x17, 0x00, 0xae, 0x42,
	0x56, 0xf0, 0x16, 0x17, 0xd4, 0xa4, 0x6b, 0xe7, 0xba, 0x21, 0x7b, 0xc6, 0x13, 0x3b, 0x53, 0xf7,
	0x3c, 0x41, 0x2f, 0x7e, 0xd8, 0x9a, 0x9b, 0x11, 0xdc, 0x13, 0x14, 0x3f, 0x85, 0x1c, 0x8d, 0x42,
	0x25, 0xc3, 0x94, 0xff, 0xdb, 0x57, 0x64, 0xaa, 0xb3, 0x8f, 0x56, 0xc9, 0x8f, 0x27, 0xb6, 0xf6,
	0x59, 0xc2, 0x59, 0x1a, 0x85, 0x9e, 0xa0, 0xd6, 0x3e, 0x14, 0x96, 0x46, 0xc3, 0x45, 0xd0, 0xdf,
	0xb1, 0x51, 0xf2, 0x21, 0x5c, 0x19, 0xe2, 0x9b, 0x90, 0x79, 0xef, 0x9f, 0x46, 0xf2, 0xfd, 0xe5,
	0x5d, 0x72, 0xb0, 0x9a, 0x60, 0x2c, 0x5e, 0x74, 0x05, 0xb4, 0xbf, 0x0c, 0xad, 0xdf, 0xe8, 0x4c,
	0xf5, 0x71, 0xea, 0x91, 0x66, 0x3d, 0x80, 0xfc, 0xfc, 0xe5, 0x97, 0x85, 0x33, 0xff, 0x71, 0x23,
	0xb9, 0xad, 0x7d, 0x28, 0x2c, 0x6d, 0x19, 0xe7, 0x40, 0xf7, 0x1a, 0xb5, 0x22, 0xc2, 0x06, 0x64,
	0x5e, 0xbc, 0x3e, 0x6c, 0x54, 0x8b, 0x1a, 0xce, 0x43, 0xda, 0x7d, 0x7e, 0x50, 0x2b, 0xa6, 0x64,
	0xe4, 0xd5, 0x0e, 0x8f, 0x8a, 0x7a, 0x65, 0xf7, 0xf7, 0x2f, 0x82, 0xce, 0x63, 0x82, 0x3e, 0xc6,
	0x04, 0x8d, 0x63, 0x82, 0xbe, 0xc5, 0x04, 0x7d, 0x8f, 0x09, 0xfa, 0x19, 0x13, 0xf4, 0x27, 0x26,
	0xe8, 0x7c, 0x4a, 0xd0, 0xc5, 0x94, 0xa0, 0x2f, 0x53, 0x82, 0xbe, 0x4e, 0x89, 0xd6, 0xd4, 0x8e,
	0xb3, 0x6a, 0xb1, 0xf7, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xe3, 0x59, 0x90, 0x78, 0x0f, 0x04,
	0x00, 0x00,
}
