// Code generated by MockGen. DO NOT EDIT.
// Source: ocp-video-api/internal/producer (interfaces: Producer)

// Package mock is a generated GoMock package.
package mock

import (
	producer "ocp-video-api/internal/producer"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockProducer) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockProducerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockProducer)(nil).Close))
}

// Init mocks base method.
func (m *MockProducer) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockProducerMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockProducer)(nil).Init))
}

// SendEvent mocks base method.
func (m *MockProducer) SendEvent(arg0 producer.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEvent indicates an expected call of SendEvent.
func (mr *MockProducerMockRecorder) SendEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEvent", reflect.TypeOf((*MockProducer)(nil).SendEvent), arg0)
}