// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jbruno/dumbwaiter/common (interfaces: RPi)

// Package mock_common is a generated GoMock package.
package mock_common

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/jbruno/dumbwaiter/common"
	reflect "reflect"
)

// MockRPi is a mock of RPi interface
type MockRPi struct {
	ctrl     *gomock.Controller
	recorder *MockRPiMockRecorder
}

// MockRPiMockRecorder is the mock recorder for MockRPi
type MockRPiMockRecorder struct {
	mock *MockRPi
}

// NewMockRPi creates a new mock instance
func NewMockRPi(ctrl *gomock.Controller) *MockRPi {
	mock := &MockRPi{ctrl: ctrl}
	mock.recorder = &MockRPiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRPi) EXPECT() *MockRPiMockRecorder {
	return m.recorder
}

// GetSignal mocks base method
func (m *MockRPi) GetSignal(arg0 common.PiPin) (bool, error) {
	ret := m.ctrl.Call(m, "GetSignal", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSignal indicates an expected call of GetSignal
func (mr *MockRPiMockRecorder) GetSignal(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSignal", reflect.TypeOf((*MockRPi)(nil).GetSignal), arg0)
}

// SendSignal mocks base method
func (m *MockRPi) SendSignal(arg0 common.PiPin) error {
	ret := m.ctrl.Call(m, "SendSignal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSignal indicates an expected call of SendSignal
func (mr *MockRPiMockRecorder) SendSignal(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSignal", reflect.TypeOf((*MockRPi)(nil).SendSignal), arg0)
}
