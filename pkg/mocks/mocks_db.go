// Code generated by MockGen. DO NOT EDIT.
// Source: database/cdxj.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	gomock "github.com/golang/mock/gomock"
)

// MockCdxjRecords is a mock of CdxjRecords interface.
type MockCdxjRecords struct {
	ctrl     *gomock.Controller
	recorder *MockCdxjRecordsMockRecorder
}

// MockCdxjRecordsMockRecorder is the mock recorder for MockCdxjRecords.
type MockCdxjRecordsMockRecorder struct {
	mock *MockCdxjRecords
}

// NewMockCdxjRecords creates a new mock instance.
func NewMockCdxjRecords(ctrl *gomock.Controller) *MockCdxjRecords {
	mock := &MockCdxjRecords{ctrl: ctrl}
	mock.recorder = &MockCdxjRecordsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCdxjRecords) EXPECT() *MockCdxjRecordsMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockCdxjRecords) Find(surt string, types []string, before, after int64, match string, limit int) (models.Records, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", surt, types, before, after, match, limit)
	ret0, _ := ret[0].(models.Records)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockCdxjRecordsMockRecorder) Find(surt, types, before, after, match, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCdxjRecords)(nil).Find), surt, types, before, after, match, limit)
}

// Get mocks base method.
func (m *MockCdxjRecords) Get(recordId string) (models.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", recordId)
	ret0, _ := ret[0].(models.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCdxjRecordsMockRecorder) Get(recordId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCdxjRecords)(nil).Get), recordId)
}

// Set mocks base method.
func (m *MockCdxjRecords) Set(rec models.Record) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", rec)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockCdxjRecordsMockRecorder) Set(rec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCdxjRecords)(nil).Set), rec)
}
