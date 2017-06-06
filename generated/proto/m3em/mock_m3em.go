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

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/m3db/m3em/generated/proto/m3em (interfaces: OperatorClient,Operator_PushFileClient,Operator_PullFileClient)

package m3em

import (
	gomock "github.com/golang/mock/gomock"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// Mock of OperatorClient interface
type MockOperatorClient struct {
	ctrl     *gomock.Controller
	recorder *_MockOperatorClientRecorder
}

// Recorder for MockOperatorClient (not exported)
type _MockOperatorClientRecorder struct {
	mock *MockOperatorClient
}

func NewMockOperatorClient(ctrl *gomock.Controller) *MockOperatorClient {
	mock := &MockOperatorClient{ctrl: ctrl}
	mock.recorder = &_MockOperatorClientRecorder{mock}
	return mock
}

func (_m *MockOperatorClient) EXPECT() *_MockOperatorClientRecorder {
	return _m.recorder
}

func (_m *MockOperatorClient) PullFile(_param0 context.Context, _param1 *PullFileRequest, _param2 ...grpc.CallOption) (Operator_PullFileClient, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "PullFile", _s...)
	ret0, _ := ret[0].(Operator_PullFileClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperatorClientRecorder) PullFile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PullFile", _s...)
}

func (_m *MockOperatorClient) PushFile(_param0 context.Context, _param1 ...grpc.CallOption) (Operator_PushFileClient, error) {
	_s := []interface{}{_param0}
	for _, _x := range _param1 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "PushFile", _s...)
	ret0, _ := ret[0].(Operator_PushFileClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperatorClientRecorder) PushFile(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0}, arg1...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PushFile", _s...)
}

func (_m *MockOperatorClient) Setup(_param0 context.Context, _param1 *SetupRequest, _param2 ...grpc.CallOption) (*SetupResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Setup", _s...)
	ret0, _ := ret[0].(*SetupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperatorClientRecorder) Setup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Setup", _s...)
}

func (_m *MockOperatorClient) Start(_param0 context.Context, _param1 *StartRequest, _param2 ...grpc.CallOption) (*StartResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Start", _s...)
	ret0, _ := ret[0].(*StartResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperatorClientRecorder) Start(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start", _s...)
}

func (_m *MockOperatorClient) Stop(_param0 context.Context, _param1 *StopRequest, _param2 ...grpc.CallOption) (*StopResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Stop", _s...)
	ret0, _ := ret[0].(*StopResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperatorClientRecorder) Stop(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop", _s...)
}

func (_m *MockOperatorClient) Teardown(_param0 context.Context, _param1 *TeardownRequest, _param2 ...grpc.CallOption) (*TeardownResponse, error) {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Teardown", _s...)
	ret0, _ := ret[0].(*TeardownResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperatorClientRecorder) Teardown(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Teardown", _s...)
}

// Mock of Operator_PushFileClient interface
type MockOperator_PushFileClient struct {
	ctrl     *gomock.Controller
	recorder *_MockOperator_PushFileClientRecorder
}

// Recorder for MockOperator_PushFileClient (not exported)
type _MockOperator_PushFileClientRecorder struct {
	mock *MockOperator_PushFileClient
}

func NewMockOperator_PushFileClient(ctrl *gomock.Controller) *MockOperator_PushFileClient {
	mock := &MockOperator_PushFileClient{ctrl: ctrl}
	mock.recorder = &_MockOperator_PushFileClientRecorder{mock}
	return mock
}

func (_m *MockOperator_PushFileClient) EXPECT() *_MockOperator_PushFileClientRecorder {
	return _m.recorder
}

func (_m *MockOperator_PushFileClient) CloseAndRecv() (*PushFileResponse, error) {
	ret := _m.ctrl.Call(_m, "CloseAndRecv")
	ret0, _ := ret[0].(*PushFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperator_PushFileClientRecorder) CloseAndRecv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CloseAndRecv")
}

func (_m *MockOperator_PushFileClient) CloseSend() error {
	ret := _m.ctrl.Call(_m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PushFileClientRecorder) CloseSend() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CloseSend")
}

func (_m *MockOperator_PushFileClient) Context() context.Context {
	ret := _m.ctrl.Call(_m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

func (_mr *_MockOperator_PushFileClientRecorder) Context() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Context")
}

func (_m *MockOperator_PushFileClient) Header() (metadata.MD, error) {
	ret := _m.ctrl.Call(_m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperator_PushFileClientRecorder) Header() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Header")
}

func (_m *MockOperator_PushFileClient) RecvMsg(_param0 interface{}) error {
	ret := _m.ctrl.Call(_m, "RecvMsg", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PushFileClientRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RecvMsg", arg0)
}

func (_m *MockOperator_PushFileClient) Send(_param0 *PushFileRequest) error {
	ret := _m.ctrl.Call(_m, "Send", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PushFileClientRecorder) Send(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Send", arg0)
}

func (_m *MockOperator_PushFileClient) SendMsg(_param0 interface{}) error {
	ret := _m.ctrl.Call(_m, "SendMsg", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PushFileClientRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SendMsg", arg0)
}

func (_m *MockOperator_PushFileClient) Trailer() metadata.MD {
	ret := _m.ctrl.Call(_m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

func (_mr *_MockOperator_PushFileClientRecorder) Trailer() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Trailer")
}

// Mock of Operator_PullFileClient interface
type MockOperator_PullFileClient struct {
	ctrl     *gomock.Controller
	recorder *_MockOperator_PullFileClientRecorder
}

// Recorder for MockOperator_PullFileClient (not exported)
type _MockOperator_PullFileClientRecorder struct {
	mock *MockOperator_PullFileClient
}

func NewMockOperator_PullFileClient(ctrl *gomock.Controller) *MockOperator_PullFileClient {
	mock := &MockOperator_PullFileClient{ctrl: ctrl}
	mock.recorder = &_MockOperator_PullFileClientRecorder{mock}
	return mock
}

func (_m *MockOperator_PullFileClient) EXPECT() *_MockOperator_PullFileClientRecorder {
	return _m.recorder
}

func (_m *MockOperator_PullFileClient) CloseSend() error {
	ret := _m.ctrl.Call(_m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PullFileClientRecorder) CloseSend() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CloseSend")
}

func (_m *MockOperator_PullFileClient) Context() context.Context {
	ret := _m.ctrl.Call(_m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

func (_mr *_MockOperator_PullFileClientRecorder) Context() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Context")
}

func (_m *MockOperator_PullFileClient) Header() (metadata.MD, error) {
	ret := _m.ctrl.Call(_m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperator_PullFileClientRecorder) Header() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Header")
}

func (_m *MockOperator_PullFileClient) Recv() (*PullFileResponse, error) {
	ret := _m.ctrl.Call(_m, "Recv")
	ret0, _ := ret[0].(*PullFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOperator_PullFileClientRecorder) Recv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Recv")
}

func (_m *MockOperator_PullFileClient) RecvMsg(_param0 interface{}) error {
	ret := _m.ctrl.Call(_m, "RecvMsg", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PullFileClientRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RecvMsg", arg0)
}

func (_m *MockOperator_PullFileClient) SendMsg(_param0 interface{}) error {
	ret := _m.ctrl.Call(_m, "SendMsg", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperator_PullFileClientRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SendMsg", arg0)
}

func (_m *MockOperator_PullFileClient) Trailer() metadata.MD {
	ret := _m.ctrl.Call(_m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

func (_mr *_MockOperator_PullFileClientRecorder) Trailer() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Trailer")
}
