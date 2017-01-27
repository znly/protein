// Code generated by protoc-gen-gogo.
// source: protobuf_schema.proto
// DO NOT EDIT!

/*
	Package protoscan is a generated protocol buffer package.

	It is generated from these files:
		protobuf_schema.proto

	It has these top-level messages:
		ProtobufSchema
*/
package protoscan

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "google/protobuf"
import _ "protobuf/gogoproto"

import strings "strings"
import reflect "reflect"
import github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"

import io "io"

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
	// uid is the unique, deterministic & versioned identifier for this schema.
	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	// fq_name is the fully-qualified name for this schema,
	// e.g. `.google.protobuf.Timestamp`.
	FqName string `protobuf:"bytes,2,opt,name=fq_name,json=fqName,proto3" json:"fq_name,omitempty"`
	// descriptor is either a Message or an Enum protobuf descriptor.
	//
	// Types that are valid to be assigned to Descriptor_:
	//	*ProtobufSchema_Msg
	//	*ProtobufSchema_Enm
	Descriptor_ isProtobufSchema_Descriptor_ `protobuf_oneof:"descriptor"`
	// deps contains every direct and indirect dependencies required by this
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

type isProtobufSchema_Descriptor_ interface {
	isProtobufSchema_Descriptor_()
	MarshalTo([]byte) (int, error)
	Size() int
	ProtoSize() int
}

type ProtobufSchema_Msg struct {
	Msg *google_protobuf.DescriptorProto `protobuf:"bytes,30,opt,name=msg,oneof"`
}
type ProtobufSchema_Enm struct {
	Enm *google_protobuf.EnumDescriptorProto `protobuf:"bytes,31,opt,name=enm,oneof"`
}

func (*ProtobufSchema_Msg) isProtobufSchema_Descriptor_() {}
func (*ProtobufSchema_Enm) isProtobufSchema_Descriptor_() {}

func (m *ProtobufSchema) GetDescriptor_() isProtobufSchema_Descriptor_ {
	if m != nil {
		return m.Descriptor_
	}
	return nil
}

func (m *ProtobufSchema) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *ProtobufSchema) GetFqName() string {
	if m != nil {
		return m.FqName
	}
	return ""
}

func (m *ProtobufSchema) GetMsg() *google_protobuf.DescriptorProto {
	if x, ok := m.GetDescriptor_().(*ProtobufSchema_Msg); ok {
		return x.Msg
	}
	return nil
}

func (m *ProtobufSchema) GetEnm() *google_protobuf.EnumDescriptorProto {
	if x, ok := m.GetDescriptor_().(*ProtobufSchema_Enm); ok {
		return x.Enm
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
		(*ProtobufSchema_Msg)(nil),
		(*ProtobufSchema_Enm)(nil),
	}
}

func _ProtobufSchema_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ProtobufSchema)
	// descriptor
	switch x := m.Descriptor_.(type) {
	case *ProtobufSchema_Msg:
		_ = b.EncodeVarint(30<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Msg); err != nil {
			return err
		}
	case *ProtobufSchema_Enm:
		_ = b.EncodeVarint(31<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Enm); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ProtobufSchema.Descriptor_ has unexpected type %T", x)
	}
	return nil
}

func _ProtobufSchema_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ProtobufSchema)
	switch tag {
	case 30: // descriptor.msg
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.DescriptorProto)
		err := b.DecodeMessage(msg)
		m.Descriptor_ = &ProtobufSchema_Msg{msg}
		return true, err
	case 31: // descriptor.enm
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.EnumDescriptorProto)
		err := b.DecodeMessage(msg)
		m.Descriptor_ = &ProtobufSchema_Enm{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ProtobufSchema_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ProtobufSchema)
	// descriptor
	switch x := m.Descriptor_.(type) {
	case *ProtobufSchema_Msg:
		s := proto.Size(x.Msg)
		n += proto.SizeVarint(30<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ProtobufSchema_Enm:
		s := proto.Size(x.Enm)
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
	s = append(s, "Uid: "+fmt.Sprintf("%#v", this.Uid)+",\n")
	s = append(s, "FqName: "+fmt.Sprintf("%#v", this.FqName)+",\n")
	if this.Descriptor_ != nil {
		s = append(s, "Descriptor_: "+fmt.Sprintf("%#v", this.Descriptor_)+",\n")
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
func (this *ProtobufSchema_Msg) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&protoscan.ProtobufSchema_Msg{` +
		`Msg:` + fmt.Sprintf("%#v", this.Msg) + `}`}, ", ")
	return s
}
func (this *ProtobufSchema_Enm) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&protoscan.ProtobufSchema_Enm{` +
		`Enm:` + fmt.Sprintf("%#v", this.Enm) + `}`}, ", ")
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
func (m *ProtobufSchema) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtobufSchema) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Uid) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintProtobufSchema(dAtA, i, uint64(len(m.Uid)))
		i += copy(dAtA[i:], m.Uid)
	}
	if len(m.FqName) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintProtobufSchema(dAtA, i, uint64(len(m.FqName)))
		i += copy(dAtA[i:], m.FqName)
	}
	if len(m.Deps) > 0 {
		for k, _ := range m.Deps {
			dAtA[i] = 0x22
			i++
			v := m.Deps[k]
			mapSize := 1 + len(k) + sovProtobufSchema(uint64(len(k))) + 1 + len(v) + sovProtobufSchema(uint64(len(v)))
			i = encodeVarintProtobufSchema(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintProtobufSchema(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintProtobufSchema(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if m.Descriptor_ != nil {
		nn1, err := m.Descriptor_.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	return i, nil
}

func (m *ProtobufSchema_Msg) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Msg != nil {
		dAtA[i] = 0xf2
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintProtobufSchema(dAtA, i, uint64(m.Msg.ProtoSize()))
		n2, err := m.Msg.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func (m *ProtobufSchema_Enm) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Enm != nil {
		dAtA[i] = 0xfa
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintProtobufSchema(dAtA, i, uint64(m.Enm.ProtoSize()))
		n3, err := m.Enm.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}
func encodeFixed64ProtobufSchema(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32ProtobufSchema(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintProtobufSchema(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ProtobufSchema) Size() (n int) {
	var l int
	_ = l
	l = len(m.Uid)
	if l > 0 {
		n += 1 + l + sovProtobufSchema(uint64(l))
	}
	l = len(m.FqName)
	if l > 0 {
		n += 1 + l + sovProtobufSchema(uint64(l))
	}
	if len(m.Deps) > 0 {
		for k, v := range m.Deps {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovProtobufSchema(uint64(len(k))) + 1 + len(v) + sovProtobufSchema(uint64(len(v)))
			n += mapEntrySize + 1 + sovProtobufSchema(uint64(mapEntrySize))
		}
	}
	if m.Descriptor_ != nil {
		n += m.Descriptor_.Size()
	}
	return n
}

func (m *ProtobufSchema_Msg) Size() (n int) {
	var l int
	_ = l
	if m.Msg != nil {
		l = m.Msg.Size()
		n += 2 + l + sovProtobufSchema(uint64(l))
	}
	return n
}
func (m *ProtobufSchema_Enm) Size() (n int) {
	var l int
	_ = l
	if m.Enm != nil {
		l = m.Enm.Size()
		n += 2 + l + sovProtobufSchema(uint64(l))
	}
	return n
}

func sovProtobufSchema(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozProtobufSchema(x uint64) (n int) {
	return sovProtobufSchema(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProtobufSchema) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtobufSchema
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProtobufSchema: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtobufSchema: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FqName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FqName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deps", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var keykey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				keykey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			var stringLenmapkey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLenmapkey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLenmapkey := int(stringLenmapkey)
			if intStringLenmapkey < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			postStringIndexmapkey := iNdEx + intStringLenmapkey
			if postStringIndexmapkey > l {
				return io.ErrUnexpectedEOF
			}
			mapkey := string(dAtA[iNdEx:postStringIndexmapkey])
			iNdEx = postStringIndexmapkey
			if m.Deps == nil {
				m.Deps = make(map[string]string)
			}
			if iNdEx < postIndex {
				var valuekey uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowProtobufSchema
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					valuekey |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				var stringLenmapvalue uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowProtobufSchema
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLenmapvalue |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLenmapvalue := int(stringLenmapvalue)
				if intStringLenmapvalue < 0 {
					return ErrInvalidLengthProtobufSchema
				}
				postStringIndexmapvalue := iNdEx + intStringLenmapvalue
				if postStringIndexmapvalue > l {
					return io.ErrUnexpectedEOF
				}
				mapvalue := string(dAtA[iNdEx:postStringIndexmapvalue])
				iNdEx = postStringIndexmapvalue
				m.Deps[mapkey] = mapvalue
			} else {
				var mapvalue string
				m.Deps[mapkey] = mapvalue
			}
			iNdEx = postIndex
		case 30:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &google_protobuf.DescriptorProto{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Descriptor_ = &ProtobufSchema_Msg{v}
			iNdEx = postIndex
		case 31:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enm", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &google_protobuf.EnumDescriptorProto{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Descriptor_ = &ProtobufSchema_Enm{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProtobufSchema(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProtobufSchema
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProtobufSchema(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProtobufSchema
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProtobufSchema
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthProtobufSchema
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowProtobufSchema
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipProtobufSchema(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthProtobufSchema = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProtobufSchema   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("protobuf_schema.proto", fileDescriptorProtobufSchema) }

var fileDescriptorProtobufSchema = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0x2a, 0x4d, 0x8b, 0x2f, 0x4e, 0xce, 0x48, 0xcd, 0x4d, 0xd4, 0x03, 0xf3, 0x85, 0x38,
	0xc1, 0x54, 0x71, 0x72, 0x62, 0x9e, 0x94, 0x42, 0x7a, 0x7e, 0x7e, 0x7a, 0x4e, 0xaa, 0x3e, 0x4c,
	0xa1, 0x7e, 0x4a, 0x6a, 0x71, 0x72, 0x51, 0x66, 0x41, 0x49, 0x7e, 0x11, 0x44, 0xb1, 0x94, 0x2c,
	0x5c, 0x2a, 0x3d, 0x3f, 0x3d, 0x1f, 0xcc, 0x01, 0xb3, 0x20, 0xd2, 0x4a, 0x1b, 0x99, 0xb8, 0xf8,
	0x02, 0xa0, 0x2a, 0x82, 0xc1, 0x96, 0x08, 0x09, 0x70, 0x31, 0x97, 0x66, 0xa6, 0x48, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x06, 0x81, 0x98, 0x42, 0xe2, 0x5c, 0xec, 0x69, 0x85, 0xf1, 0x79, 0x89, 0xb9,
	0xa9, 0x12, 0x4c, 0x60, 0x51, 0xb6, 0xb4, 0x42, 0xbf, 0xc4, 0xdc, 0x54, 0x21, 0x5b, 0x2e, 0x96,
	0x94, 0xd4, 0x82, 0x62, 0x09, 0x16, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x65, 0x3d, 0xb8, 0xc3, 0xf4,
	0x50, 0xcd, 0xd4, 0x73, 0x49, 0x2d, 0x28, 0x76, 0xcd, 0x2b, 0x29, 0xaa, 0x74, 0x62, 0x39, 0x71,
	0x4f, 0x9e, 0x21, 0x08, 0xac, 0x4d, 0xc8, 0x84, 0x8b, 0x39, 0xb7, 0x38, 0x5d, 0x42, 0x4e, 0x81,
	0x51, 0x83, 0xdb, 0x48, 0x41, 0x0f, 0xe2, 0x17, 0x3d, 0x98, 0x83, 0xf5, 0x5c, 0xe0, 0x7e, 0x01,
	0x9b, 0xe6, 0xc1, 0x10, 0x04, 0x52, 0x2e, 0x64, 0xc1, 0xc5, 0x9c, 0x9a, 0x97, 0x2b, 0x21, 0x0f,
	0xd6, 0xa5, 0x82, 0xa1, 0xcb, 0x35, 0xaf, 0x34, 0x17, 0x8b, 0xce, 0xd4, 0xbc, 0x5c, 0x29, 0x73,
	0x2e, 0x4e, 0xb8, 0x43, 0x40, 0xde, 0xcc, 0x4e, 0xad, 0x84, 0x79, 0x33, 0x3b, 0xb5, 0x52, 0x48,
	0x84, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x14, 0xe6, 0x49, 0x08, 0xc7, 0x8a, 0xc9, 0x82, 0xd1, 0x89,
	0x87, 0x8b, 0x0b, 0x11, 0xb0, 0x4e, 0x06, 0x1f, 0x1e, 0xca, 0x31, 0x36, 0x3c, 0x92, 0x63, 0xe8,
	0x78, 0x24, 0xc7, 0x70, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x37, 0x1e, 0xc9, 0x31,
	0x3c, 0x78, 0x24, 0xc7, 0xf8, 0xe3, 0x91, 0x1c, 0x43, 0xc3, 0x63, 0x39, 0x86, 0x09, 0x8f, 0xe5,
	0x18, 0x16, 0x3c, 0x96, 0x63, 0xdc, 0xf0, 0x58, 0x8e, 0x31, 0x89, 0x0d, 0xec, 0x3a, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x11, 0x67, 0x08, 0x27, 0xd1, 0x01, 0x00, 0x00,
}
