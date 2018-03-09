// Copyright (c) 2018 Uber Technologies, Inc.
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

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3em/build/types.go

package build

import (
	gomock "github.com/golang/mock/gomock"
	fs "github.com/m3db/m3em/os/fs"
	reflect "reflect"
)

// MockIterableBytesWithID is a mock of IterableBytesWithID interface
type MockIterableBytesWithID struct {
	ctrl     *gomock.Controller
	recorder *MockIterableBytesWithIDMockRecorder
}

// MockIterableBytesWithIDMockRecorder is the mock recorder for MockIterableBytesWithID
type MockIterableBytesWithIDMockRecorder struct {
	mock *MockIterableBytesWithID
}

// NewMockIterableBytesWithID creates a new mock instance
func NewMockIterableBytesWithID(ctrl *gomock.Controller) *MockIterableBytesWithID {
	mock := &MockIterableBytesWithID{ctrl: ctrl}
	mock.recorder = &MockIterableBytesWithIDMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockIterableBytesWithID) EXPECT() *MockIterableBytesWithIDMockRecorder {
	return _m.recorder
}

// ID mocks base method
func (_m *MockIterableBytesWithID) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (_mr *MockIterableBytesWithIDMockRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "ID", reflect.TypeOf((*MockIterableBytesWithID)(nil).ID))
}

// Iter mocks base method
func (_m *MockIterableBytesWithID) Iter(bufferSize int) (fs.FileReaderIter, error) {
	ret := _m.ctrl.Call(_m, "Iter", bufferSize)
	ret0, _ := ret[0].(fs.FileReaderIter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Iter indicates an expected call of Iter
func (_mr *MockIterableBytesWithIDMockRecorder) Iter(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Iter", reflect.TypeOf((*MockIterableBytesWithID)(nil).Iter), arg0)
}

// MockServiceBuild is a mock of ServiceBuild interface
type MockServiceBuild struct {
	ctrl     *gomock.Controller
	recorder *MockServiceBuildMockRecorder
}

// MockServiceBuildMockRecorder is the mock recorder for MockServiceBuild
type MockServiceBuildMockRecorder struct {
	mock *MockServiceBuild
}

// NewMockServiceBuild creates a new mock instance
func NewMockServiceBuild(ctrl *gomock.Controller) *MockServiceBuild {
	mock := &MockServiceBuild{ctrl: ctrl}
	mock.recorder = &MockServiceBuildMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockServiceBuild) EXPECT() *MockServiceBuildMockRecorder {
	return _m.recorder
}

// ID mocks base method
func (_m *MockServiceBuild) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (_mr *MockServiceBuildMockRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "ID", reflect.TypeOf((*MockServiceBuild)(nil).ID))
}

// Iter mocks base method
func (_m *MockServiceBuild) Iter(bufferSize int) (fs.FileReaderIter, error) {
	ret := _m.ctrl.Call(_m, "Iter", bufferSize)
	ret0, _ := ret[0].(fs.FileReaderIter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Iter indicates an expected call of Iter
func (_mr *MockServiceBuildMockRecorder) Iter(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Iter", reflect.TypeOf((*MockServiceBuild)(nil).Iter), arg0)
}

// SourcePath mocks base method
func (_m *MockServiceBuild) SourcePath() string {
	ret := _m.ctrl.Call(_m, "SourcePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// SourcePath indicates an expected call of SourcePath
func (_mr *MockServiceBuildMockRecorder) SourcePath() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SourcePath", reflect.TypeOf((*MockServiceBuild)(nil).SourcePath))
}

// MockServiceConfiguration is a mock of ServiceConfiguration interface
type MockServiceConfiguration struct {
	ctrl     *gomock.Controller
	recorder *MockServiceConfigurationMockRecorder
}

// MockServiceConfigurationMockRecorder is the mock recorder for MockServiceConfiguration
type MockServiceConfigurationMockRecorder struct {
	mock *MockServiceConfiguration
}

// NewMockServiceConfiguration creates a new mock instance
func NewMockServiceConfiguration(ctrl *gomock.Controller) *MockServiceConfiguration {
	mock := &MockServiceConfiguration{ctrl: ctrl}
	mock.recorder = &MockServiceConfigurationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockServiceConfiguration) EXPECT() *MockServiceConfigurationMockRecorder {
	return _m.recorder
}

// ID mocks base method
func (_m *MockServiceConfiguration) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (_mr *MockServiceConfigurationMockRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "ID", reflect.TypeOf((*MockServiceConfiguration)(nil).ID))
}

// Iter mocks base method
func (_m *MockServiceConfiguration) Iter(bufferSize int) (fs.FileReaderIter, error) {
	ret := _m.ctrl.Call(_m, "Iter", bufferSize)
	ret0, _ := ret[0].(fs.FileReaderIter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Iter indicates an expected call of Iter
func (_mr *MockServiceConfigurationMockRecorder) Iter(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Iter", reflect.TypeOf((*MockServiceConfiguration)(nil).Iter), arg0)
}

// Bytes mocks base method
func (_m *MockServiceConfiguration) Bytes() ([]byte, error) {
	ret := _m.ctrl.Call(_m, "Bytes")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bytes indicates an expected call of Bytes
func (_mr *MockServiceConfigurationMockRecorder) Bytes() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Bytes", reflect.TypeOf((*MockServiceConfiguration)(nil).Bytes))
}
