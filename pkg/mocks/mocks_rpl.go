// Code generated by MockGen. DO NOT EDIT.
// Source: replayer/replayer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/crossedbot/hermes-archiver/pkg/replayer/models"
	gomock "github.com/golang/mock/gomock"
)

// MockReplayer is a mock of Replayer interface.
type MockReplayer struct {
	ctrl     *gomock.Controller
	recorder *MockReplayerMockRecorder
}

// MockReplayerMockRecorder is the mock recorder for MockReplayer.
type MockReplayerMockRecorder struct {
	mock *MockReplayer
}

// NewMockReplayer creates a new mock instance.
func NewMockReplayer(ctrl *gomock.Controller) *MockReplayer {
	mock := &MockReplayer{ctrl: ctrl}
	mock.recorder = &MockReplayerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReplayer) EXPECT() *MockReplayerMockRecorder {
	return m.recorder
}

// Replay mocks base method.
func (m *MockReplayer) Replay(id string, key []byte) (models.Replay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replay", id, key)
	ret0, _ := ret[0].(models.Replay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Replay indicates an expected call of Replay.
func (mr *MockReplayerMockRecorder) Replay(id, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replay", reflect.TypeOf((*MockReplayer)(nil).Replay), id, key)
}
