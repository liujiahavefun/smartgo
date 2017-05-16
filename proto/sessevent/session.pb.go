// Code generated by protoc-gen-go.
// source: session.proto
// DO NOT EDIT!

/*
Package sessevent is a generated protocol buffer package.

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
package sessevent

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
	proto.RegisterType((*SessionAccepted)(nil), "sessevent.SessionAccepted")
	proto.RegisterType((*SessionAcceptFailed)(nil), "sessevent.SessionAcceptFailed")
	proto.RegisterType((*SessionConnected)(nil), "sessevent.SessionConnected")
	proto.RegisterType((*SessionConnectFailed)(nil), "sessevent.SessionConnectFailed")
	proto.RegisterType((*SessionClosed)(nil), "sessevent.SessionClosed")
	proto.RegisterType((*SessionError)(nil), "sessevent.SessionError")
}

func init() { proto.RegisterFile("session.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2e,
	0xce, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x04, 0x71, 0x53, 0xcb, 0x52,
	0xf3, 0x4a, 0x94, 0x04, 0xb9, 0xf8, 0x83, 0x21, 0x72, 0x8e, 0xc9, 0xc9, 0xa9, 0x05, 0x25, 0xa9,
	0x29, 0x4a, 0xba, 0x5c, 0xc2, 0x28, 0x42, 0x6e, 0x89, 0x99, 0x39, 0xa9, 0x29, 0x42, 0x62, 0x5c,
	0x6c, 0x41, 0xa9, 0x89, 0xc5, 0xf9, 0x79, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x50, 0x9e,
	0x92, 0x10, 0x97, 0x00, 0x54, 0xb9, 0x73, 0x7e, 0x5e, 0x5e, 0x6a, 0x32, 0xc8, 0x08, 0x3d, 0x2e,
	0x11, 0x54, 0x31, 0x02, 0x66, 0xa8, 0x73, 0xf1, 0xc2, 0xd4, 0xe7, 0xe4, 0x17, 0xe3, 0x51, 0xa8,
	0xc6, 0xc5, 0x03, 0x55, 0xe8, 0x5a, 0x54, 0x94, 0x5f, 0x84, 0x4b, 0x5d, 0x12, 0x1b, 0xd8, 0xa3,
	0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x26, 0x5a, 0x76, 0xf9, 0x00, 0x00, 0x00,
}