package announcement_tests

import (
	"errors"
	"testing"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/announcement"
	"github.com/Vractos/kloni/usecases/common"
	common_mock "github.com/Vractos/kloni/usecases/common/mock"
	"github.com/Vractos/kloni/usecases/store"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func testLogger() metrics.Logger {
	return *metrics.NewLogger("info")
}

// ---------- TestRetrieveAnnouncements ----------

func TestRetrieveAnnouncements(t *testing.T) {
	creds := store.Credentials{
		ID: entity.ID(uuid.New()),
		MeliCredential: &common.MeliCredential{
			AccessToken: "test-access-token",
			UserID:      "user-123",
		},
	}

	t.Run("success with less than 20 announcements", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		expected := &[]common.MeliAnnouncement{
			{ID: "MLB001", Title: "Product 1", Quantity: 5, Price: 10.0, Sku: "SKU001"},
			{ID: "MLB002", Title: "Product 2", Quantity: 3, Price: 20.0, Sku: "SKU001"},
		}

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-123", "test-access-token").
			Return([]string{"MLB001", "MLB002"}, nil)
		mockMeli.EXPECT().GetAnnouncements([]string{"MLB001", "MLB002"}, "test-access-token").
			Return(expected, nil)

		result, err := service.RetrieveAnnouncements("SKU001", creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || len(*result) != 2 {
			t.Fatalf("expected 2 announcements, got %v", result)
		}
	})

	t.Run("empty announcement IDs returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU-EMPTY", "user-123", "test-access-token").
			Return([]string{}, nil)

		result, err := service.RetrieveAnnouncements("SKU-EMPTY", creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("error from GetAnnouncementsIDsViaSKU", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU-ERR", "user-123", "test-access-token").
			Return(nil, errors.New("api error"))

		_, err := service.RetrieveAnnouncements("SKU-ERR", creds)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error from GetAnnouncements", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-123", "test-access-token").
			Return([]string{"MLB001"}, nil)
		mockMeli.EXPECT().GetAnnouncements(gomock.Any(), "test-access-token").
			Return(nil, errors.New("fetch error"))

		_, err := service.RetrieveAnnouncements("SKU001", creds)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success with 20 or more announcements (chunked)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		ids := make([]string, 21)
		for i := range ids {
			ids[i] = uuid.NewString()
		}

		anns20 := make([]common.MeliAnnouncement, 20)
		for i := range anns20 {
			anns20[i] = common.MeliAnnouncement{ID: ids[i], Sku: "SKU001"}
		}
		anns1 := make([]common.MeliAnnouncement, 1)
		anns1[0] = common.MeliAnnouncement{ID: ids[20], Sku: "SKU001"}

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-123", "test-access-token").
			Return(ids, nil)

		gomock.InOrder(
			mockMeli.EXPECT().GetAnnouncements(ids[:20], "test-access-token").Return(&anns20, nil),
			mockMeli.EXPECT().GetAnnouncements(ids[20:], "test-access-token").Return(&anns1, nil),
		)

		result, err := service.RetrieveAnnouncements("SKU001", creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || len(*result) != 21 {
			t.Fatalf("expected 21 announcements, got %d", len(*result))
		}
	})

	t.Run("error from GetAnnouncements in chunked path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		ids := make([]string, 20)
		for i := range ids {
			ids[i] = uuid.NewString()
		}

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-123", "test-access-token").
			Return(ids, nil)
		mockMeli.EXPECT().GetAnnouncements(gomock.Any(), "test-access-token").
			Return(nil, errors.New("chunk error"))

		_, err := service.RetrieveAnnouncements("SKU001", creds)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- TestRetrieveAnnouncementsFromAllAccounts ----------

func TestRetrieveAnnouncementsFromAllAccounts(t *testing.T) {
	accountID1 := entity.ID(uuid.New())
	accountID2 := entity.ID(uuid.New())
	accountName := "Account1"

	credentials := &[]store.Credentials{
		{
			ID:          accountID1,
			AccountName: &accountName,
			MeliCredential: &common.MeliCredential{
				AccessToken: "token-1",
				UserID:      "user-1",
			},
		},
		{
			ID: accountID2,
			MeliCredential: &common.MeliCredential{
				AccessToken: "token-2",
				UserID:      "user-2",
			},
		},
	}

	t.Run("success with multiple credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		anns1 := &[]common.MeliAnnouncement{{ID: "MLB001", Sku: "SKU001"}}
		anns2 := &[]common.MeliAnnouncement{{ID: "MLB002", Sku: "SKU001"}}

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-1", "token-1").Return([]string{"MLB001"}, nil)
		mockMeli.EXPECT().GetAnnouncements([]string{"MLB001"}, "token-1").Return(anns1, nil)
		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-2", "token-2").Return([]string{"MLB002"}, nil)
		mockMeli.EXPECT().GetAnnouncements([]string{"MLB002"}, "token-2").Return(anns2, nil)

		result, err := service.RetrieveAnnouncementsFromAllAccounts("SKU001", credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || len(*result) != 2 {
			t.Fatalf("expected 2 account results, got %v", result)
		}
		if (*result)[0].AccountName != "Account1" {
			t.Errorf("expected account name 'Account1', got %q", (*result)[0].AccountName)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().GetAnnouncementsIDsViaSKU("SKU001", "user-1", "token-1").
			Return(nil, errors.New("api error"))

		_, err := service.RetrieveAnnouncementsFromAllAccounts("SKU001", credentials)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- TestUpdateQuantity ----------

func TestUpdateQuantity(t *testing.T) {
	creds := store.Credentials{
		ID: entity.ID(uuid.New()),
		MeliCredential: &common.MeliCredential{
			AccessToken: "test-token",
			UserID:      "user-1",
		},
	}

	t.Run("success without variation IDs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().UpdateQuantity(10, "MLB001", "test-token").Return(nil)

		err := service.UpdateQuantity("MLB001", 10, creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("success with variation IDs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().UpdateQuantity(5, "MLB001", "test-token", 100, 200).Return(nil)

		err := service.UpdateQuantity("MLB001", 5, creds, 100, 200)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error from meli UpdateQuantity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		mockMeli.EXPECT().UpdateQuantity(10, "MLB001", "test-token").Return(errors.New("update failed"))

		err := service.UpdateQuantity("MLB001", 10, creds)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// ---------- TestCloneAnnouncement ----------

func TestCloneAnnouncement(t *testing.T) {
	rootAccountID := entity.ID(uuid.New())
	destAccountID := entity.ID(uuid.New())

	credentials := &[]store.Credentials{
		{
			ID: rootAccountID,
			MeliCredential: &common.MeliCredential{
				AccessToken: "root-token",
				UserID:      "root-user",
			},
		},
		{
			ID: destAccountID,
			MeliCredential: &common.MeliCredential{
				AccessToken: "dest-token",
				UserID:      "dest-user",
			},
		},
	}

	rootAnnouncement := &common.MeliAnnouncement{
		ID:            "MLB-ROOT",
		Title:         "Root Product",
		Quantity:      10,
		Price:         100.0,
		Sku:           "SKU-ROOT",
		CategoryID:    "CAT001",
		Condition:     "new",
		ListingTypeID: "gold_pro",
		Pictures:      []string{"http://pic1.jpg"},
		Description:   "Product description",
		Channels:      []string{"marketplace"},
	}

	description := "Product description"

	t.Run("success - clone with one title to one destination", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			Titles:          []string{"Custom Title"},
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID},
		}

		// getAnnouncement: GetAnnouncement + GetDescription
		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-ROOT").Return(&description, nil)

		// 2 announcements (classic + 1 titled) published to 1 destination = 2 publishes
		newID := "MLB-NEW"
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(&newID, nil).Times(2)
		mockMeli.EXPECT().AddDescription(description, newID, "dest-token").Return(nil).Times(2)

		// cloneCompatibilityProducts for each publish
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-ROOT", "root-token").Return(nil, nil).Times(2)
		mockMeli.EXPECT().AddCompatibilityException(newID, "dest-token").Return(nil).Times(2)

		err := service.CloneAnnouncement(input, credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("success - no extra titles, only classic", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			Titles:          []string{},
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID},
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-ROOT").Return(&description, nil)

		newID := "MLB-CLASSIC"
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(&newID, nil)
		mockMeli.EXPECT().AddDescription(description, newID, "dest-token").Return(nil)
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-ROOT", "root-token").Return(nil, nil)
		mockMeli.EXPECT().AddCompatibilityException(newID, "dest-token").Return(nil)

		err := service.CloneAnnouncement(input, credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error from getAnnouncement - GetAnnouncement fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID},
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(nil, errors.New("not found"))

		err := service.CloneAnnouncement(input, credentials)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error from getAnnouncement - GetDescription fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID},
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-ROOT").Return(nil, errors.New("desc error"))

		err := service.CloneAnnouncement(input, credentials)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error from PublishAnnouncement", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			Titles:          []string{},
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID},
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-ROOT").Return(&description, nil)
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(nil, errors.New("publish failed"))

		err := service.CloneAnnouncement(input, credentials)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success with compatibilities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			Titles:          []string{},
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID},
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-ROOT").Return(&description, nil)

		newID := "MLB-NEW"
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(&newID, nil)
		mockMeli.EXPECT().AddDescription(description, newID, "dest-token").Return(nil)

		compatProducts := []common.AnnouncementCompatibilityProduct{
			{ID: "COMPAT-1", DomainID: "domain-1"},
		}
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-ROOT", "root-token").Return(compatProducts, nil)
		mockMeli.EXPECT().AddCompatibilities(newID, "dest-token", &compatProducts).Return(nil)

		err := service.CloneAnnouncement(input, credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("multiple destinations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dest2ID := entity.ID(uuid.New())
		credsWithTwo := &[]store.Credentials{
			{
				ID: rootAccountID,
				MeliCredential: &common.MeliCredential{
					AccessToken: "root-token",
					UserID:      "root-user",
				},
			},
			{
				ID: destAccountID,
				MeliCredential: &common.MeliCredential{
					AccessToken: "dest-token",
					UserID:      "dest-user",
				},
			},
			{
				ID: dest2ID,
				MeliCredential: &common.MeliCredential{
					AccessToken: "dest2-token",
					UserID:      "dest2-user",
				},
			},
		}

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.CloneAnnouncementDtoInput{
			RootID:          "MLB-ROOT",
			Titles:          []string{},
			RootAccountID:   rootAccountID,
			DestinyAccounts: []entity.ID{destAccountID, dest2ID},
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-ROOT", "root-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-ROOT").Return(&description, nil)

		newID := "MLB-NEW"
		// 1 classic announcement published to 2 destinations = 2 publishes
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), gomock.Any()).Return(&newID, nil).Times(2)
		mockMeli.EXPECT().AddDescription(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-ROOT", "root-token").Return(nil, nil).Times(2)
		mockMeli.EXPECT().AddCompatibilityException(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		err := service.CloneAnnouncement(input, credsWithTwo)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// ---------- TestImportAnnouncement ----------

func TestImportAnnouncement(t *testing.T) {
	originAccountID := entity.ID(uuid.New())
	destAccountID := entity.ID(uuid.New())

	credentials := &[]store.Credentials{
		{
			ID: originAccountID,
			MeliCredential: &common.MeliCredential{
				AccessToken: "origin-token",
				UserID:      "origin-user",
			},
		},
		{
			ID: destAccountID,
			MeliCredential: &common.MeliCredential{
				AccessToken: "dest-token",
				UserID:      "dest-user",
			},
		},
	}

	rootAnnouncement := &common.MeliAnnouncement{
		ID:            "MLB-IMPORT",
		Title:         "Import Product",
		Quantity:      5,
		Price:         50.0,
		Sku:           "SKU-IMP",
		CategoryID:    "CAT002",
		Condition:     "new",
		ListingTypeID: "gold_pro",
		Pictures:      []string{"http://pic1.jpg"},
		Description:   "Import description",
		Channels:      []string{"marketplace"},
	}

	description := "Import description"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.ImportAnnouncementDtoInput{
			AnnouncementID: "MLB-IMPORT",
			AccountOrigin:  originAccountID,
			AccountDestiny: destAccountID,
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-IMPORT", "origin-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-IMPORT").Return(&description, nil)

		newID := "MLB-IMPORTED"
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(&newID, nil)
		mockMeli.EXPECT().AddDescription(description, newID, "dest-token").Return(nil)
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-IMPORT", "origin-token").Return(nil, nil)
		mockMeli.EXPECT().AddCompatibilityException(newID, "dest-token").Return(nil)

		err := service.ImportAnnouncement(input, credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error from getAnnouncement", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.ImportAnnouncementDtoInput{
			AnnouncementID: "MLB-IMPORT",
			AccountOrigin:  originAccountID,
			AccountDestiny: destAccountID,
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-IMPORT", "origin-token").Return(nil, errors.New("not found"))

		err := service.ImportAnnouncement(input, credentials)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error from PublishAnnouncement", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.ImportAnnouncementDtoInput{
			AnnouncementID: "MLB-IMPORT",
			AccountOrigin:  originAccountID,
			AccountDestiny: destAccountID,
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-IMPORT", "origin-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-IMPORT").Return(&description, nil)
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(nil, errors.New("publish failed"))

		err := service.ImportAnnouncement(input, credentials)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success with compatibilities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.ImportAnnouncementDtoInput{
			AnnouncementID: "MLB-IMPORT",
			AccountOrigin:  originAccountID,
			AccountDestiny: destAccountID,
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-IMPORT", "origin-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-IMPORT").Return(&description, nil)

		newID := "MLB-IMPORTED"
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(&newID, nil)
		mockMeli.EXPECT().AddDescription(description, newID, "dest-token").Return(nil)

		compatProducts := []common.AnnouncementCompatibilityProduct{
			{ID: "COMPAT-1", DomainID: "domain-1"},
		}
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-IMPORT", "origin-token").Return(compatProducts, nil)
		mockMeli.EXPECT().AddCompatibilities(newID, "dest-token", &compatProducts).Return(nil)

		err := service.ImportAnnouncement(input, credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("compatibility error is logged but doesn't fail import", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMeli := common_mock.NewMockMercadoLivre(ctrl)
		service := announcement.NewAnnouncementService(mockMeli, nil, testLogger())

		input := announcement.ImportAnnouncementDtoInput{
			AnnouncementID: "MLB-IMPORT",
			AccountOrigin:  originAccountID,
			AccountDestiny: destAccountID,
		}

		mockMeli.EXPECT().GetAnnouncement("MLB-IMPORT", "origin-token").Return(rootAnnouncement, nil)
		mockMeli.EXPECT().GetDescription("MLB-IMPORT").Return(&description, nil)

		newID := "MLB-IMPORTED"
		mockMeli.EXPECT().PublishAnnouncement(gomock.Any(), "dest-token").Return(&newID, nil)
		mockMeli.EXPECT().AddDescription(description, newID, "dest-token").Return(nil)
		mockMeli.EXPECT().GetAnnouncementCompatibilities("MLB-IMPORT", "origin-token").
			Return(nil, errors.New("compat error"))

		err := service.ImportAnnouncement(input, credentials)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// ---------- TestAnnouncementError ----------

func TestAnnouncementError(t *testing.T) {
	t.Run("error with announcement ID", func(t *testing.T) {
		err := &announcement.AnnouncementError{
			Message:        "test error",
			AnnouncementID: "MLB001",
		}
		expected := "Message: test error - Announcement ID: MLB001"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
		}
	})

	t.Run("error with SKU", func(t *testing.T) {
		err := &announcement.AnnouncementError{
			Message: "test error",
			Sku:     "SKU001",
		}
		expected := "Message: test error - Announcement SKU: SKU001"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
		}
	})

	t.Run("error with only message", func(t *testing.T) {
		err := &announcement.AnnouncementError{
			Message: "test error",
		}
		expected := "Message: test error"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
		}
	})
}
