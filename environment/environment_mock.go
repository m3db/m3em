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
// Source: github.com/m3db/m3em/environment (interfaces: M3DBInstance,M3DBEnvironment,Options)

package environment

import (
	gomock "github.com/golang/mock/gomock"
	build "github.com/m3db/m3em/build"

	operator "github.com/m3db/m3em/operator"
	services "github.com/m3db/m3cluster/services"
	shard "github.com/m3db/m3cluster/shard"
	instrument "github.com/m3db/m3x/instrument"
	retry "github.com/m3db/m3x/retry"
	time "time"
)

// Mock of M3DBInstance interface
type MockM3DBInstance struct {
	ctrl     *gomock.Controller
	recorder *_MockM3DBInstanceRecorder
}

// Recorder for MockM3DBInstance (not exported)
type _MockM3DBInstanceRecorder struct {
	mock *MockM3DBInstance
}

func NewMockM3DBInstance(ctrl *gomock.Controller) *MockM3DBInstance {
	mock := &MockM3DBInstance{ctrl: ctrl}
	mock.recorder = &_MockM3DBInstanceRecorder{mock}
	return mock
}

func (_m *MockM3DBInstance) EXPECT() *_MockM3DBInstanceRecorder {
	return _m.recorder
}

func (_m *MockM3DBInstance) Endpoint() string {
	ret := _m.ctrl.Call(_m, "Endpoint")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Endpoint() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Endpoint")
}

func (_m *MockM3DBInstance) Health() (M3DBInstanceHealth, error) {
	ret := _m.ctrl.Call(_m, "Health")
	ret0, _ := ret[0].(M3DBInstanceHealth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockM3DBInstanceRecorder) Health() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Health")
}

func (_m *MockM3DBInstance) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

func (_m *MockM3DBInstance) Operator() operator.Operator {
	ret := _m.ctrl.Call(_m, "Operator")
	ret0, _ := ret[0].(operator.Operator)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Operator() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Operator")
}

func (_m *MockM3DBInstance) OverrideConfiguration(_param0 build.ServiceConfiguration) error {
	ret := _m.ctrl.Call(_m, "OverrideConfiguration", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) OverrideConfiguration(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "OverrideConfiguration", arg0)
}

func (_m *MockM3DBInstance) Rack() string {
	ret := _m.ctrl.Call(_m, "Rack")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Rack() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Rack")
}

func (_m *MockM3DBInstance) Reset() error {
	ret := _m.ctrl.Call(_m, "Reset")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Reset() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Reset")
}

func (_m *MockM3DBInstance) SetEndpoint(_param0 string) services.PlacementInstance {
	ret := _m.ctrl.Call(_m, "SetEndpoint", _param0)
	ret0, _ := ret[0].(services.PlacementInstance)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) SetEndpoint(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetEndpoint", arg0)
}

func (_m *MockM3DBInstance) SetID(_param0 string) services.PlacementInstance {
	ret := _m.ctrl.Call(_m, "SetID", _param0)
	ret0, _ := ret[0].(services.PlacementInstance)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) SetID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetID", arg0)
}

func (_m *MockM3DBInstance) SetRack(_param0 string) services.PlacementInstance {
	ret := _m.ctrl.Call(_m, "SetRack", _param0)
	ret0, _ := ret[0].(services.PlacementInstance)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) SetRack(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetRack", arg0)
}

func (_m *MockM3DBInstance) SetShards(_param0 shard.Shards) services.PlacementInstance {
	ret := _m.ctrl.Call(_m, "SetShards", _param0)
	ret0, _ := ret[0].(services.PlacementInstance)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) SetShards(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetShards", arg0)
}

func (_m *MockM3DBInstance) SetWeight(_param0 uint32) services.PlacementInstance {
	ret := _m.ctrl.Call(_m, "SetWeight", _param0)
	ret0, _ := ret[0].(services.PlacementInstance)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) SetWeight(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetWeight", arg0)
}

func (_m *MockM3DBInstance) SetZone(_param0 string) services.PlacementInstance {
	ret := _m.ctrl.Call(_m, "SetZone", _param0)
	ret0, _ := ret[0].(services.PlacementInstance)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) SetZone(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetZone", arg0)
}

func (_m *MockM3DBInstance) Setup(_param0 build.ServiceBuild, _param1 build.ServiceConfiguration) error {
	ret := _m.ctrl.Call(_m, "Setup", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Setup(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Setup", arg0, arg1)
}

func (_m *MockM3DBInstance) Shards() shard.Shards {
	ret := _m.ctrl.Call(_m, "Shards")
	ret0, _ := ret[0].(shard.Shards)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Shards() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Shards")
}

func (_m *MockM3DBInstance) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

func (_m *MockM3DBInstance) Status() InstanceStatus {
	ret := _m.ctrl.Call(_m, "Status")
	ret0, _ := ret[0].(InstanceStatus)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Status() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Status")
}

func (_m *MockM3DBInstance) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}

func (_m *MockM3DBInstance) String() string {
	ret := _m.ctrl.Call(_m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) String() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "String")
}

func (_m *MockM3DBInstance) Teardown() error {
	ret := _m.ctrl.Call(_m, "Teardown")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Teardown() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Teardown")
}

func (_m *MockM3DBInstance) Weight() uint32 {
	ret := _m.ctrl.Call(_m, "Weight")
	ret0, _ := ret[0].(uint32)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Weight() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Weight")
}

func (_m *MockM3DBInstance) Zone() string {
	ret := _m.ctrl.Call(_m, "Zone")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockM3DBInstanceRecorder) Zone() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Zone")
}

// Mock of M3DBEnvironment interface
type MockM3DBEnvironment struct {
	ctrl     *gomock.Controller
	recorder *_MockM3DBEnvironmentRecorder
}

// Recorder for MockM3DBEnvironment (not exported)
type _MockM3DBEnvironmentRecorder struct {
	mock *MockM3DBEnvironment
}

func NewMockM3DBEnvironment(ctrl *gomock.Controller) *MockM3DBEnvironment {
	mock := &MockM3DBEnvironment{ctrl: ctrl}
	mock.recorder = &_MockM3DBEnvironmentRecorder{mock}
	return mock
}

func (_m *MockM3DBEnvironment) EXPECT() *_MockM3DBEnvironmentRecorder {
	return _m.recorder
}

func (_m *MockM3DBEnvironment) Instances() M3DBInstances {
	ret := _m.ctrl.Call(_m, "Instances")
	ret0, _ := ret[0].(M3DBInstances)
	return ret0
}

func (_mr *_MockM3DBEnvironmentRecorder) Instances() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Instances")
}

func (_m *MockM3DBEnvironment) InstancesByID() map[string]M3DBInstance {
	ret := _m.ctrl.Call(_m, "InstancesByID")
	ret0, _ := ret[0].(map[string]M3DBInstance)
	return ret0
}

func (_mr *_MockM3DBEnvironmentRecorder) InstancesByID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InstancesByID")
}

func (_m *MockM3DBEnvironment) Status() map[string]InstanceStatus {
	ret := _m.ctrl.Call(_m, "Status")
	ret0, _ := ret[0].(map[string]InstanceStatus)
	return ret0
}

func (_mr *_MockM3DBEnvironmentRecorder) Status() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Status")
}

// Mock of Options interface
type MockOptions struct {
	ctrl     *gomock.Controller
	recorder *_MockOptionsRecorder
}

// Recorder for MockOptions (not exported)
type _MockOptionsRecorder struct {
	mock *MockOptions
}

func NewMockOptions(ctrl *gomock.Controller) *MockOptions {
	mock := &MockOptions{ctrl: ctrl}
	mock.recorder = &_MockOptionsRecorder{mock}
	return mock
}

func (_m *MockOptions) EXPECT() *_MockOptionsRecorder {
	return _m.recorder
}

func (_m *MockOptions) InstanceOperationRetrier() retry.Retrier {
	ret := _m.ctrl.Call(_m, "InstanceOperationRetrier")
	ret0, _ := ret[0].(retry.Retrier)
	return ret0
}

func (_mr *_MockOptionsRecorder) InstanceOperationRetrier() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InstanceOperationRetrier")
}

func (_m *MockOptions) InstanceOperationTimeout() time.Duration {
	ret := _m.ctrl.Call(_m, "InstanceOperationTimeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

func (_mr *_MockOptionsRecorder) InstanceOperationTimeout() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InstanceOperationTimeout")
}

func (_m *MockOptions) InstrumentOptions() instrument.Options {
	ret := _m.ctrl.Call(_m, "InstrumentOptions")
	ret0, _ := ret[0].(instrument.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) InstrumentOptions() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InstrumentOptions")
}

func (_m *MockOptions) OperatorOptions() operator.Options {
	ret := _m.ctrl.Call(_m, "OperatorOptions")
	ret0, _ := ret[0].(operator.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) OperatorOptions() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "OperatorOptions")
}

func (_m *MockOptions) SessionOverride() bool {
	ret := _m.ctrl.Call(_m, "SessionOverride")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockOptionsRecorder) SessionOverride() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SessionOverride")
}

func (_m *MockOptions) SessionToken() string {
	ret := _m.ctrl.Call(_m, "SessionToken")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockOptionsRecorder) SessionToken() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SessionToken")
}

func (_m *MockOptions) SetInstanceOperationRetrier(_param0 retry.Retrier) Options {
	ret := _m.ctrl.Call(_m, "SetInstanceOperationRetrier", _param0)
	ret0, _ := ret[0].(Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetInstanceOperationRetrier(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetInstanceOperationRetrier", arg0)
}

func (_m *MockOptions) SetInstanceOperationTimeout(_param0 time.Duration) Options {
	ret := _m.ctrl.Call(_m, "SetInstanceOperationTimeout", _param0)
	ret0, _ := ret[0].(Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetInstanceOperationTimeout(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetInstanceOperationTimeout", arg0)
}

func (_m *MockOptions) SetInstrumentOptions(_param0 instrument.Options) Options {
	ret := _m.ctrl.Call(_m, "SetInstrumentOptions", _param0)
	ret0, _ := ret[0].(Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetInstrumentOptions(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetInstrumentOptions", arg0)
}

func (_m *MockOptions) SetOperatorOptions(_param0 operator.Options) Options {
	ret := _m.ctrl.Call(_m, "SetOperatorOptions", _param0)
	ret0, _ := ret[0].(Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetOperatorOptions(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetOperatorOptions", arg0)
}

func (_m *MockOptions) SetSessionOverride(_param0 bool) Options {
	ret := _m.ctrl.Call(_m, "SetSessionOverride", _param0)
	ret0, _ := ret[0].(Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetSessionOverride(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetSessionOverride", arg0)
}

func (_m *MockOptions) SetSessionToken(_param0 string) Options {
	ret := _m.ctrl.Call(_m, "SetSessionToken", _param0)
	ret0, _ := ret[0].(Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetSessionToken(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetSessionToken", arg0)
}
