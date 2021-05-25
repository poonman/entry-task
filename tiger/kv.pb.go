// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kv.proto

package kv

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type HelloReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReq) Reset()         { *m = HelloReq{} }
func (m *HelloReq) String() string { return proto.CompactTextString(m) }
func (*HelloReq) ProtoMessage()    {}
func (*HelloReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216fe83c9c12408, []int{0}
}

func (m *HelloReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloReq.Unmarshal(m, b)
}
func (m *HelloReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloReq.Marshal(b, m, deterministic)
}
func (m *HelloReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReq.Merge(m, src)
}
func (m *HelloReq) XXX_Size() int {
	return xxx_messageInfo_HelloReq.Size(m)
}
func (m *HelloReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReq.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReq proto.InternalMessageInfo

func (m *HelloReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type HelloRsp struct {
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloRsp) Reset()         { *m = HelloRsp{} }
func (m *HelloRsp) String() string { return proto.CompactTextString(m) }
func (*HelloRsp) ProtoMessage()    {}
func (*HelloRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216fe83c9c12408, []int{1}
}

func (m *HelloRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloRsp.Unmarshal(m, b)
}
func (m *HelloRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloRsp.Marshal(b, m, deterministic)
}
func (m *HelloRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRsp.Merge(m, src)
}
func (m *HelloRsp) XXX_Size() int {
	return xxx_messageInfo_HelloRsp.Size(m)
}
func (m *HelloRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRsp.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRsp proto.InternalMessageInfo

func (m *HelloRsp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloReq)(nil), "kv.HelloReq")
	proto.RegisterType((*HelloRsp)(nil), "kv.HelloRsp")
}

func init() {
	proto.RegisterFile("kv.proto", fileDescriptor_2216fe83c9c12408)
}

var fileDescriptor_2216fe83c9c12408 = []byte{
	// 124 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xc8, 0x2e, 0xd3, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xca, 0x2e, 0x53, 0x92, 0xe3, 0xe2, 0xf0, 0x48, 0xcd, 0xc9,
	0xc9, 0x0f, 0x4a, 0x2d, 0x14, 0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60,
	0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95, 0x64, 0x60, 0xf2, 0xc5, 0x05, 0x42, 0x02, 0x5c, 0xcc, 0xb9,
	0xc5, 0xe9, 0x12, 0x4c, 0x60, 0x69, 0x10, 0xd3, 0x48, 0x87, 0x8b, 0xc9, 0x3b, 0x4c, 0x48, 0x8d,
	0x8b, 0x23, 0x38, 0xb1, 0x12, 0xac, 0x4c, 0x88, 0x47, 0x2f, 0xbb, 0x4c, 0x0f, 0x66, 0xa2, 0x14,
	0x12, 0xaf, 0xb8, 0xc0, 0x89, 0x25, 0x8a, 0x29, 0xbb, 0x2c, 0x89, 0x0d, 0x6c, 0xb9, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x51, 0xd8, 0x04, 0xa9, 0x88, 0x00, 0x00, 0x00,
}
