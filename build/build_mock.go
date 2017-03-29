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
// Source: /Users/prungta/code/gocode/src/github.com/m3db/m3em/build/types.go

package build

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of ServiceBuild interface
type MockServiceBuild struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceBuildRecorder
}

// Recorder for MockServiceBuild (not exported)
type _MockServiceBuildRecorder struct {
	mock *MockServiceBuild
}

func NewMockServiceBuild(ctrl *gomock.Controller) *MockServiceBuild {
	mock := &MockServiceBuild{ctrl: ctrl}
	mock.recorder = &_MockServiceBuildRecorder{mock}
	return mock
}

func (_m *MockServiceBuild) EXPECT() *_MockServiceBuildRecorder {
	return _m.recorder
}

func (_m *MockServiceBuild) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceBuildRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

func (_m *MockServiceBuild) SourcePath() string {
	ret := _m.ctrl.Call(_m, "SourcePath")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceBuildRecorder) SourcePath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SourcePath")
}

// Mock of ServiceConfiguration interface
type MockServiceConfiguration struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceConfigurationRecorder
}

// Recorder for MockServiceConfiguration (not exported)
type _MockServiceConfigurationRecorder struct {
	mock *MockServiceConfiguration
}

func NewMockServiceConfiguration(ctrl *gomock.Controller) *MockServiceConfiguration {
	mock := &MockServiceConfiguration{ctrl: ctrl}
	mock.recorder = &_MockServiceConfigurationRecorder{mock}
	return mock
}

func (_m *MockServiceConfiguration) EXPECT() *_MockServiceConfigurationRecorder {
	return _m.recorder
}

func (_m *MockServiceConfiguration) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceConfigurationRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

func (_m *MockServiceConfiguration) MarshalText() ([]byte, error) {
	ret := _m.ctrl.Call(_m, "MarshalText")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceConfigurationRecorder) MarshalText() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MarshalText")
}

func (_m *MockServiceConfiguration) UnmarshalText(_param0 []byte) error {
	ret := _m.ctrl.Call(_m, "UnmarshalText", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceConfigurationRecorder) UnmarshalText(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UnmarshalText", arg0)
}
