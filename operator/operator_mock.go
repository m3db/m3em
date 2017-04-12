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
// Source: github.com/m3db/m3em/operator (interfaces: Operator)

package operator

import (
	gomock "github.com/golang/mock/gomock"
	build "github.com/m3db/m3em/build"

)

// Mock of Operator interface
type MockOperator struct {
	ctrl     *gomock.Controller
	recorder *_MockOperatorRecorder
}

// Recorder for MockOperator (not exported)
type _MockOperatorRecorder struct {
	mock *MockOperator
}

func NewMockOperator(ctrl *gomock.Controller) *MockOperator {
	mock := &MockOperator{ctrl: ctrl}
	mock.recorder = &_MockOperatorRecorder{mock}
	return mock
}

func (_m *MockOperator) EXPECT() *_MockOperatorRecorder {
	return _m.recorder
}

func (_m *MockOperator) DeregisterListener(_param0 ListenerID) {
	_m.ctrl.Call(_m, "DeregisterListener", _param0)
}

func (_mr *_MockOperatorRecorder) DeregisterListener(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeregisterListener", arg0)
}

func (_m *MockOperator) RegisterListener(_param0 Listener) ListenerID {
	ret := _m.ctrl.Call(_m, "RegisterListener", _param0)
	ret0, _ := ret[0].(ListenerID)
	return ret0
}

func (_mr *_MockOperatorRecorder) RegisterListener(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RegisterListener", arg0)
}

func (_m *MockOperator) Reset() error {
	ret := _m.ctrl.Call(_m, "Reset")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperatorRecorder) Reset() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Reset")
}

func (_m *MockOperator) Setup(_param0 build.ServiceBuild, _param1 build.ServiceConfiguration, _param2 string, _param3 bool) error {
	ret := _m.ctrl.Call(_m, "Setup", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperatorRecorder) Setup(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Setup", arg0, arg1, arg2, arg3)
}

func (_m *MockOperator) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperatorRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

func (_m *MockOperator) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperatorRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}

func (_m *MockOperator) Teardown() error {
	ret := _m.ctrl.Call(_m, "Teardown")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOperatorRecorder) Teardown() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Teardown")
}
