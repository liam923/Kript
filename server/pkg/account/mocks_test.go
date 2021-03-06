// Code generated by MockGen. DO NOT EDIT.
// Source: firebase.go

// Package account is a generated GoMock package.
package account

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

// fetchUserById mocks base method
func (m *Mockdatabase) fetchUserById(ctx context.Context, userId string) (*user, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fetchUserById", ctx, userId)
	ret0, _ := ret[0].(*user)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// fetchUserById indicates an expected call of fetchUserById
func (mr *MockdatabaseMockRecorder) fetchUserById(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fetchUserById", reflect.TypeOf((*Mockdatabase)(nil).fetchUserById), ctx, userId)
}

// fetchUserByUsername mocks base method
func (m *Mockdatabase) fetchUserByUsername(ctx context.Context, username string) (*user, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fetchUserByUsername", ctx, username)
	ret0, _ := ret[0].(*user)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// fetchUserByUsername indicates an expected call of fetchUserByUsername
func (mr *MockdatabaseMockRecorder) fetchUserByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fetchUserByUsername", reflect.TypeOf((*Mockdatabase)(nil).fetchUserByUsername), ctx, username)
}

// isUsernameAvailable mocks base method
func (m *Mockdatabase) isUsernameAvailable(ctx context.Context, username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isUsernameAvailable", ctx, username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// isUsernameAvailable indicates an expected call of isUsernameAvailable
func (mr *MockdatabaseMockRecorder) isUsernameAvailable(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isUsernameAvailable", reflect.TypeOf((*Mockdatabase)(nil).isUsernameAvailable), ctx, username)
}

// updateUser mocks base method
func (m *Mockdatabase) updateUser(ctx context.Context, userId string, user *user) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "updateUser", ctx, userId, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// updateUser indicates an expected call of updateUser
func (mr *MockdatabaseMockRecorder) updateUser(ctx, userId, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "updateUser", reflect.TypeOf((*Mockdatabase)(nil).updateUser), ctx, userId, user)
}

// createUser mocks base method
func (m *Mockdatabase) createUser(ctx context.Context, user *user) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createUser", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createUser indicates an expected call of createUser
func (mr *MockdatabaseMockRecorder) createUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createUser", reflect.TypeOf((*Mockdatabase)(nil).createUser), ctx, user)
}

// addVerificationTokenCode mocks base method
func (m *Mockdatabase) addVerificationTokenCode(ctx context.Context, userId, tokenId, code string, confirmDestination *twoFactorOption) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "addVerificationTokenCode", ctx, userId, tokenId, code, confirmDestination)
	ret0, _ := ret[0].(error)
	return ret0
}

// addVerificationTokenCode indicates an expected call of addVerificationTokenCode
func (mr *MockdatabaseMockRecorder) addVerificationTokenCode(ctx, userId, tokenId, code, confirmDestination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "addVerificationTokenCode", reflect.TypeOf((*Mockdatabase)(nil).addVerificationTokenCode), ctx, userId, tokenId, code, confirmDestination)
}

// verifyVerificationTokenCode mocks base method
func (m *Mockdatabase) verifyVerificationTokenCode(ctx context.Context, userId, tokenId, code string) (*twoFactorOption, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "verifyVerificationTokenCode", ctx, userId, tokenId, code)
	ret0, _ := ret[0].(*twoFactorOption)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// verifyVerificationTokenCode indicates an expected call of verifyVerificationTokenCode
func (mr *MockdatabaseMockRecorder) verifyVerificationTokenCode(ctx, userId, tokenId, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "verifyVerificationTokenCode", reflect.TypeOf((*Mockdatabase)(nil).verifyVerificationTokenCode), ctx, userId, tokenId, code)
}
