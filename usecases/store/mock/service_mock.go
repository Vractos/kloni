// Code generated by MockGen. DO NOT EDIT.
// Source: usecases/store/interface.go
//
// Generated by this command:
//
//	mockgen -source=usecases/store/interface.go -destination=usecases/store/mock/service_mock.go
//

// Package mock_store is a generated GoMock package.
package mock_store

import (
	reflect "reflect"

	entity "github.com/Vractos/kloni/entity"
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

// RefreshMeliCredential mocks base method.
func (m *MockUseCase) RefreshMeliCredential(accountId entity.ID, refreshToken string) (*store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshMeliCredential", accountId, refreshToken)
	ret0, _ := ret[0].(*store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshMeliCredential indicates an expected call of RefreshMeliCredential.
func (mr *MockUseCaseMockRecorder) RefreshMeliCredential(accountId, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshMeliCredential", reflect.TypeOf((*MockUseCase)(nil).RefreshMeliCredential), accountId, refreshToken)
}

// RegisterMeliCredentials mocks base method.
func (m *MockUseCase) RegisterMeliCredentials(input store.RegisterMeliCredentialsDtoInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterMeliCredentials", input)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterMeliCredentials indicates an expected call of RegisterMeliCredentials.
func (mr *MockUseCaseMockRecorder) RegisterMeliCredentials(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterMeliCredentials", reflect.TypeOf((*MockUseCase)(nil).RegisterMeliCredentials), input)
}

// RegisterStore mocks base method.
func (m *MockUseCase) RegisterStore(input store.RegisterStoreDtoInput) (entity.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterStore", input)
	ret0, _ := ret[0].(entity.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterStore indicates an expected call of RegisterStore.
func (mr *MockUseCaseMockRecorder) RegisterStore(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterStore", reflect.TypeOf((*MockUseCase)(nil).RegisterStore), input)
}

// RetrieveMeliCredentialsFromMeliUserID mocks base method.
func (m *MockUseCase) RetrieveMeliCredentialsFromMeliUserID(id string) (*[]store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMeliCredentialsFromMeliUserID", id)
	ret0, _ := ret[0].(*[]store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMeliCredentialsFromMeliUserID indicates an expected call of RetrieveMeliCredentialsFromMeliUserID.
func (mr *MockUseCaseMockRecorder) RetrieveMeliCredentialsFromMeliUserID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMeliCredentialsFromMeliUserID", reflect.TypeOf((*MockUseCase)(nil).RetrieveMeliCredentialsFromMeliUserID), id)
}

// RetrieveMeliCredentialsFromStoreID mocks base method.
func (m *MockUseCase) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMeliCredentialsFromStoreID", id)
	ret0, _ := ret[0].(*[]store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMeliCredentialsFromStoreID indicates an expected call of RetrieveMeliCredentialsFromStoreID.
func (mr *MockUseCaseMockRecorder) RetrieveMeliCredentialsFromStoreID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMeliCredentialsFromStoreID", reflect.TypeOf((*MockUseCase)(nil).RetrieveMeliCredentialsFromStoreID), id)
}

// MockRepoReader is a mock of RepoReader interface.
type MockRepoReader struct {
	ctrl     *gomock.Controller
	recorder *MockRepoReaderMockRecorder
}

// MockRepoReaderMockRecorder is the mock recorder for MockRepoReader.
type MockRepoReaderMockRecorder struct {
	mock *MockRepoReader
}

// NewMockRepoReader creates a new mock instance.
func NewMockRepoReader(ctrl *gomock.Controller) *MockRepoReader {
	mock := &MockRepoReader{ctrl: ctrl}
	mock.recorder = &MockRepoReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoReader) EXPECT() *MockRepoReaderMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRepoReader) Get(id string) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepoReaderMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepoReader)(nil).Get), id)
}

// RetrieveMeliCredentialsFromMeliUserID mocks base method.
func (m *MockRepoReader) RetrieveMeliCredentialsFromMeliUserID(accountId string) (*[]store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMeliCredentialsFromMeliUserID", accountId)
	ret0, _ := ret[0].(*[]store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMeliCredentialsFromMeliUserID indicates an expected call of RetrieveMeliCredentialsFromMeliUserID.
func (mr *MockRepoReaderMockRecorder) RetrieveMeliCredentialsFromMeliUserID(accountId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMeliCredentialsFromMeliUserID", reflect.TypeOf((*MockRepoReader)(nil).RetrieveMeliCredentialsFromMeliUserID), accountId)
}

// RetrieveMeliCredentialsFromStoreID mocks base method.
func (m *MockRepoReader) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMeliCredentialsFromStoreID", id)
	ret0, _ := ret[0].(*[]store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMeliCredentialsFromStoreID indicates an expected call of RetrieveMeliCredentialsFromStoreID.
func (mr *MockRepoReaderMockRecorder) RetrieveMeliCredentialsFromStoreID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMeliCredentialsFromStoreID", reflect.TypeOf((*MockRepoReader)(nil).RetrieveMeliCredentialsFromStoreID), id)
}

// MockRepoWriter is a mock of RepoWriter interface.
type MockRepoWriter struct {
	ctrl     *gomock.Controller
	recorder *MockRepoWriterMockRecorder
}

// MockRepoWriterMockRecorder is the mock recorder for MockRepoWriter.
type MockRepoWriterMockRecorder struct {
	mock *MockRepoWriter
}

// NewMockRepoWriter creates a new mock instance.
func NewMockRepoWriter(ctrl *gomock.Controller) *MockRepoWriter {
	mock := &MockRepoWriter{ctrl: ctrl}
	mock.recorder = &MockRepoWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoWriter) EXPECT() *MockRepoWriterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepoWriter) Create(e *entity.Store) (entity.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", e)
	ret0, _ := ret[0].(entity.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepoWriterMockRecorder) Create(e any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepoWriter)(nil).Create), e)
}

// Delete mocks base method.
func (m *MockRepoWriter) Delete(id entity.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepoWriterMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepoWriter)(nil).Delete), id)
}

// RegisterMeliCredential mocks base method.
func (m *MockRepoWriter) RegisterMeliCredential(id, owner_id entity.ID, c *common.MeliCredential, account_name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterMeliCredential", id, owner_id, c, account_name)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterMeliCredential indicates an expected call of RegisterMeliCredential.
func (mr *MockRepoWriterMockRecorder) RegisterMeliCredential(id, owner_id, c, account_name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterMeliCredential", reflect.TypeOf((*MockRepoWriter)(nil).RegisterMeliCredential), id, owner_id, c, account_name)
}

// Update mocks base method.
func (m *MockRepoWriter) Update(e *entity.Store) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", e)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepoWriterMockRecorder) Update(e any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepoWriter)(nil).Update), e)
}

// UpdateMeliCredentials mocks base method.
func (m *MockRepoWriter) UpdateMeliCredentials(accountId entity.ID, c *common.MeliCredential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMeliCredentials", accountId, c)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMeliCredentials indicates an expected call of UpdateMeliCredentials.
func (mr *MockRepoWriterMockRecorder) UpdateMeliCredentials(accountId, c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMeliCredentials", reflect.TypeOf((*MockRepoWriter)(nil).UpdateMeliCredentials), accountId, c)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(e *entity.Store) (entity.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", e)
	ret0, _ := ret[0].(entity.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(e any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), e)
}

// Delete mocks base method.
func (m *MockRepository) Delete(id entity.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockRepository) Get(id string) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), id)
}

// RegisterMeliCredential mocks base method.
func (m *MockRepository) RegisterMeliCredential(id, owner_id entity.ID, c *common.MeliCredential, account_name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterMeliCredential", id, owner_id, c, account_name)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterMeliCredential indicates an expected call of RegisterMeliCredential.
func (mr *MockRepositoryMockRecorder) RegisterMeliCredential(id, owner_id, c, account_name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterMeliCredential", reflect.TypeOf((*MockRepository)(nil).RegisterMeliCredential), id, owner_id, c, account_name)
}

// RetrieveMeliCredentialsFromMeliUserID mocks base method.
func (m *MockRepository) RetrieveMeliCredentialsFromMeliUserID(accountId string) (*[]store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMeliCredentialsFromMeliUserID", accountId)
	ret0, _ := ret[0].(*[]store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMeliCredentialsFromMeliUserID indicates an expected call of RetrieveMeliCredentialsFromMeliUserID.
func (mr *MockRepositoryMockRecorder) RetrieveMeliCredentialsFromMeliUserID(accountId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMeliCredentialsFromMeliUserID", reflect.TypeOf((*MockRepository)(nil).RetrieveMeliCredentialsFromMeliUserID), accountId)
}

// RetrieveMeliCredentialsFromStoreID mocks base method.
func (m *MockRepository) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]store.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveMeliCredentialsFromStoreID", id)
	ret0, _ := ret[0].(*[]store.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveMeliCredentialsFromStoreID indicates an expected call of RetrieveMeliCredentialsFromStoreID.
func (mr *MockRepositoryMockRecorder) RetrieveMeliCredentialsFromStoreID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveMeliCredentialsFromStoreID", reflect.TypeOf((*MockRepository)(nil).RetrieveMeliCredentialsFromStoreID), id)
}

// Update mocks base method.
func (m *MockRepository) Update(e *entity.Store) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", e)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(e any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), e)
}

// UpdateMeliCredentials mocks base method.
func (m *MockRepository) UpdateMeliCredentials(accountId entity.ID, c *common.MeliCredential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMeliCredentials", accountId, c)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMeliCredentials indicates an expected call of UpdateMeliCredentials.
func (mr *MockRepositoryMockRecorder) UpdateMeliCredentials(accountId, c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMeliCredentials", reflect.TypeOf((*MockRepository)(nil).UpdateMeliCredentials), accountId, c)
}
