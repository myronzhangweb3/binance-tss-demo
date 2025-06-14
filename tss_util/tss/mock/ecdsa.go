// Code generated by MockGen. DO NOT EDIT.
// Source: ./tss/ecdsa/keygen/keygen.go

// Package mock_tss is a generated GoMock package.
package mock_tss

import (
	reflect "reflect"
	"tss-demo/tss_util/keyshare"

	gomock "github.com/golang/mock/gomock"
)

// MockECDSAKeyshareStorer is a mock of ECDSAKeyshareStorer interface.
type MockECDSAKeyshareStorer struct {
	ctrl     *gomock.Controller
	recorder *MockECDSAKeyshareStorerMockRecorder
}

// MockECDSAKeyshareStorerMockRecorder is the mock recorder for MockECDSAKeyshareStorer.
type MockECDSAKeyshareStorerMockRecorder struct {
	mock *MockECDSAKeyshareStorer
}

// NewMockECDSAKeyshareStorer creates a new mock instance.
func NewMockECDSAKeyshareStorer(ctrl *gomock.Controller) *MockECDSAKeyshareStorer {
	mock := &MockECDSAKeyshareStorer{ctrl: ctrl}
	mock.recorder = &MockECDSAKeyshareStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockECDSAKeyshareStorer) EXPECT() *MockECDSAKeyshareStorerMockRecorder {
	return m.recorder
}

// GetKeyshare mocks base method.
func (m *MockECDSAKeyshareStorer) GetKeyshare() (keyshare.ECDSAKeyshare, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeyshare")
	ret0, _ := ret[0].(keyshare.ECDSAKeyshare)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeyshare indicates an expected call of GetKeyshare.
func (mr *MockECDSAKeyshareStorerMockRecorder) GetKeyshare() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeyshare", reflect.TypeOf((*MockECDSAKeyshareStorer)(nil).GetKeyshare))
}

// LockKeyshare mocks base method.
func (m *MockECDSAKeyshareStorer) LockKeyshare() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LockKeyshare")
}

// LockKeyshare indicates an expected call of LockKeyshare.
func (mr *MockECDSAKeyshareStorerMockRecorder) LockKeyshare() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockKeyshare", reflect.TypeOf((*MockECDSAKeyshareStorer)(nil).LockKeyshare))
}

// StoreKeyshare mocks base method.
func (m *MockECDSAKeyshareStorer) StoreKeyshare(keyshare keyshare.ECDSAKeyshare) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreKeyshare", keyshare)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreKeyshare indicates an expected call of StoreKeyshare.
func (mr *MockECDSAKeyshareStorerMockRecorder) StoreKeyshare(keyshare interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreKeyshare", reflect.TypeOf((*MockECDSAKeyshareStorer)(nil).StoreKeyshare), keyshare)
}

// UnlockKeyshare mocks base method.
func (m *MockECDSAKeyshareStorer) UnlockKeyshare() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnlockKeyshare")
}

// UnlockKeyshare indicates an expected call of UnlockKeyshare.
func (mr *MockECDSAKeyshareStorerMockRecorder) UnlockKeyshare() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockKeyshare", reflect.TypeOf((*MockECDSAKeyshareStorer)(nil).UnlockKeyshare))
}
