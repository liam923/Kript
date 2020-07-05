// Code generated by MockGen. DO NOT EDIT.
// Source: firebase.go

// Package data is a generated GoMock package.
package data

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mockdatabase is a mock of database interface
type Mockdatabase struct {
	ctrl     *gomock.Controller
	recorder *MockdatabaseMockRecorder
}

// MockdatabaseMockRecorder is the mock recorder for Mockdatabase
type MockdatabaseMockRecorder struct {
	mock *Mockdatabase
}

// NewMockdatabase creates a new mock instance
func NewMockdatabase(ctrl *gomock.Controller) *Mockdatabase {
	mock := &Mockdatabase{ctrl: ctrl}
	mock.recorder = &MockdatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockdatabase) EXPECT() *MockdatabaseMockRecorder {
	return m.recorder
}

// fetchDatum mocks base method
func (m *Mockdatabase) fetchDatum(ctx context.Context, id string) (*datum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fetchDatum", ctx, id)
	ret0, _ := ret[0].(*datum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// fetchDatum indicates an expected call of fetchDatum
func (mr *MockdatabaseMockRecorder) fetchDatum(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fetchDatum", reflect.TypeOf((*Mockdatabase)(nil).fetchDatum), ctx, id)
}

// fetchDataForUser mocks base method
func (m *Mockdatabase) fetchDataForUser(ctx context.Context, userId string) (*[]idedDatum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fetchDataForUser", ctx, userId)
	ret0, _ := ret[0].(*[]idedDatum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// fetchDataForUser indicates an expected call of fetchDataForUser
func (mr *MockdatabaseMockRecorder) fetchDataForUser(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fetchDataForUser", reflect.TypeOf((*Mockdatabase)(nil).fetchDataForUser), ctx, userId)
}

// createDatum mocks base method
func (m *Mockdatabase) createDatum(ctx context.Context, datum *datum) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createDatum", ctx, datum)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createDatum indicates an expected call of createDatum
func (mr *MockdatabaseMockRecorder) createDatum(ctx, datum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createDatum", reflect.TypeOf((*Mockdatabase)(nil).createDatum), ctx, datum)
}

// updateDatum mocks base method
func (m *Mockdatabase) updateDatum(ctx context.Context, datum *datum, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "updateDatum", ctx, datum, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// updateDatum indicates an expected call of updateDatum
func (mr *MockdatabaseMockRecorder) updateDatum(ctx, datum, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "updateDatum", reflect.TypeOf((*Mockdatabase)(nil).updateDatum), ctx, datum, id)
}

// deleteDatum mocks base method
func (m *Mockdatabase) deleteDatum(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "deleteDatum", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// deleteDatum indicates an expected call of deleteDatum
func (mr *MockdatabaseMockRecorder) deleteDatum(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteDatum", reflect.TypeOf((*Mockdatabase)(nil).deleteDatum), ctx, id)
}
