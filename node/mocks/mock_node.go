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
// Source: github.com/m3db/m3em/node (interfaces: ServiceNode,Options)

package node

import (
	gomock "github.com/golang/mock/gomock"
	build "github.com/m3db/m3em/build"
	node "github.com/m3db/m3em/node"
	placementpb "github.com/m3db/m3cluster/generated/proto/placementpb"
	placement "github.com/m3db/m3cluster/placement"
	shard "github.com/m3db/m3cluster/shard"
	instrument "github.com/m3db/m3x/instrument"
	retry "github.com/m3db/m3x/retry"
	time "time"
)

// Mock of ServiceNode interface
type MockServiceNode struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceNodeRecorder
}

// Recorder for MockServiceNode (not exported)
type _MockServiceNodeRecorder struct {
	mock *MockServiceNode
}

func NewMockServiceNode(ctrl *gomock.Controller) *MockServiceNode {
	mock := &MockServiceNode{ctrl: ctrl}
	mock.recorder = &_MockServiceNodeRecorder{mock}
	return mock
}

func (_m *MockServiceNode) EXPECT() *_MockServiceNodeRecorder {
	return _m.recorder
}

func (_m *MockServiceNode) Clone() placement.Instance {
	ret := _m.ctrl.Call(_m, "Clone")
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Clone() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Clone")
}

func (_m *MockServiceNode) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockServiceNode) DeregisterListener(_param0 node.ListenerID) {
	_m.ctrl.Call(_m, "DeregisterListener", _param0)
}

func (_mr *_MockServiceNodeRecorder) DeregisterListener(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeregisterListener", arg0)
}

func (_m *MockServiceNode) Endpoint() string {
	ret := _m.ctrl.Call(_m, "Endpoint")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Endpoint() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Endpoint")
}

func (_m *MockServiceNode) GetRemoteOutput(_param0 node.RemoteOutputType, _param1 string) (bool, error) {
	ret := _m.ctrl.Call(_m, "GetRemoteOutput", _param0, _param1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceNodeRecorder) GetRemoteOutput(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRemoteOutput", arg0, arg1)
}

func (_m *MockServiceNode) Hostname() string {
	ret := _m.ctrl.Call(_m, "Hostname")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Hostname() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Hostname")
}

func (_m *MockServiceNode) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

func (_m *MockServiceNode) IsInitializing() bool {
	ret := _m.ctrl.Call(_m, "IsInitializing")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) IsInitializing() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsInitializing")
}

func (_m *MockServiceNode) IsLeaving() bool {
	ret := _m.ctrl.Call(_m, "IsLeaving")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) IsLeaving() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsLeaving")
}

func (_m *MockServiceNode) Port() uint32 {
	ret := _m.ctrl.Call(_m, "Port")
	ret0, _ := ret[0].(uint32)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Port() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Port")
}

func (_m *MockServiceNode) Proto() (*placementpb.Instance, error) {
	ret := _m.ctrl.Call(_m, "Proto")
	ret0, _ := ret[0].(*placementpb.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceNodeRecorder) Proto() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Proto")
}

func (_m *MockServiceNode) Rack() string {
	ret := _m.ctrl.Call(_m, "Rack")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Rack() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Rack")
}

func (_m *MockServiceNode) RegisterListener(_param0 node.Listener) node.ListenerID {
	ret := _m.ctrl.Call(_m, "RegisterListener", _param0)
	ret0, _ := ret[0].(node.ListenerID)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) RegisterListener(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RegisterListener", arg0)
}

func (_m *MockServiceNode) SetEndpoint(_param0 string) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetEndpoint", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetEndpoint(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetEndpoint", arg0)
}

func (_m *MockServiceNode) SetHostname(_param0 string) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetHostname", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetHostname(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetHostname", arg0)
}

func (_m *MockServiceNode) SetID(_param0 string) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetID", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetID", arg0)
}

func (_m *MockServiceNode) SetPort(_param0 uint32) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetPort", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetPort(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetPort", arg0)
}

func (_m *MockServiceNode) SetRack(_param0 string) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetRack", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetRack(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetRack", arg0)
}

func (_m *MockServiceNode) SetShardSetID(_param0 uint32) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetShardSetID", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetShardSetID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetShardSetID", arg0)
}

func (_m *MockServiceNode) SetShards(_param0 shard.Shards) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetShards", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetShards(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetShards", arg0)
}

func (_m *MockServiceNode) SetWeight(_param0 uint32) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetWeight", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetWeight(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetWeight", arg0)
}

func (_m *MockServiceNode) SetZone(_param0 string) placement.Instance {
	ret := _m.ctrl.Call(_m, "SetZone", _param0)
	ret0, _ := ret[0].(placement.Instance)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) SetZone(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetZone", arg0)
}

func (_m *MockServiceNode) Setup(_param0 build.ServiceBuild, _param1 build.ServiceConfiguration, _param2 string, _param3 bool) error {
	ret := _m.ctrl.Call(_m, "Setup", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Setup(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Setup", arg0, arg1, arg2, arg3)
}

func (_m *MockServiceNode) ShardSetID() uint32 {
	ret := _m.ctrl.Call(_m, "ShardSetID")
	ret0, _ := ret[0].(uint32)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) ShardSetID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ShardSetID")
}

func (_m *MockServiceNode) Shards() shard.Shards {
	ret := _m.ctrl.Call(_m, "Shards")
	ret0, _ := ret[0].(shard.Shards)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Shards() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Shards")
}

func (_m *MockServiceNode) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

func (_m *MockServiceNode) Status() node.Status {
	ret := _m.ctrl.Call(_m, "Status")
	ret0, _ := ret[0].(node.Status)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Status() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Status")
}

func (_m *MockServiceNode) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}

func (_m *MockServiceNode) String() string {
	ret := _m.ctrl.Call(_m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) String() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "String")
}

func (_m *MockServiceNode) Teardown() error {
	ret := _m.ctrl.Call(_m, "Teardown")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Teardown() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Teardown")
}

func (_m *MockServiceNode) TransferLocalFile(_param0 string, _param1 []string, _param2 bool) error {
	ret := _m.ctrl.Call(_m, "TransferLocalFile", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) TransferLocalFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TransferLocalFile", arg0, arg1, arg2)
}

func (_m *MockServiceNode) Weight() uint32 {
	ret := _m.ctrl.Call(_m, "Weight")
	ret0, _ := ret[0].(uint32)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Weight() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Weight")
}

func (_m *MockServiceNode) Zone() string {
	ret := _m.ctrl.Call(_m, "Zone")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceNodeRecorder) Zone() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Zone")
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

func (_m *MockOptions) HeartbeatOptions() node.HeartbeatOptions {
	ret := _m.ctrl.Call(_m, "HeartbeatOptions")
	ret0, _ := ret[0].(node.HeartbeatOptions)
	return ret0
}

func (_mr *_MockOptionsRecorder) HeartbeatOptions() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HeartbeatOptions")
}

func (_m *MockOptions) InstrumentOptions() instrument.Options {
	ret := _m.ctrl.Call(_m, "InstrumentOptions")
	ret0, _ := ret[0].(instrument.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) InstrumentOptions() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InstrumentOptions")
}

func (_m *MockOptions) MaxPullSize() int64 {
	ret := _m.ctrl.Call(_m, "MaxPullSize")
	ret0, _ := ret[0].(int64)
	return ret0
}

func (_mr *_MockOptionsRecorder) MaxPullSize() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MaxPullSize")
}

func (_m *MockOptions) OperationTimeout() time.Duration {
	ret := _m.ctrl.Call(_m, "OperationTimeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

func (_mr *_MockOptionsRecorder) OperationTimeout() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "OperationTimeout")
}

func (_m *MockOptions) OperatorClientFn() node.OperatorClientFn {
	ret := _m.ctrl.Call(_m, "OperatorClientFn")
	ret0, _ := ret[0].(node.OperatorClientFn)
	return ret0
}

func (_mr *_MockOptionsRecorder) OperatorClientFn() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "OperatorClientFn")
}

func (_m *MockOptions) Retrier() retry.Retrier {
	ret := _m.ctrl.Call(_m, "Retrier")
	ret0, _ := ret[0].(retry.Retrier)
	return ret0
}

func (_mr *_MockOptionsRecorder) Retrier() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Retrier")
}

func (_m *MockOptions) SetHeartbeatOptions(_param0 node.HeartbeatOptions) node.Options {
	ret := _m.ctrl.Call(_m, "SetHeartbeatOptions", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetHeartbeatOptions(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetHeartbeatOptions", arg0)
}

func (_m *MockOptions) SetInstrumentOptions(_param0 instrument.Options) node.Options {
	ret := _m.ctrl.Call(_m, "SetInstrumentOptions", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetInstrumentOptions(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetInstrumentOptions", arg0)
}

func (_m *MockOptions) SetMaxPullSize(_param0 int64) node.Options {
	ret := _m.ctrl.Call(_m, "SetMaxPullSize", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetMaxPullSize(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetMaxPullSize", arg0)
}

func (_m *MockOptions) SetOperationTimeout(_param0 time.Duration) node.Options {
	ret := _m.ctrl.Call(_m, "SetOperationTimeout", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetOperationTimeout(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetOperationTimeout", arg0)
}

func (_m *MockOptions) SetOperatorClientFn(_param0 node.OperatorClientFn) node.Options {
	ret := _m.ctrl.Call(_m, "SetOperatorClientFn", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetOperatorClientFn(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetOperatorClientFn", arg0)
}

func (_m *MockOptions) SetRetrier(_param0 retry.Retrier) node.Options {
	ret := _m.ctrl.Call(_m, "SetRetrier", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetRetrier(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetRetrier", arg0)
}

func (_m *MockOptions) SetTransferBufferSize(_param0 int) node.Options {
	ret := _m.ctrl.Call(_m, "SetTransferBufferSize", _param0)
	ret0, _ := ret[0].(node.Options)
	return ret0
}

func (_mr *_MockOptionsRecorder) SetTransferBufferSize(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetTransferBufferSize", arg0)
}

func (_m *MockOptions) TransferBufferSize() int {
	ret := _m.ctrl.Call(_m, "TransferBufferSize")
	ret0, _ := ret[0].(int)
	return ret0
}

func (_mr *_MockOptionsRecorder) TransferBufferSize() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TransferBufferSize")
}

func (_m *MockOptions) Validate() error {
	ret := _m.ctrl.Call(_m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOptionsRecorder) Validate() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Validate")
}
