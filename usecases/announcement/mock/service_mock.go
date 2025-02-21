// Code generated by MockGen. DO NOT EDIT.
// Source: usecases/announcement/interface.go
//
// Generated by this command:
//
//	mockgen -source=usecases/announcement/interface.go -destination=usecases/announcement/mock/service_mock.go
//

// Package mock_announcement is a generated GoMock package.
package mock_announcement

import (
	reflect "reflect"

	announcement "github.com/Vractos/kloni/usecases/announcement"
	common "github.com/Vractos/kloni/usecases/common"
	store "github.com/Vractos/kloni/usecases/store"
	gomock "go.uber.org/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// CloneAnnouncement mocks base method.
func (m *MockUseCase) CloneAnnouncement(input announcement.CloneAnnouncementDtoInput, credentials *[]store.Credentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloneAnnouncement", input, credentials)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloneAnnouncement indicates an expected call of CloneAnnouncement.
func (mr *MockUseCaseMockRecorder) CloneAnnouncement(input, credentials any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloneAnnouncement", reflect.TypeOf((*MockUseCase)(nil).CloneAnnouncement), input, credentials)
}

// ImportAnnouncement mocks base method.
func (m *MockUseCase) ImportAnnouncement(input announcement.ImportAnnouncementDtoInput, credentials *[]store.Credentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportAnnouncement", input, credentials)
	ret0, _ := ret[0].(error)
	return ret0
}

// ImportAnnouncement indicates an expected call of ImportAnnouncement.
func (mr *MockUseCaseMockRecorder) ImportAnnouncement(input, credentials any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportAnnouncement", reflect.TypeOf((*MockUseCase)(nil).ImportAnnouncement), input, credentials)
}

// RetrieveAnnouncements mocks base method.
func (m *MockUseCase) RetrieveAnnouncements(sku string, credentials store.Credentials) (*[]common.MeliAnnouncement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveAnnouncements", sku, credentials)
	ret0, _ := ret[0].(*[]common.MeliAnnouncement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveAnnouncements indicates an expected call of RetrieveAnnouncements.
func (mr *MockUseCaseMockRecorder) RetrieveAnnouncements(sku, credentials any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveAnnouncements", reflect.TypeOf((*MockUseCase)(nil).RetrieveAnnouncements), sku, credentials)
}

// RetrieveAnnouncementsFromAllAccounts mocks base method.
func (m *MockUseCase) RetrieveAnnouncementsFromAllAccounts(sku string, credentials *[]store.Credentials) (*[]announcement.Announcements, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveAnnouncementsFromAllAccounts", sku, credentials)
	ret0, _ := ret[0].(*[]announcement.Announcements)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveAnnouncementsFromAllAccounts indicates an expected call of RetrieveAnnouncementsFromAllAccounts.
func (mr *MockUseCaseMockRecorder) RetrieveAnnouncementsFromAllAccounts(sku, credentials any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveAnnouncementsFromAllAccounts", reflect.TypeOf((*MockUseCase)(nil).RetrieveAnnouncementsFromAllAccounts), sku, credentials)
}

// UpdateQuantity mocks base method.
func (m *MockUseCase) UpdateQuantity(id string, quantity int, credentials store.Credentials, variationIDs ...int) error {
	m.ctrl.T.Helper()
	varargs := []any{id, quantity, credentials}
	for _, a := range variationIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateQuantity", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuantity indicates an expected call of UpdateQuantity.
func (mr *MockUseCaseMockRecorder) UpdateQuantity(id, quantity, credentials any, variationIDs ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{id, quantity, credentials}, variationIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuantity", reflect.TypeOf((*MockUseCase)(nil).UpdateQuantity), varargs...)
}
