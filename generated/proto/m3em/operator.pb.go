// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code generated by protoc-gen-go.
// source: operator.proto
// DO NOT EDIT!

/*
Package m3em is a generated protocol buffer package.

It is generated from these files:
	operator.proto

It has these top-level messages:
	SetupRequest
	SetupResponse
	TransferRequest
	TransferResponse
	StartRequest
	StartResponse
	StopRequest
	StopResponse
	TeardownRequest
	TeardownResponse
*/
package m3em

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FileType int32

const (
	FileType_UNKNOWN        FileType = 0
	FileType_SERVICE_BINARY FileType = 1
	FileType_SERVICE_CONFIG FileType = 2
)

var FileType_name = map[int32]string{
	0: "UNKNOWN",
	1: "SERVICE_BINARY",
	2: "SERVICE_CONFIG",
}
var FileType_value = map[string]int32{
	"UNKNOWN":        0,
	"SERVICE_BINARY": 1,
	"SERVICE_CONFIG": 2,
}

func (x FileType) String() string {
	return proto.EnumName(FileType_name, int32(x))
}
func (FileType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SetupRequest struct {
	SessionToken           string `protobuf:"bytes,1,opt,name=session_token,json=sessionToken" json:"session_token,omitempty"`
	OperatorUuid           string `protobuf:"bytes,2,opt,name=operator_uuid,json=operatorUuid" json:"operator_uuid,omitempty"`
	Force                  bool   `protobuf:"varint,3,opt,name=force" json:"force,omitempty"`
	HeartbeatEnabled       bool   `protobuf:"varint,4,opt,name=heartbeat_enabled,json=heartbeatEnabled" json:"heartbeat_enabled,omitempty"`
	HeartbeatEndpoint      string `protobuf:"bytes,5,opt,name=heartbeat_endpoint,json=heartbeatEndpoint" json:"heartbeat_endpoint,omitempty"`
	HeartbeatFrequencySecs uint32 `protobuf:"varint,6,opt,name=heartbeat_frequency_secs,json=heartbeatFrequencySecs" json:"heartbeat_frequency_secs,omitempty"`
}

func (m *SetupRequest) Reset()                    { *m = SetupRequest{} }
func (m *SetupRequest) String() string            { return proto.CompactTextString(m) }
func (*SetupRequest) ProtoMessage()               {}
func (*SetupRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SetupRequest) GetSessionToken() string {
	if m != nil {
		return m.SessionToken
	}
	return ""
}

func (m *SetupRequest) GetOperatorUuid() string {
	if m != nil {
		return m.OperatorUuid
	}
	return ""
}

func (m *SetupRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

func (m *SetupRequest) GetHeartbeatEnabled() bool {
	if m != nil {
		return m.HeartbeatEnabled
	}
	return false
}

func (m *SetupRequest) GetHeartbeatEndpoint() string {
	if m != nil {
		return m.HeartbeatEndpoint
	}
	return ""
}

func (m *SetupRequest) GetHeartbeatFrequencySecs() uint32 {
	if m != nil {
		return m.HeartbeatFrequencySecs
	}
	return 0
}

type SetupResponse struct {
}

func (m *SetupResponse) Reset()                    { *m = SetupResponse{} }
func (m *SetupResponse) String() string            { return proto.CompactTextString(m) }
func (*SetupResponse) ProtoMessage()               {}
func (*SetupResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type TransferRequest struct {
	Type       FileType `protobuf:"varint,1,opt,name=type,enum=m3em.FileType" json:"type,omitempty"`
	Filename   string   `protobuf:"bytes,2,opt,name=filename" json:"filename,omitempty"`
	Overwrite  bool     `protobuf:"varint,3,opt,name=overwrite" json:"overwrite,omitempty"`
	ChunkIdx   int32    `protobuf:"varint,4,opt,name=chunk_idx,json=chunkIdx" json:"chunk_idx,omitempty"`
	ChunkBytes []byte   `protobuf:"bytes,5,opt,name=chunk_bytes,json=chunkBytes,proto3" json:"chunk_bytes,omitempty"`
}

func (m *TransferRequest) Reset()                    { *m = TransferRequest{} }
func (m *TransferRequest) String() string            { return proto.CompactTextString(m) }
func (*TransferRequest) ProtoMessage()               {}
func (*TransferRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TransferRequest) GetType() FileType {
	if m != nil {
		return m.Type
	}
	return FileType_UNKNOWN
}

func (m *TransferRequest) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *TransferRequest) GetOverwrite() bool {
	if m != nil {
		return m.Overwrite
	}
	return false
}

func (m *TransferRequest) GetChunkIdx() int32 {
	if m != nil {
		return m.ChunkIdx
	}
	return 0
}

func (m *TransferRequest) GetChunkBytes() []byte {
	if m != nil {
		return m.ChunkBytes
	}
	return nil
}

type TransferResponse struct {
	FileChecksum   uint32 `protobuf:"varint,1,opt,name=file_checksum,json=fileChecksum" json:"file_checksum,omitempty"`
	NumChunksRecvd int32  `protobuf:"varint,2,opt,name=num_chunks_recvd,json=numChunksRecvd" json:"num_chunks_recvd,omitempty"`
}

func (m *TransferResponse) Reset()                    { *m = TransferResponse{} }
func (m *TransferResponse) String() string            { return proto.CompactTextString(m) }
func (*TransferResponse) ProtoMessage()               {}
func (*TransferResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *TransferResponse) GetFileChecksum() uint32 {
	if m != nil {
		return m.FileChecksum
	}
	return 0
}

func (m *TransferResponse) GetNumChunksRecvd() int32 {
	if m != nil {
		return m.NumChunksRecvd
	}
	return 0
}

type StartRequest struct {
}

func (m *StartRequest) Reset()                    { *m = StartRequest{} }
func (m *StartRequest) String() string            { return proto.CompactTextString(m) }
func (*StartRequest) ProtoMessage()               {}
func (*StartRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type StartResponse struct {
}

func (m *StartResponse) Reset()                    { *m = StartResponse{} }
func (m *StartResponse) String() string            { return proto.CompactTextString(m) }
func (*StartResponse) ProtoMessage()               {}
func (*StartResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type StopRequest struct {
}

func (m *StopRequest) Reset()                    { *m = StopRequest{} }
func (m *StopRequest) String() string            { return proto.CompactTextString(m) }
func (*StopRequest) ProtoMessage()               {}
func (*StopRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type StopResponse struct {
}

func (m *StopResponse) Reset()                    { *m = StopResponse{} }
func (m *StopResponse) String() string            { return proto.CompactTextString(m) }
func (*StopResponse) ProtoMessage()               {}
func (*StopResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type TeardownRequest struct {
}

func (m *TeardownRequest) Reset()                    { *m = TeardownRequest{} }
func (m *TeardownRequest) String() string            { return proto.CompactTextString(m) }
func (*TeardownRequest) ProtoMessage()               {}
func (*TeardownRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type TeardownResponse struct {
}

func (m *TeardownResponse) Reset()                    { *m = TeardownResponse{} }
func (m *TeardownResponse) String() string            { return proto.CompactTextString(m) }
func (*TeardownResponse) ProtoMessage()               {}
func (*TeardownResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func init() {
	proto.RegisterType((*SetupRequest)(nil), "m3em.SetupRequest")
	proto.RegisterType((*SetupResponse)(nil), "m3em.SetupResponse")
	proto.RegisterType((*TransferRequest)(nil), "m3em.TransferRequest")
	proto.RegisterType((*TransferResponse)(nil), "m3em.TransferResponse")
	proto.RegisterType((*StartRequest)(nil), "m3em.StartRequest")
	proto.RegisterType((*StartResponse)(nil), "m3em.StartResponse")
	proto.RegisterType((*StopRequest)(nil), "m3em.StopRequest")
	proto.RegisterType((*StopResponse)(nil), "m3em.StopResponse")
	proto.RegisterType((*TeardownRequest)(nil), "m3em.TeardownRequest")
	proto.RegisterType((*TeardownResponse)(nil), "m3em.TeardownResponse")
	proto.RegisterEnum("m3em.FileType", FileType_name, FileType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Operator service

type OperatorClient interface {
	Setup(ctx context.Context, in *SetupRequest, opts ...grpc.CallOption) (*SetupResponse, error)
	Transfer(ctx context.Context, opts ...grpc.CallOption) (Operator_TransferClient, error)
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error)
	Teardown(ctx context.Context, in *TeardownRequest, opts ...grpc.CallOption) (*TeardownResponse, error)
}

type operatorClient struct {
	cc *grpc.ClientConn
}

func NewOperatorClient(cc *grpc.ClientConn) OperatorClient {
	return &operatorClient{cc}
}

func (c *operatorClient) Setup(ctx context.Context, in *SetupRequest, opts ...grpc.CallOption) (*SetupResponse, error) {
	out := new(SetupResponse)
	err := grpc.Invoke(ctx, "/m3em.Operator/Setup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operatorClient) Transfer(ctx context.Context, opts ...grpc.CallOption) (Operator_TransferClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Operator_serviceDesc.Streams[0], c.cc, "/m3em.Operator/Transfer", opts...)
	if err != nil {
		return nil, err
	}
	x := &operatorTransferClient{stream}
	return x, nil
}

type Operator_TransferClient interface {
	Send(*TransferRequest) error
	CloseAndRecv() (*TransferResponse, error)
	grpc.ClientStream
}

type operatorTransferClient struct {
	grpc.ClientStream
}

func (x *operatorTransferClient) Send(m *TransferRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *operatorTransferClient) CloseAndRecv() (*TransferResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(TransferResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *operatorClient) Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error) {
	out := new(StartResponse)
	err := grpc.Invoke(ctx, "/m3em.Operator/Start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operatorClient) Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error) {
	out := new(StopResponse)
	err := grpc.Invoke(ctx, "/m3em.Operator/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operatorClient) Teardown(ctx context.Context, in *TeardownRequest, opts ...grpc.CallOption) (*TeardownResponse, error) {
	out := new(TeardownResponse)
	err := grpc.Invoke(ctx, "/m3em.Operator/Teardown", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Operator service

type OperatorServer interface {
	Setup(context.Context, *SetupRequest) (*SetupResponse, error)
	Transfer(Operator_TransferServer) error
	Start(context.Context, *StartRequest) (*StartResponse, error)
	Stop(context.Context, *StopRequest) (*StopResponse, error)
	Teardown(context.Context, *TeardownRequest) (*TeardownResponse, error)
}

func RegisterOperatorServer(s *grpc.Server, srv OperatorServer) {
	s.RegisterService(&_Operator_serviceDesc, srv)
}

func _Operator_Setup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperatorServer).Setup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3em.Operator/Setup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperatorServer).Setup(ctx, req.(*SetupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operator_Transfer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OperatorServer).Transfer(&operatorTransferServer{stream})
}

type Operator_TransferServer interface {
	SendAndClose(*TransferResponse) error
	Recv() (*TransferRequest, error)
	grpc.ServerStream
}

type operatorTransferServer struct {
	grpc.ServerStream
}

func (x *operatorTransferServer) SendAndClose(m *TransferResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *operatorTransferServer) Recv() (*TransferRequest, error) {
	m := new(TransferRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Operator_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperatorServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3em.Operator/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperatorServer).Start(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operator_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperatorServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3em.Operator/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperatorServer).Stop(ctx, req.(*StopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operator_Teardown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TeardownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperatorServer).Teardown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3em.Operator/Teardown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperatorServer).Teardown(ctx, req.(*TeardownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Operator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m3em.Operator",
	HandlerType: (*OperatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Setup",
			Handler:    _Operator_Setup_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _Operator_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Operator_Stop_Handler,
		},
		{
			MethodName: "Teardown",
			Handler:    _Operator_Teardown_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Transfer",
			Handler:       _Operator_Transfer_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "operator.proto",
}

func init() { proto.RegisterFile("operator.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 551 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x93, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x49, 0x69, 0x47, 0xf6, 0xd6, 0x66, 0xa9, 0x81, 0x29, 0x2a, 0x48, 0x4c, 0xe1, 0x52,
	0x81, 0x36, 0xa1, 0xed, 0x02, 0xe2, 0x80, 0x58, 0xb5, 0xa1, 0x0a, 0xa9, 0x95, 0xdc, 0x0e, 0xc4,
	0x29, 0x4a, 0x93, 0x57, 0x35, 0x6a, 0x63, 0x07, 0xdb, 0xd9, 0xd6, 0x3b, 0x5f, 0x81, 0x8f, 0xc1,
	0x77, 0x44, 0xb1, 0x93, 0x34, 0xeb, 0xd1, 0xbf, 0xf7, 0x6f, 0x9e, 0xfd, 0x7b, 0xaf, 0xe0, 0xf0,
	0x0c, 0x45, 0xa8, 0xb8, 0x38, 0xcf, 0x04, 0x57, 0x9c, 0xb4, 0xd3, 0x4b, 0x4c, 0xfd, 0x3f, 0x2d,
	0xe8, 0xce, 0x50, 0xe5, 0x19, 0xc5, 0xdf, 0x39, 0x4a, 0x45, 0xde, 0x42, 0x4f, 0xa2, 0x94, 0x09,
	0x67, 0x81, 0xe2, 0x6b, 0x64, 0x9e, 0x75, 0x6a, 0x0d, 0x0f, 0x69, 0xb7, 0x84, 0xf3, 0x82, 0x15,
	0xa1, 0xea, 0x6b, 0x41, 0x9e, 0x27, 0xb1, 0xd7, 0x32, 0xa1, 0x0a, 0xde, 0xe6, 0x49, 0x4c, 0x5e,
	0x40, 0x67, 0xc9, 0x45, 0x84, 0xde, 0xd3, 0x53, 0x6b, 0x68, 0x53, 0x73, 0x20, 0xef, 0xa1, 0xbf,
	0xc2, 0x50, 0xa8, 0x05, 0x86, 0x2a, 0x40, 0x16, 0x2e, 0x36, 0x18, 0x7b, 0x6d, 0x9d, 0x70, 0xeb,
	0xc2, 0xb5, 0xe1, 0xe4, 0x0c, 0x48, 0x33, 0x1c, 0x67, 0x3c, 0x61, 0xca, 0xeb, 0xe8, 0x66, 0xfd,
	0x46, 0xda, 0x14, 0xc8, 0x47, 0xf0, 0x76, 0xf1, 0xa5, 0x28, 0x5e, 0xc4, 0xa2, 0x6d, 0x20, 0x31,
	0x92, 0xde, 0xc1, 0xa9, 0x35, 0xec, 0xd1, 0x93, 0xba, 0x7e, 0x53, 0x95, 0x67, 0x18, 0x49, 0xff,
	0x18, 0x7a, 0xa5, 0x05, 0x99, 0x71, 0x26, 0xd1, 0xff, 0x67, 0xc1, 0xf1, 0x5c, 0x84, 0x4c, 0x2e,
	0x51, 0x54, 0x6a, 0x7c, 0x68, 0xab, 0x6d, 0x86, 0xda, 0x88, 0x73, 0xe1, 0x9c, 0x17, 0x02, 0xcf,
	0x6f, 0x92, 0x0d, 0xce, 0xb7, 0x19, 0x52, 0x5d, 0x23, 0x03, 0xb0, 0x97, 0xc9, 0x06, 0x59, 0x98,
	0x62, 0x29, 0xa5, 0x3e, 0x93, 0xd7, 0x70, 0xc8, 0xef, 0x50, 0xdc, 0x8b, 0x44, 0x55, 0x52, 0x76,
	0x80, 0xbc, 0x82, 0xc3, 0x68, 0x95, 0xb3, 0x75, 0x90, 0xc4, 0x0f, 0x5a, 0x48, 0x87, 0xda, 0x1a,
	0x8c, 0xe3, 0x07, 0xf2, 0x06, 0x8e, 0x4c, 0x71, 0xb1, 0x55, 0x28, 0xb5, 0x81, 0x2e, 0x05, 0x8d,
	0xae, 0x0a, 0xe2, 0x87, 0xe0, 0xee, 0xae, 0x6b, 0xde, 0x50, 0x4c, 0xa9, 0xe8, 0x1d, 0x44, 0x2b,
	0x8c, 0xd6, 0x32, 0x4f, 0xf5, 0xc5, 0x7b, 0xb4, 0x5b, 0xc0, 0x51, 0xc9, 0xc8, 0x10, 0x5c, 0x96,
	0xa7, 0x81, 0xfe, 0x94, 0x0c, 0x04, 0x46, 0x77, 0x66, 0x9a, 0x1d, 0xea, 0xb0, 0x3c, 0x1d, 0x69,
	0x4c, 0x0b, 0xea, 0x3b, 0xd0, 0x9d, 0xa9, 0x50, 0xa8, 0x52, 0x87, 0x76, 0x66, 0xce, 0xa5, 0xb3,
	0x1e, 0x1c, 0xcd, 0x14, 0xaf, 0x36, 0xc9, 0xe4, 0xf9, 0x4e, 0x69, 0x1f, 0x8e, 0xe7, 0x18, 0x8a,
	0x98, 0xdf, 0xb3, 0x2a, 0x42, 0xc0, 0xdd, 0x21, 0x13, 0x7b, 0xf7, 0x05, 0xec, 0xca, 0x29, 0x39,
	0x82, 0x67, 0xb7, 0x93, 0xef, 0x93, 0xe9, 0xcf, 0x89, 0xfb, 0x84, 0x10, 0x70, 0x66, 0xd7, 0xf4,
	0xc7, 0x78, 0x74, 0x1d, 0x5c, 0x8d, 0x27, 0x5f, 0xe9, 0x2f, 0xd7, 0x6a, 0xb2, 0xd1, 0x74, 0x72,
	0x33, 0xfe, 0xe6, 0xb6, 0x2e, 0xfe, 0xb6, 0xc0, 0x9e, 0x96, 0x8b, 0x48, 0x3e, 0x40, 0x47, 0x0f,
	0x96, 0x10, 0x33, 0xae, 0xe6, 0xae, 0x0f, 0x9e, 0x3f, 0x62, 0xa5, 0xb5, 0xcf, 0x60, 0x57, 0x26,
	0xc9, 0x4b, 0x13, 0xd8, 0x5b, 0x84, 0xc1, 0xc9, 0x3e, 0x36, 0x3f, 0x1d, 0x5a, 0xba, 0x5d, 0xe1,
	0xa4, 0x6e, 0xd7, 0x10, 0x56, 0xb7, 0x6b, 0x4a, 0x23, 0x67, 0xd0, 0x2e, 0x2c, 0x91, 0x7e, 0x55,
	0xac, 0x05, 0x0e, 0x48, 0x13, 0x95, 0xf1, 0x4f, 0x60, 0x57, 0xc6, 0xea, 0xdb, 0x3d, 0x96, 0x5a,
	0xdf, 0x6e, 0x4f, 0xec, 0xe2, 0x40, 0xff, 0xef, 0x2f, 0xff, 0x07, 0x00, 0x00, 0xff, 0xff, 0xac,
	0xa6, 0x95, 0x32, 0x09, 0x04, 0x00, 0x00,
}
