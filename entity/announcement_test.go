package entity

import (
	"testing"

	"github.com/Vractos/kloni/usecases/common"
	"github.com/google/go-cmp/cmp"
)

func TestNewAnnouncement(t *testing.T) {
	type args struct {
		rootAnn *common.MeliAnnouncement
	}
	tests := []struct {
		name     string
		args     args
		newTitle string
		want     *Announcement
		wantErr  bool
	}{
		{
			name:     "TestNewAnnouncement",
			newTitle: "Radiador De Água Ford Ka 1.6 2000",
			args: args{
				rootAnn: &common.MeliAnnouncement{
					ID:            "MLB123456789",
					Title:         "Radiador De Água Ford Ka 1.6 1997 A 2013",
					Quantity:      1,
					Price:         100.0,
					ThumbnailURL:  "https://http2.mlstatic.com/D_NQ_NP_2X_905878-MLB31068648685_062019-F.webp",
					Sku:           "SKU-123456789",
					Link:          "https://www.mercadolivre.com.br",
					CategoryID:    "MLB5672",
					Condition:     "new",
					ListingTypeID: "gold_pro",
					Pictures: []string{
						"https://http2.mlstatic.com/D_NQ_NP_2X_905878-MLB31068648685_062019-F.webp",
					},
					Description: "TestNewAnnouncement",
					Channels:    []string{"marketplace"},
					SaleTerms: []struct {
						ID          string
						Name        string
						ValueID     interface{}
						ValueName   string
						ValueStruct struct {
							Number int
							Unit   string
						}
						Values []struct {
							ID     interface{}
							Name   string
							Struct struct {
								Number int
								Unit   string
							}
						}
						ValueType string
					}{
						{
							ID:        "WARRANTY_TYPE",
							ValueName: "Garantia do vendedor",
						},
						{
							ID:        "WARRANTY_TIME",
							ValueName: "90 dias",
						},
						{
							ID:        "WARRANTY_TIME_TYPE",
							ValueName: "dias",
						},
					},
					Attributes: []struct {
						ID          string
						Name        string
						ValueID     string
						ValueName   string
						ValueStruct interface{}
						Values      []struct {
							ID     string
							Name   string
							Struct interface{}
						}
						AttributeGroupID   string
						AttributeGroupName string
						ValueType          string
					}{
						{
							ID:        "BRAND",
							ValueID:   "215",
							ValueName: "Apple",
						},
						{
							ID:        "MODEL",
							ValueID:   "215",
							ValueName: "Apple",
						},
					},
				},
			},
			want: &Announcement{
				Title:             "Radiador De Água Ford Ka 1.6 1997 A 2013",
				AvailableQuantity: 1,
				Price:             100.0,
				CurrencyID:        "BRL",
				BuyingMode:        "buy_it_now",
				CategoryID:        "MLB5672",
				Condition:         "new",
				ListingTypeID:     "gold_pro",
				Channels:          []string{"marketplace"},
				Attributes: []attributes{
					{
						ID:        "BRAND",
						ValueID:   "215",
						ValueName: "Apple",
					},
					{
						ID:        "MODEL",
						ValueID:   "215",
						ValueName: "Apple",
					},
				},
				Pictures: []pictures{
					{
						Source: "https://http2.mlstatic.com/D_NQ_NP_2X_905878-MLB31068648685_062019-F.webp",
					},
				},
				SaleTerms: []saleTerms{
					{
						ID:        "WARRANTY_TYPE",
						ValueName: "Garantia do vendedor",
					},
					{
						ID:        "WARRANTY_TIME",
						ValueName: "90 dias",
					},
					{
						ID:        "WARRANTY_TIME_TYPE",
						ValueName: "dias",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := NewAnnouncement(tt.args.rootAnn)
		t.Run(tt.name, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnnouncement() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("NewAnnouncement() diff:\n%v", cmp.Diff(got, tt.want))
			}
		})
		t.Run("generating classic announcement", func(t *testing.T) {
			tt.want.ListingTypeID = "gold_special"
			tt.want.Price = 95.0
			got.GenerateClassic()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GenerateClassic() diff:\n%v", cmp.Diff(got, tt.want))
			}
		})
		t.Run("changing title", func(t *testing.T) {
			tt.want.Title = tt.newTitle
			got.ChangeTitle(tt.newTitle)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ChangeTitle() diff:\n%v", cmp.Diff(got, tt.want))
			}
		})
	}
}
