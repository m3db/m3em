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
	CleanDataDirectoryRequest
	CleanDataDirectoryResponse
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
	FileType_UNKNOWN     FileType = 0
	FileType_M3DB_BINARY FileType = 1
	FileType_M3DB_CONFIG FileType = 2
)

var FileType_name = map[int32]string{
	0: "UNKNOWN",
	1: "M3DB_BINARY",
	2: "M3DB_CONFIG",
}
var FileType_value = map[string]int32{
	"UNKNOWN":     0,
	"M3DB_BINARY": 1,
	"M3DB_CONFIG": 2,
}

func (x FileType) String() string {
	return proto.EnumName(FileType_name, int32(x))
}
func (FileType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SetupRequest struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Force bool   `protobuf:"varint,2,opt,name=force" json:"force,omitempty"`
}

func (m *SetupRequest) Reset()                    { *m = SetupRequest{} }
func (m *SetupRequest) String() string            { return proto.CompactTextString(m) }
func (*SetupRequest) ProtoMessage()               {}
func (*SetupRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SetupRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SetupRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
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

type CleanDataDirectoryRequest struct {
}

func (m *CleanDataDirectoryRequest) Reset()                    { *m = CleanDataDirectoryRequest{} }
func (m *CleanDataDirectoryRequest) String() string            { return proto.CompactTextString(m) }
func (*CleanDataDirectoryRequest) ProtoMessage()               {}
func (*CleanDataDirectoryRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type CleanDataDirectoryResponse struct {
}

func (m *CleanDataDirectoryResponse) Reset()                    { *m = CleanDataDirectoryResponse{} }
func (m *CleanDataDirectoryResponse) String() string            { return proto.CompactTextString(m) }
func (*CleanDataDirectoryResponse) ProtoMessage()               {}
func (*CleanDataDirectoryResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

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
	proto.RegisterType((*CleanDataDirectoryRequest)(nil), "m3em.CleanDataDirectoryRequest")
	proto.RegisterType((*CleanDataDirectoryResponse)(nil), "m3em.CleanDataDirectoryResponse")
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
	// rpc CleanDataDirectory(CleanDataDirectoryRequest) returns (CleanDataDirectoryResponse);
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
	// rpc CleanDataDirectory(CleanDataDirectoryRequest) returns (CleanDataDirectoryResponse);
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
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x93, 0xdf, 0x8f, 0xd2, 0x40,
	0x10, 0xc7, 0x2d, 0x52, 0x2d, 0x03, 0x94, 0xde, 0xfa, 0x23, 0xb5, 0x77, 0x89, 0xa4, 0xbe, 0x34,
	0x26, 0x12, 0x73, 0x3c, 0x9d, 0x3e, 0x09, 0xe4, 0x0c, 0x31, 0x96, 0x64, 0xc1, 0x18, 0x9f, 0x9a,
	0x5e, 0x19, 0x72, 0x0d, 0xb4, 0x5b, 0xb7, 0xdb, 0xbb, 0xe3, 0x0f, 0xf1, 0xcf, 0xf0, 0x7f, 0x34,
	0xdd, 0x6d, 0x0b, 0x47, 0x7c, 0x9c, 0xcf, 0xcc, 0xec, 0x77, 0xfa, 0x9d, 0x29, 0x98, 0x2c, 0x43,
	0x1e, 0x0a, 0xc6, 0x47, 0x19, 0x67, 0x82, 0x91, 0x76, 0x32, 0xc6, 0xc4, 0xfd, 0x04, 0xbd, 0x25,
	0x8a, 0x22, 0xa3, 0xf8, 0xbb, 0xc0, 0x5c, 0x90, 0x97, 0xa0, 0x0b, 0xb6, 0xc5, 0xd4, 0xd6, 0x86,
	0x9a, 0xd7, 0xa1, 0x2a, 0x28, 0xe9, 0x86, 0xf1, 0x08, 0xed, 0xd6, 0x50, 0xf3, 0x0c, 0xaa, 0x02,
	0x77, 0x00, 0xfd, 0xaa, 0x37, 0xcf, 0x58, 0x9a, 0xa3, 0xfb, 0x57, 0x83, 0xc1, 0x8a, 0x87, 0x69,
	0xbe, 0x41, 0x5e, 0x3f, 0xe8, 0x42, 0x5b, 0xec, 0x33, 0x94, 0xef, 0x99, 0x97, 0xe6, 0xa8, 0x54,
	0x1d, 0x5d, 0xc7, 0x3b, 0x5c, 0xed, 0x33, 0xa4, 0x32, 0x47, 0x1c, 0x30, 0x36, 0xf1, 0x0e, 0xd3,
	0x30, 0x51, 0x0a, 0x1d, 0xda, 0xc4, 0xe4, 0x02, 0x3a, 0xec, 0x0e, 0xf9, 0x3d, 0x8f, 0x05, 0xda,
	0x4f, 0xa5, 0xfc, 0x01, 0x90, 0x73, 0xe8, 0x44, 0xb7, 0x45, 0xba, 0x0d, 0xe2, 0xf5, 0x83, 0xdd,
	0x1e, 0x6a, 0x9e, 0x4e, 0x0d, 0x09, 0xe6, 0xeb, 0x07, 0xf2, 0x16, 0xba, 0x2a, 0x79, 0xb3, 0x17,
	0x98, 0xdb, 0xfa, 0x50, 0xf3, 0x7a, 0x14, 0x24, 0x9a, 0x94, 0xc4, 0x0d, 0xc1, 0x3a, 0x8c, 0xab,
	0xbe, 0x81, 0xbc, 0x83, 0x7e, 0xa9, 0x1d, 0x44, 0xb7, 0x18, 0x6d, 0xf3, 0x22, 0x91, 0x83, 0xf7,
	0x69, 0xaf, 0x84, 0xd3, 0x8a, 0x11, 0x0f, 0xac, 0xb4, 0x48, 0x02, 0xf9, 0x54, 0x1e, 0x70, 0x8c,
	0xee, 0xd6, 0x72, 0x70, 0x9d, 0x9a, 0x69, 0x91, 0x4c, 0x25, 0xa6, 0x25, 0x75, 0x4d, 0xe8, 0x2d,
	0x45, 0xc8, 0x45, 0x65, 0x87, 0xf4, 0x4c, 0xc5, 0x95, 0x67, 0x7d, 0xe8, 0x2e, 0x05, 0xab, 0xfd,
	0x57, 0xf5, 0xec, 0x60, 0xe9, 0x19, 0x0c, 0x56, 0x18, 0xf2, 0x35, 0xbb, 0x4f, 0xeb, 0x12, 0x02,
	0xd6, 0x01, 0x55, 0x65, 0xe7, 0xf0, 0x66, 0xba, 0xc3, 0x30, 0x9d, 0x85, 0x22, 0x9c, 0xc5, 0x1c,
	0x23, 0xc1, 0xf8, 0xbe, 0x6e, 0xb8, 0x00, 0xe7, 0x7f, 0x49, 0xd5, 0xfa, 0xfe, 0x0a, 0x8c, 0x7a,
	0x1d, 0xa4, 0x0b, 0xcf, 0x7f, 0xf8, 0xdf, 0xfc, 0xc5, 0x4f, 0xdf, 0x7a, 0x42, 0x06, 0xd0, 0xfd,
	0x3e, 0x9e, 0x4d, 0x82, 0xc9, 0xdc, 0xff, 0x42, 0x7f, 0x59, 0x5a, 0x03, 0xa6, 0x0b, 0xff, 0x7a,
	0xfe, 0xd5, 0x6a, 0x5d, 0xfe, 0x69, 0x81, 0xb1, 0xa8, 0xae, 0x8a, 0x7c, 0x04, 0x5d, 0x5e, 0x03,
	0x21, 0x6a, 0xc7, 0xc7, 0x67, 0xe5, 0xbc, 0x78, 0xc4, 0x2a, 0xab, 0x3f, 0x83, 0x51, 0xdb, 0x4f,
	0x5e, 0xa9, 0x82, 0x93, 0xeb, 0x71, 0x5e, 0x9f, 0x62, 0xd5, 0xea, 0x69, 0x52, 0xae, 0x34, 0xb2,
	0x91, 0x3b, 0x72, 0xb9, 0x91, 0x3b, 0x76, 0x9a, 0x7c, 0x80, 0x76, 0x69, 0x2d, 0x39, 0xab, 0x93,
	0x8d, 0xeb, 0x0e, 0x39, 0x46, 0x55, 0xf9, 0x15, 0x18, 0xb5, 0xcd, 0xcd, 0x74, 0x8f, 0x37, 0xd1,
	0x4c, 0x77, 0xb2, 0x8d, 0x9b, 0x67, 0xf2, 0x0f, 0x1b, 0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x0a,
	0x15, 0x46, 0x00, 0x73, 0x03, 0x00, 0x00,
}