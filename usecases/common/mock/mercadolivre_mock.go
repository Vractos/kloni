// Code generated by MockGen. DO NOT EDIT.
// Source: usecases/common/mercadolivre.go
//
// Generated by this command:
//
//	mockgen -source=usecases/common/mercadolivre.go -destination=usecases/common/mock/mercadolivre_mock.go
//

// Package mock_common is a generated GoMock package.
package mock_common

import (
	reflect "reflect"

	common "github.com/Vractos/kloni/usecases/common"
	gomock "go.uber.org/mock/gomock"
)

// MockmeliReaderStore is a mock of meliReaderStore interface.
type MockmeliReaderStore struct {
	ctrl     *gomock.Controller
	recorder *MockmeliReaderStoreMockRecorder
}

// MockmeliReaderStoreMockRecorder is the mock recorder for MockmeliReaderStore.
type MockmeliReaderStoreMockRecorder struct {
	mock *MockmeliReaderStore
}

// NewMockmeliReaderStore creates a new mock instance.
func NewMockmeliReaderStore(ctrl *gomock.Controller) *MockmeliReaderStore {
	mock := &MockmeliReaderStore{ctrl: ctrl}
	mock.recorder = &MockmeliReaderStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmeliReaderStore) EXPECT() *MockmeliReaderStoreMockRecorder {
	return m.recorder
}

// MockmeliWriterStore is a mock of meliWriterStore interface.
type MockmeliWriterStore struct {
	ctrl     *gomock.Controller
	recorder *MockmeliWriterStoreMockRecorder
}

// MockmeliWriterStoreMockRecorder is the mock recorder for MockmeliWriterStore.
type MockmeliWriterStoreMockRecorder struct {
	mock *MockmeliWriterStore
}

// NewMockmeliWriterStore creates a new mock instance.
func NewMockmeliWriterStore(ctrl *gomock.Controller) *MockmeliWriterStore {
	mock := &MockmeliWriterStore{ctrl: ctrl}
	mock.recorder = &MockmeliWriterStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmeliWriterStore) EXPECT() *MockmeliWriterStoreMockRecorder {
	return m.recorder
}

// RefreshCredentials mocks base method.
func (m *MockmeliWriterStore) RefreshCredentials(refreshToken string) (*common.MeliCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshCredentials", refreshToken)
	ret0, _ := ret[0].(*common.MeliCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshCredentials indicates an expected call of RefreshCredentials.
func (mr *MockmeliWriterStoreMockRecorder) RefreshCredentials(refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshCredentials", reflect.TypeOf((*MockmeliWriterStore)(nil).RefreshCredentials), refreshToken)
}

// RegisterCredential mocks base method.
func (m *MockmeliWriterStore) RegisterCredential(code string) (*common.MeliCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCredential", code)
	ret0, _ := ret[0].(*common.MeliCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterCredential indicates an expected call of RegisterCredential.
func (mr *MockmeliWriterStoreMockRecorder) RegisterCredential(code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCredential", reflect.TypeOf((*MockmeliWriterStore)(nil).RegisterCredential), code)
}

// MockmeliReaderOrder is a mock of meliReaderOrder interface.
type MockmeliReaderOrder struct {
	ctrl     *gomock.Controller
	recorder *MockmeliReaderOrderMockRecorder
}

// MockmeliReaderOrderMockRecorder is the mock recorder for MockmeliReaderOrder.
type MockmeliReaderOrderMockRecorder struct {
	mock *MockmeliReaderOrder
}

// NewMockmeliReaderOrder creates a new mock instance.
func NewMockmeliReaderOrder(ctrl *gomock.Controller) *MockmeliReaderOrder {
	mock := &MockmeliReaderOrder{ctrl: ctrl}
	mock.recorder = &MockmeliReaderOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmeliReaderOrder) EXPECT() *MockmeliReaderOrderMockRecorder {
	return m.recorder
}

// FetchOrder mocks base method.
func (m *MockmeliReaderOrder) FetchOrder(orderId, accessToken string) (*common.MeliOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOrder", orderId, accessToken)
	ret0, _ := ret[0].(*common.MeliOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOrder indicates an expected call of FetchOrder.
func (mr *MockmeliReaderOrderMockRecorder) FetchOrder(orderId, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOrder", reflect.TypeOf((*MockmeliReaderOrder)(nil).FetchOrder), orderId, accessToken)
}

// MockmeliReaderAnnouncement is a mock of meliReaderAnnouncement interface.
type MockmeliReaderAnnouncement struct {
	ctrl     *gomock.Controller
	recorder *MockmeliReaderAnnouncementMockRecorder
}

// MockmeliReaderAnnouncementMockRecorder is the mock recorder for MockmeliReaderAnnouncement.
type MockmeliReaderAnnouncementMockRecorder struct {
	mock *MockmeliReaderAnnouncement
}

// NewMockmeliReaderAnnouncement creates a new mock instance.
func NewMockmeliReaderAnnouncement(ctrl *gomock.Controller) *MockmeliReaderAnnouncement {
	mock := &MockmeliReaderAnnouncement{ctrl: ctrl}
	mock.recorder = &MockmeliReaderAnnouncementMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmeliReaderAnnouncement) EXPECT() *MockmeliReaderAnnouncementMockRecorder {
	return m.recorder
}

// GetAnnouncement mocks base method.
func (m *MockmeliReaderAnnouncement) GetAnnouncement(id string) (*common.MeliAnnouncement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncement", id)
	ret0, _ := ret[0].(*common.MeliAnnouncement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncement indicates an expected call of GetAnnouncement.
func (mr *MockmeliReaderAnnouncementMockRecorder) GetAnnouncement(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncement", reflect.TypeOf((*MockmeliReaderAnnouncement)(nil).GetAnnouncement), id)
}

// GetAnnouncements mocks base method.
func (m *MockmeliReaderAnnouncement) GetAnnouncements(ids []string, accessToken string) (*[]common.MeliAnnouncement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncements", ids, accessToken)
	ret0, _ := ret[0].(*[]common.MeliAnnouncement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncements indicates an expected call of GetAnnouncements.
func (mr *MockmeliReaderAnnouncementMockRecorder) GetAnnouncements(ids, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncements", reflect.TypeOf((*MockmeliReaderAnnouncement)(nil).GetAnnouncements), ids, accessToken)
}

// GetAnnouncementsIDsViaSKU mocks base method.
func (m *MockmeliReaderAnnouncement) GetAnnouncementsIDsViaSKU(sku, userId, accessToken string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncementsIDsViaSKU", sku, userId, accessToken)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncementsIDsViaSKU indicates an expected call of GetAnnouncementsIDsViaSKU.
func (mr *MockmeliReaderAnnouncementMockRecorder) GetAnnouncementsIDsViaSKU(sku, userId, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncementsIDsViaSKU", reflect.TypeOf((*MockmeliReaderAnnouncement)(nil).GetAnnouncementsIDsViaSKU), sku, userId, accessToken)
}

// GetDescription mocks base method.
func (m *MockmeliReaderAnnouncement) GetDescription(id string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDescription", id)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDescription indicates an expected call of GetDescription.
func (mr *MockmeliReaderAnnouncementMockRecorder) GetDescription(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDescription", reflect.TypeOf((*MockmeliReaderAnnouncement)(nil).GetDescription), id)
}

// MockmeliWriterAnnouncement is a mock of meliWriterAnnouncement interface.
type MockmeliWriterAnnouncement struct {
	ctrl     *gomock.Controller
	recorder *MockmeliWriterAnnouncementMockRecorder
}

// MockmeliWriterAnnouncementMockRecorder is the mock recorder for MockmeliWriterAnnouncement.
type MockmeliWriterAnnouncementMockRecorder struct {
	mock *MockmeliWriterAnnouncement
}

// NewMockmeliWriterAnnouncement creates a new mock instance.
func NewMockmeliWriterAnnouncement(ctrl *gomock.Controller) *MockmeliWriterAnnouncement {
	mock := &MockmeliWriterAnnouncement{ctrl: ctrl}
	mock.recorder = &MockmeliWriterAnnouncementMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmeliWriterAnnouncement) EXPECT() *MockmeliWriterAnnouncementMockRecorder {
	return m.recorder
}

// AddDescription mocks base method.
func (m *MockmeliWriterAnnouncement) AddDescription(description, announcementId, accessToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDescription", description, announcementId, accessToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDescription indicates an expected call of AddDescription.
func (mr *MockmeliWriterAnnouncementMockRecorder) AddDescription(description, announcementId, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDescription", reflect.TypeOf((*MockmeliWriterAnnouncement)(nil).AddDescription), description, announcementId, accessToken)
}

// PublishAnnouncement mocks base method.
func (m *MockmeliWriterAnnouncement) PublishAnnouncement(announcementJson []byte, accessToken string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishAnnouncement", announcementJson, accessToken)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishAnnouncement indicates an expected call of PublishAnnouncement.
func (mr *MockmeliWriterAnnouncementMockRecorder) PublishAnnouncement(announcementJson, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishAnnouncement", reflect.TypeOf((*MockmeliWriterAnnouncement)(nil).PublishAnnouncement), announcementJson, accessToken)
}

// UpdateQuantity mocks base method.
func (m *MockmeliWriterAnnouncement) UpdateQuantity(quantity int, announcementId, accessToken string, variationIDs ...int) error {
	m.ctrl.T.Helper()
	varargs := []any{quantity, announcementId, accessToken}
	for _, a := range variationIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateQuantity", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuantity indicates an expected call of UpdateQuantity.
func (mr *MockmeliWriterAnnouncementMockRecorder) UpdateQuantity(quantity, announcementId, accessToken any, variationIDs ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{quantity, announcementId, accessToken}, variationIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuantity", reflect.TypeOf((*MockmeliWriterAnnouncement)(nil).UpdateQuantity), varargs...)
}

// MockMercadoLivre is a mock of MercadoLivre interface.
type MockMercadoLivre struct {
	ctrl     *gomock.Controller
	recorder *MockMercadoLivreMockRecorder
}

// MockMercadoLivreMockRecorder is the mock recorder for MockMercadoLivre.
type MockMercadoLivreMockRecorder struct {
	mock *MockMercadoLivre
}

// NewMockMercadoLivre creates a new mock instance.
func NewMockMercadoLivre(ctrl *gomock.Controller) *MockMercadoLivre {
	mock := &MockMercadoLivre{ctrl: ctrl}
	mock.recorder = &MockMercadoLivreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMercadoLivre) EXPECT() *MockMercadoLivreMockRecorder {
	return m.recorder
}

// AddDescription mocks base method.
func (m *MockMercadoLivre) AddDescription(description, announcementId, accessToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDescription", description, announcementId, accessToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDescription indicates an expected call of AddDescription.
func (mr *MockMercadoLivreMockRecorder) AddDescription(description, announcementId, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDescription", reflect.TypeOf((*MockMercadoLivre)(nil).AddDescription), description, announcementId, accessToken)
}

// FetchOrder mocks base method.
func (m *MockMercadoLivre) FetchOrder(orderId, accessToken string) (*common.MeliOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOrder", orderId, accessToken)
	ret0, _ := ret[0].(*common.MeliOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOrder indicates an expected call of FetchOrder.
func (mr *MockMercadoLivreMockRecorder) FetchOrder(orderId, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOrder", reflect.TypeOf((*MockMercadoLivre)(nil).FetchOrder), orderId, accessToken)
}

// GetAnnouncement mocks base method.
func (m *MockMercadoLivre) GetAnnouncement(id string) (*common.MeliAnnouncement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncement", id)
	ret0, _ := ret[0].(*common.MeliAnnouncement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncement indicates an expected call of GetAnnouncement.
func (mr *MockMercadoLivreMockRecorder) GetAnnouncement(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncement", reflect.TypeOf((*MockMercadoLivre)(nil).GetAnnouncement), id)
}

// GetAnnouncements mocks base method.
func (m *MockMercadoLivre) GetAnnouncements(ids []string, accessToken string) (*[]common.MeliAnnouncement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncements", ids, accessToken)
	ret0, _ := ret[0].(*[]common.MeliAnnouncement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncements indicates an expected call of GetAnnouncements.
func (mr *MockMercadoLivreMockRecorder) GetAnnouncements(ids, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncements", reflect.TypeOf((*MockMercadoLivre)(nil).GetAnnouncements), ids, accessToken)
}

// GetAnnouncementsIDsViaSKU mocks base method.
func (m *MockMercadoLivre) GetAnnouncementsIDsViaSKU(sku, userId, accessToken string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncementsIDsViaSKU", sku, userId, accessToken)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncementsIDsViaSKU indicates an expected call of GetAnnouncementsIDsViaSKU.
func (mr *MockMercadoLivreMockRecorder) GetAnnouncementsIDsViaSKU(sku, userId, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncementsIDsViaSKU", reflect.TypeOf((*MockMercadoLivre)(nil).GetAnnouncementsIDsViaSKU), sku, userId, accessToken)
}

// GetDescription mocks base method.
func (m *MockMercadoLivre) GetDescription(id string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDescription", id)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDescription indicates an expected call of GetDescription.
func (mr *MockMercadoLivreMockRecorder) GetDescription(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDescription", reflect.TypeOf((*MockMercadoLivre)(nil).GetDescription), id)
}

// PublishAnnouncement mocks base method.
func (m *MockMercadoLivre) PublishAnnouncement(announcementJson []byte, accessToken string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishAnnouncement", announcementJson, accessToken)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishAnnouncement indicates an expected call of PublishAnnouncement.
func (mr *MockMercadoLivreMockRecorder) PublishAnnouncement(announcementJson, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishAnnouncement", reflect.TypeOf((*MockMercadoLivre)(nil).PublishAnnouncement), announcementJson, accessToken)
}

// RefreshCredentials mocks base method.
func (m *MockMercadoLivre) RefreshCredentials(refreshToken string) (*common.MeliCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshCredentials", refreshToken)
	ret0, _ := ret[0].(*common.MeliCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshCredentials indicates an expected call of RefreshCredentials.
func (mr *MockMercadoLivreMockRecorder) RefreshCredentials(refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshCredentials", reflect.TypeOf((*MockMercadoLivre)(nil).RefreshCredentials), refreshToken)
}

// RegisterCredential mocks base method.
func (m *MockMercadoLivre) RegisterCredential(code string) (*common.MeliCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCredential", code)
	ret0, _ := ret[0].(*common.MeliCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterCredential indicates an expected call of RegisterCredential.
func (mr *MockMercadoLivreMockRecorder) RegisterCredential(code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCredential", reflect.TypeOf((*MockMercadoLivre)(nil).RegisterCredential), code)
}

// UpdateQuantity mocks base method.
func (m *MockMercadoLivre) UpdateQuantity(quantity int, announcementId, accessToken string, variationIDs ...int) error {
	m.ctrl.T.Helper()
	varargs := []any{quantity, announcementId, accessToken}
	for _, a := range variationIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateQuantity", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuantity indicates an expected call of UpdateQuantity.
func (mr *MockMercadoLivreMockRecorder) UpdateQuantity(quantity, announcementId, accessToken any, variationIDs ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{quantity, announcementId, accessToken}, variationIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuantity", reflect.TypeOf((*MockMercadoLivre)(nil).UpdateQuantity), varargs...)
}
