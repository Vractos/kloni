package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/contexttools"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/announcement"
	common "github.com/Vractos/kloni/usecases/common"
	mock_announcement "github.com/Vractos/kloni/usecases/announcement/mock"
	"github.com/Vractos/kloni/usecases/order"
	mock_order "github.com/Vractos/kloni/usecases/order/mock"
	"github.com/Vractos/kloni/usecases/store"
	mock_store "github.com/Vractos/kloni/usecases/store/mock"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func testLogger() metrics.Logger {
	return *metrics.NewLogger("info")
}

// ---------- Store Handler Tests ----------

func TestRegisterStore(t *testing.T) {
	logger := testLogger()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mock_store.NewMockUseCase(ctrl)
		storeID := entity.ID(uuid.New())

		mockStore.EXPECT().RegisterStore(store.RegisterStoreDtoInput{
			Email: "test@test.com",
			Name:  "Test Store",
		}).Return(storeID, nil)

		body, _ := json.Marshal(map[string]string{"email": "test@test.com", "name": "Test Store"})
		req := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader(body))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		MakeStoreHandlers(r, mockStore, logger)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mock_store.NewMockUseCase(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		MakeStoreHandlers(r, mockStore, logger)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mock_store.NewMockUseCase(ctrl)
		mockStore.EXPECT().RegisterStore(gomock.Any()).Return(entity.ID(uuid.Nil), errors.New("db error"))

		body, _ := json.Marshal(map[string]string{"email": "test@test.com", "name": "Test Store"})
		req := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader(body))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		MakeStoreHandlers(r, mockStore, logger)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestRegisterMeliCredentials(t *testing.T) {
	logger := testLogger()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mock_store.NewMockUseCase(ctrl)
		storeID := entity.ID(uuid.New())
		mockStore.EXPECT().RegisterMeliCredentials(gomock.Any()).Return(nil)

		body, _ := json.Marshal(map[string]string{"code": "auth-code", "account_name": "My Account"})
		req := httptest.NewRequest(http.MethodPost, "/store/meli-credentials", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := registerMeliCredentials(mockStore, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("empty store ID in context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mock_store.NewMockUseCase(ctrl)

		body, _ := json.Marshal(map[string]string{"code": "auth-code"})
		req := httptest.NewRequest(http.MethodPost, "/store/meli-credentials", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, "")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := registerMeliCredentials(mockStore, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mock_store.NewMockUseCase(ctrl)
		storeID := entity.ID(uuid.New())
		mockStore.EXPECT().RegisterMeliCredentials(gomock.Any()).Return(errors.New("cred error"))

		body, _ := json.Marshal(map[string]string{"code": "auth-code"})
		req := httptest.NewRequest(http.MethodPost, "/store/meli-credentials", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := registerMeliCredentials(mockStore, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

// ---------- Order Handler Tests ----------

func TestReceiveMeliOrderNotification(t *testing.T) {
	logger := testLogger()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrder := mock_order.NewMockUseCase(ctrl)
		mockOrder.EXPECT().ProcessWebhook(gomock.Any()).Return(nil)

		input := order.OrderWebhookDtoInput{
			ID:       "test-id",
			Resource: "/orders/123",
			UserID:   1,
			Topic:    "orders_v2",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/order/meli-notification", bytes.NewReader(body))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		MakeOrderHandlers(r, mockOrder, logger)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrder := mock_order.NewMockUseCase(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/order/meli-notification", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		MakeOrderHandlers(r, mockOrder, logger)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrder := mock_order.NewMockUseCase(ctrl)
		mockOrder.EXPECT().ProcessWebhook(gomock.Any()).Return(errors.New("queue error"))

		input := order.OrderWebhookDtoInput{
			ID:       "test-id",
			Resource: "/orders/123",
			UserID:   1,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/order/meli-notification", bytes.NewReader(body))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		MakeOrderHandlers(r, mockOrder, logger)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

// ---------- Announcement Handler Tests ----------

func TestCloneAnnouncementHandler(t *testing.T) {
	logger := testLogger()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().CloneAnnouncement(gomock.Any(), creds).Return(nil)

		input := announcement.CloneAnnouncementDtoInput{RootID: "MLB001"}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/announcement", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := cloneAnnouncement(mockAnn, mockStoreUC, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/announcement", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()

		handler := cloneAnnouncement(mockAnn, mockStoreUC, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("credential retrieval error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(nil, errors.New("no creds"))

		body, _ := json.Marshal(announcement.CloneAnnouncementDtoInput{RootID: "MLB001"})
		req := httptest.NewRequest(http.MethodPost, "/announcement", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := cloneAnnouncement(mockAnn, mockStoreUC, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("clone service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().CloneAnnouncement(gomock.Any(), creds).Return(errors.New("clone failed"))

		body, _ := json.Marshal(announcement.CloneAnnouncementDtoInput{RootID: "MLB001"})
		req := httptest.NewRequest(http.MethodPost, "/announcement", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := cloneAnnouncement(mockAnn, mockStoreUC, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestGetAnnouncementsHandler(t *testing.T) {
	logger := testLogger()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		accountID := entity.ID(uuid.New())
		accountName := "Account1"

		creds := &[]store.Credentials{{
			ID:             accountID,
			AccountName:    &accountName,
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		anns := &[]announcement.Announcements{{
			AccountID:   accountID,
			AccountName: "Account1",
			Announcements: &[]common.MeliAnnouncement{{
				ID:    "MLB001",
				Title: "Test Product",
				Price: 100.0,
				Sku:   "SKU001",
			}},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().RetrieveAnnouncementsFromAllAccounts("SKU001", creds).Return(anns, nil)

		r := chi.NewRouter()
		r.Get("/announcement/{sku}", func(w http.ResponseWriter, r *http.Request) {
			handler := getAnnouncements(mockAnn, mockStoreUC, logger)
			handler.ServeHTTP(w, r)
		})

		req := httptest.NewRequest(http.MethodGet, "/announcement/SKU001", nil)
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().RetrieveAnnouncementsFromAllAccounts("NONE", creds).Return(nil, nil)

		r := chi.NewRouter()
		r.Get("/announcement/{sku}", func(w http.ResponseWriter, r *http.Request) {
			handler := getAnnouncements(mockAnn, mockStoreUC, logger)
			handler.ServeHTTP(w, r)
		})

		req := httptest.NewRequest(http.MethodGet, "/announcement/NONE", nil)
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().RetrieveAnnouncementsFromAllAccounts("SKU001", creds).Return(nil, errors.New("api error"))

		r := chi.NewRouter()
		r.Get("/announcement/{sku}", func(w http.ResponseWriter, r *http.Request) {
			handler := getAnnouncements(mockAnn, mockStoreUC, logger)
			handler.ServeHTTP(w, r)
		})

		req := httptest.NewRequest(http.MethodGet, "/announcement/SKU001", nil)
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestImportAnnouncementHandler(t *testing.T) {
	logger := testLogger()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().ImportAnnouncement(gomock.Any(), creds).Return(nil)

		input := announcement.ImportAnnouncementDtoInput{
			AnnouncementID: "MLB001",
			AccountOrigin:  entity.ID(uuid.New()),
			AccountDestiny: entity.ID(uuid.New()),
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/announcement/import", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := importAnnouncement(mockAnn, mockStoreUC, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAnn := mock_announcement.NewMockUseCase(ctrl)
		mockStoreUC := mock_store.NewMockUseCase(ctrl)

		storeID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
		mockAnn.EXPECT().ImportAnnouncement(gomock.Any(), creds).Return(errors.New("import failed"))

		input := announcement.ImportAnnouncementDtoInput{AnnouncementID: "MLB001"}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/announcement/import", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), contexttools.ContextKeyStoreId, storeID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := importAnnouncement(mockAnn, mockStoreUC, logger)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
