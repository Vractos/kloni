package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vractos/kloni/adapter/api/handler"
	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/contexttools"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/announcement"
	mock_announcement "github.com/Vractos/kloni/usecases/announcement/mock"
	"github.com/Vractos/kloni/usecases/common"
	common_mock "github.com/Vractos/kloni/usecases/common/mock"
	"github.com/Vractos/kloni/usecases/order"
	mock_order "github.com/Vractos/kloni/usecases/order/mock"
	"github.com/Vractos/kloni/usecases/store"
	mock_store "github.com/Vractos/kloni/usecases/store/mock"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// setupRouter creates a chi router with all handlers wired up.
// Auth middleware is skipped; store ID is injected via a test middleware.
func setupRouter(
	storeUC store.UseCase,
	orderUC order.UseCase,
	annUC announcement.UseCase,
	logger metrics.Logger,
	storeID string,
) *chi.Mux {
	r := chi.NewRouter()

	// Public routes
	handler.MakeStoreHandlers(r, storeUC, logger)
	handler.MakeOrderHandlers(r, orderUC, logger)

	// Protected routes - inject store ID into context instead of JWT middleware
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), contexttools.ContextKeyStoreId, storeID)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
		handler.MakeAnnouncementHandlers(r, annUC, storeUC, logger)
	})

	return r
}

// TestE2E_RegisterStore tests the full flow: POST /store -> register store -> 201 with ID
func TestE2E_RegisterStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreUC := mock_store.NewMockUseCase(ctrl)
	mockOrderUC := mock_order.NewMockUseCase(ctrl)
	mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
	logger := *metrics.NewLogger("info")

	storeID := entity.ID(uuid.New())
	mockStoreUC.EXPECT().RegisterStore(store.RegisterStoreDtoInput{
		Email: "integration@test.com",
		Name:  "Integration Store",
	}).Return(storeID, nil)

	router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, "")

	body, _ := json.Marshal(map[string]string{
		"email": "integration@test.com",
		"name":  "Integration Store",
	})
	req := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["id"] == nil {
		t.Error("expected 'id' in response")
	}
}

// TestE2E_OrderWebhook tests: POST /order/meli-notification -> process webhook -> 200
func TestE2E_OrderWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreUC := mock_store.NewMockUseCase(ctrl)
	mockOrderUC := mock_order.NewMockUseCase(ctrl)
	mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
	logger := *metrics.NewLogger("info")

	mockOrderUC.EXPECT().ProcessWebhook(gomock.Any()).Return(nil)

	router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, "")

	input := order.OrderWebhookDtoInput{
		ID:            "webhook-001",
		Resource:      "/orders/999",
		UserID:        42,
		Topic:         "orders_v2",
		ApplicationID: 1,
		Attempts:      0,
		Sent:          "2024-01-01T00:00:00Z",
		Received:      "2024-01-01T00:00:01Z",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/order/meli-notification", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TestE2E_GetAnnouncements tests: GET /announcement/{sku} -> retrieve announcements -> 200 with JSON
func TestE2E_GetAnnouncements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreUC := mock_store.NewMockUseCase(ctrl)
	mockOrderUC := mock_order.NewMockUseCase(ctrl)
	mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
	logger := *metrics.NewLogger("info")

	testStoreID := entity.ID(uuid.New())
	accountID := entity.ID(uuid.New())
	accountName := "Main Account"

	creds := &[]store.Credentials{{
		ID:          accountID,
		AccountName: &accountName,
		MeliCredential: &common.MeliCredential{
			AccessToken: "integration-token",
			UserID:      "meli-user-1",
		},
	}}

	anns := &[]announcement.Announcements{{
		AccountID:   accountID,
		AccountName: "Main Account",
		Announcements: &[]common.MeliAnnouncement{
			{ID: "MLB-INT-001", Title: "Integration Product", Price: 99.99, Sku: "INT-SKU", Quantity: 5},
			{ID: "MLB-INT-002", Title: "Integration Product 2", Price: 149.99, Sku: "INT-SKU", Quantity: 3},
		},
	}}

	mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(testStoreID).Return(creds, nil)
	mockAnnUC.EXPECT().RetrieveAnnouncementsFromAllAccounts("INT-SKU", creds).Return(anns, nil)

	router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, testStoreID.String())

	req := httptest.NewRequest(http.MethodGet, "/announcement/INT-SKU", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("expected 2 announcements, got %d", len(response))
	}
	if response[0]["id"] != "MLB-INT-001" {
		t.Errorf("expected first announcement ID 'MLB-INT-001', got %v", response[0]["id"])
	}
}

// TestE2E_CloneAnnouncement tests: POST /announcement -> clone -> 201
func TestE2E_CloneAnnouncement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreUC := mock_store.NewMockUseCase(ctrl)
	mockOrderUC := mock_order.NewMockUseCase(ctrl)
	mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
	logger := *metrics.NewLogger("info")

	testStoreID := entity.ID(uuid.New())
	creds := &[]store.Credentials{{
		ID:             entity.ID(uuid.New()),
		MeliCredential: &common.MeliCredential{AccessToken: "token"},
	}}

	mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(testStoreID).Return(creds, nil)
	mockAnnUC.EXPECT().CloneAnnouncement(gomock.Any(), creds).Return(nil)

	router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, testStoreID.String())

	input := announcement.CloneAnnouncementDtoInput{
		RootID:          "MLB-ROOT-001",
		Titles:          []string{"Clone Title 1", "Clone Title 2"},
		RootAccountID:   entity.ID(uuid.New()),
		DestinyAccounts: []entity.ID{entity.ID(uuid.New())},
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/announcement", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

// TestE2E_ImportAnnouncement tests: POST /announcement/import -> import -> 201
func TestE2E_ImportAnnouncement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreUC := mock_store.NewMockUseCase(ctrl)
	mockOrderUC := mock_order.NewMockUseCase(ctrl)
	mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
	logger := *metrics.NewLogger("info")

	testStoreID := entity.ID(uuid.New())
	creds := &[]store.Credentials{{
		ID:             entity.ID(uuid.New()),
		MeliCredential: &common.MeliCredential{AccessToken: "token"},
	}}

	mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(testStoreID).Return(creds, nil)
	mockAnnUC.EXPECT().ImportAnnouncement(gomock.Any(), creds).Return(nil)

	router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, testStoreID.String())

	input := announcement.ImportAnnouncementDtoInput{
		AnnouncementID: "MLB-IMPORT-001",
		AccountOrigin:  entity.ID(uuid.New()),
		AccountDestiny: entity.ID(uuid.New()),
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/announcement/import", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

// TestE2E_ErrorFlows tests error scenarios across the full stack
func TestE2E_ErrorFlows(t *testing.T) {
	t.Run("register store - service failure returns 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreUC := mock_store.NewMockUseCase(ctrl)
		mockOrderUC := mock_order.NewMockUseCase(ctrl)
		mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
		logger := *metrics.NewLogger("info")

		mockStoreUC.EXPECT().RegisterStore(gomock.Any()).Return(entity.ID(uuid.Nil), errors.New("db down"))

		router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, "")

		body, _ := json.Marshal(map[string]string{"email": "test@test.com", "name": "Test"})
		req := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("webhook processing failure returns 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreUC := mock_store.NewMockUseCase(ctrl)
		mockOrderUC := mock_order.NewMockUseCase(ctrl)
		mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
		logger := *metrics.NewLogger("info")

		mockOrderUC.EXPECT().ProcessWebhook(gomock.Any()).Return(errors.New("queue full"))

		router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, "")

		body, _ := json.Marshal(order.OrderWebhookDtoInput{ID: "err-001", Resource: "/orders/1", UserID: 1})
		req := httptest.NewRequest(http.MethodPost, "/order/meli-notification", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("get announcements - credentials failure returns 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreUC := mock_store.NewMockUseCase(ctrl)
		mockOrderUC := mock_order.NewMockUseCase(ctrl)
		mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
		logger := *metrics.NewLogger("info")

		testStoreID := entity.ID(uuid.New())
		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(testStoreID).
			Return(nil, errors.New("creds not found"))

		router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, testStoreID.String())

		req := httptest.NewRequest(http.MethodGet, "/announcement/SKU001", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("get announcements - not found returns 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreUC := mock_store.NewMockUseCase(ctrl)
		mockOrderUC := mock_order.NewMockUseCase(ctrl)
		mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
		logger := *metrics.NewLogger("info")

		testStoreID := entity.ID(uuid.New())
		creds := &[]store.Credentials{{
			ID:             entity.ID(uuid.New()),
			MeliCredential: &common.MeliCredential{AccessToken: "token"},
		}}

		mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(testStoreID).Return(creds, nil)
		mockAnnUC.EXPECT().RetrieveAnnouncementsFromAllAccounts("NOSKU", creds).Return(nil, nil)

		router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, testStoreID.String())

		req := httptest.NewRequest(http.MethodGet, "/announcement/NOSKU", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

// TestE2E_MultiStepFlow tests a realistic multi-step workflow
func TestE2E_MultiStepFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreUC := mock_store.NewMockUseCase(ctrl)
	mockOrderUC := mock_order.NewMockUseCase(ctrl)
	mockAnnUC := mock_announcement.NewMockUseCase(ctrl)
	mockMeli := common_mock.NewMockMercadoLivre(ctrl)
	_ = mockMeli // Available if needed for deeper integration

	logger := *metrics.NewLogger("info")

	// Step 1: Register a store
	storeID := entity.ID(uuid.New())
	mockStoreUC.EXPECT().RegisterStore(store.RegisterStoreDtoInput{
		Email: "flow@test.com",
		Name:  "Flow Store",
	}).Return(storeID, nil)

	router := setupRouter(mockStoreUC, mockOrderUC, mockAnnUC, logger, storeID.String())

	body1, _ := json.Marshal(map[string]string{"email": "flow@test.com", "name": "Flow Store"})
	req1 := httptest.NewRequest(http.MethodPost, "/store", bytes.NewReader(body1))
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("Step 1: expected 201, got %d", w1.Code)
	}

	// Step 2: Receive an order webhook
	mockOrderUC.EXPECT().ProcessWebhook(gomock.Any()).Return(nil)

	webhookBody, _ := json.Marshal(order.OrderWebhookDtoInput{
		ID:       "flow-webhook",
		Resource: "/orders/flow-order-1",
		UserID:   100,
		Topic:    "orders_v2",
	})
	req2 := httptest.NewRequest(http.MethodPost, "/order/meli-notification", bytes.NewReader(webhookBody))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Step 2: expected 200, got %d", w2.Code)
	}

	// Step 3: Clone an announcement
	creds := &[]store.Credentials{{
		ID: entity.ID(uuid.New()),
		MeliCredential: &common.MeliCredential{
			AccessToken: "flow-token",
			UserID:      "meli-100",
		},
	}}

	mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
	mockAnnUC.EXPECT().CloneAnnouncement(gomock.Any(), creds).Return(nil)

	cloneInput := announcement.CloneAnnouncementDtoInput{
		RootID:          "MLB-FLOW-ROOT",
		Titles:          []string{"Flow Clone"},
		RootAccountID:   entity.ID(uuid.New()),
		DestinyAccounts: []entity.ID{entity.ID(uuid.New())},
	}
	body3, _ := json.Marshal(cloneInput)
	req3 := httptest.NewRequest(http.MethodPost, "/announcement", bytes.NewReader(body3))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusCreated {
		t.Fatalf("Step 3: expected 201, got %d, body: %s", w3.Code, w3.Body.String())
	}

	// Step 4: Retrieve announcements
	mockStoreUC.EXPECT().RetrieveMeliCredentialsFromStoreID(storeID).Return(creds, nil)
	mockAnnUC.EXPECT().RetrieveAnnouncementsFromAllAccounts("FLOW-SKU", creds).Return(&[]announcement.Announcements{{
		AccountID:   (*creds)[0].ID,
		AccountName: "Flow Account",
		Announcements: &[]common.MeliAnnouncement{{
			ID: "MLB-FLOW-001", Title: "Flow Product", Price: 50.0, Sku: "FLOW-SKU", Quantity: 10,
		}},
	}}, nil)

	req4 := httptest.NewRequest(http.MethodGet, "/announcement/FLOW-SKU", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusOK {
		t.Fatalf("Step 4: expected 200, got %d, body: %s", w4.Code, w4.Body.String())
	}

	var annResponse []map[string]interface{}
	if err := json.Unmarshal(w4.Body.Bytes(), &annResponse); err != nil {
		t.Fatalf("Step 4: failed to decode: %v", err)
	}
	if len(annResponse) != 1 {
		t.Fatalf("Step 4: expected 1 announcement, got %d", len(annResponse))
	}
}
