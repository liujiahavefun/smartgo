// Code generated by protoc-gen-go.
// source: session.proto
// DO NOT EDIT!

/*
Package session is a generated protocol buffer package.

It is generated from these files:
	session.proto

It has these top-level messages:
	SessionAccepted
	SessionAcceptFailed
	SessionConnected
	SessionConnectFailed
	SessionClosed
	SessionError
*/
package session

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// 一个连接被Accept
type SessionAccepted struct {
}

func (m *SessionAccepted) Reset()                    { *m = SessionAccepted{} }
func (m *SessionAccepted) String() string            { return proto.CompactTextString(m) }
func (*SessionAccepted) ProtoMessage()               {}
func (*SessionAccepted) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Accept发生错误
type SessionAcceptFailed struct {
	Reason string `protobuf:"bytes,1,opt,name=Reason" json:"Reason,omitempty"`
}

func (m *SessionAcceptFailed) Reset()                    { *m = SessionAcceptFailed{} }
func (m *SessionAcceptFailed) String() string            { return proto.CompactTextString(m) }
func (*SessionAcceptFailed) ProtoMessage()               {}
func (*SessionAcceptFailed) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// Session已经建立
type SessionConnected struct {
}

func (m *SessionConnected) Reset()                    { *m = SessionConnected{} }
func (m *SessionConnected) String() string            { return proto.CompactTextString(m) }
func (*SessionConnected) ProtoMessage()               {}
func (*SessionConnected) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// Session建立出现错误
type SessionConnectFailed struct {
	Reason string `protobuf:"bytes,1,opt,name=Reason" json:"Reason,omitempty"`
}

func (m *SessionConnectFailed) Reset()                    { *m = SessionConnectFailed{} }
func (m *SessionConnectFailed) String() string            { return proto.CompactTextString(m) }
func (*SessionConnectFailed) ProtoMessage()               {}
func (*SessionConnectFailed) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// Session关闭，因为错误被关闭，就填eason好喽
type SessionClosed struct {
	Reason string `protobuf:"bytes,1,opt,name=Reason" json:"Reason,omitempty"`
}

func (m *SessionClosed) Reset()                    { *m = SessionClosed{} }
func (m *SessionClosed) String() string            { return proto.CompactTextString(m) }
func (*SessionClosed) ProtoMessage()               {}
func (*SessionClosed) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// Session有错误
type SessionError struct {
	Reason string `protobuf:"bytes,1,opt,name=Reason" json:"Reason,omitempty"`
}

func (m *SessionError) Reset()                    { *m = SessionError{} }
func (m *SessionError) String() string            { return proto.CompactTextString(m) }
func (*SessionError) ProtoMessage()               {}
func (*SessionError) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*SessionAccepted)(nil), "session.SessionAccepted")
	proto.RegisterType((*SessionAcceptFailed)(nil), "session.SessionAcceptFailed")
	proto.RegisterType((*SessionConnected)(nil), "session.SessionConnected")
	proto.RegisterType((*SessionConnectFailed)(nil), "session.SessionConnectFailed")
	proto.RegisterType((*SessionClosed)(nil), "session.SessionClosed")
	proto.RegisterType((*SessionError)(nil), "session.SessionError")
}

func init() { proto.RegisterFile("session.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 138 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2e,
	0xce, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x04, 0xb9,
	0xf8, 0x83, 0x21, 0x4c, 0xc7, 0xe4, 0xe4, 0xd4, 0x82, 0x92, 0xd4, 0x14, 0x25, 0x5d, 0x2e, 0x61,
	0x14, 0x21, 0xb7, 0xc4, 0xcc, 0x9c, 0xd4, 0x14, 0x21, 0x31, 0x2e, 0xb6, 0xa0, 0xd4, 0xc4, 0xe2,
	0xfc, 0x3c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x28, 0x4f, 0x49, 0x88, 0x4b, 0x00, 0xaa,
	0xdc, 0x39, 0x3f, 0x2f, 0x2f, 0x35, 0x19, 0x64, 0x84, 0x1e, 0x97, 0x08, 0xaa, 0x18, 0x01, 0x33,
	0xd4, 0xb9, 0x78, 0x61, 0xea, 0x73, 0xf2, 0x8b, 0xf1, 0x28, 0x54, 0xe3, 0xe2, 0x81, 0x2a, 0x74,
	0x2d, 0x2a, 0xca, 0x2f, 0xc2, 0xa5, 0x2e, 0x89, 0x0d, 0xec, 0x4d, 0x63, 0x40, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x2d, 0xba, 0xc4, 0xd7, 0xf7, 0x00, 0x00, 0x00,
}
