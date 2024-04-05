package entity

import (
	"math"
	"strings"

	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/utils"
)

type ListingType string

const (
	Premium    OrderStatus = "gold_pro"
	Diamante   OrderStatus = "gold_premium"
	Classic    OrderStatus = "gold_special"
	Ouro       OrderStatus = "gold"
	Prata      OrderStatus = "silver"
	PaiBronzed OrderStatus = "bronze"
	Free       OrderStatus = "free"
)

type attributes struct {
	ID        string      `json:"id,omitempty"`
	ValueID   interface{} `json:"value_id,omitempty"`
	ValueName interface{} `json:"value_name,omitempty"`
}
type pictures struct {
	Source string `json:"source,omitempty"`
}

type saleTerms struct {
	ID        string      `json:"id,omitempty"`
	ValueName interface{} `json:"value_name,omitempty"`
}

type Variation struct {
	ID                int `json:"id,omitempty"`
	AvailableQuantity int `json:"available_quantity,omitempty"`
}

type Announcement struct {
	Title             string       `json:"title,omitempty"`
	AvailableQuantity int          `json:"available_quantity,omitempty"`
	Price             float64      `json:"price,omitempty"`
	CurrencyID        string       `json:"currency_id,omitempty"`
	BuyingMode        string       `json:"buying_mode,omitempty"`
	CategoryID        string       `json:"category_id,omitempty"`
	Condition         string       `json:"condition,omitempty"`
	ListingTypeID     ListingType  `json:"listing_type_id,omitempty"`
	Channels          []string     `json:"channels,omitempty"`
	Attributes        []attributes `json:"attributes,omitempty"`
	Pictures          []pictures   `json:"pictures,omitempty"`
	SaleTerms         []saleTerms  `json:"sale_terms,omitempty"`
	Variations        []Variation  `json:"variations,omitempty"`
}

// NewAnnouncement creates a new Announcement object based on the provided rootAnn MeliAnnouncement.
// It converts the attributes, pictures, and sale terms from the rootAnn into the corresponding fields of the Announcement object.
// Returns the created Announcement object and an error if any.
func NewAnnouncement(rootAnn *common.MeliAnnouncement) (*Announcement, error) {
	att := make([]attributes, len(rootAnn.Attributes))
	for i, at := range rootAnn.Attributes {
		att[i] = attributes{
			ID:        at.ID,
			ValueName: at.ValueName,
		}
	}

	pics := make([]pictures, len(rootAnn.Pictures))
	for i, p := range rootAnn.Pictures {
		pics[i] = pictures{
			Source: p,
		}
	}

	sTerms := make([]saleTerms, len(rootAnn.SaleTerms))
	for i, st := range rootAnn.SaleTerms {
		sTerms[i] = saleTerms{
			ID:        st.ID,
			ValueName: st.ValueName,
		}
	}

	for i, e := range att {
		if e.ID == "GTIN" {
			break
		}
		if i == len(att)-1 {
			att = append(att, attributes{
				ID:        "GTIN",
				ValueName: strings.ToUpper(rootAnn.Sku),
			})
		}
	}

	return &Announcement{
		Title:             rootAnn.Title,
		AvailableQuantity: rootAnn.Quantity,
		Price:             rootAnn.Price,
		CurrencyID:        "BRL",
		BuyingMode:        "buy_it_now",
		CategoryID:        rootAnn.CategoryID,
		Condition:         rootAnn.Condition,
		ListingTypeID:     ListingType(rootAnn.ListingTypeID),
		Channels:          rootAnn.Channels,
		Attributes:        att,
		Pictures:          pics,
		SaleTerms:         sTerms,
	}, nil
}

func (a *Announcement) ChangeTitle(title string) {
	a.Title = title
}

func (a *Announcement) GenerateClassic() {
	a.ListingTypeID = ListingType(Classic)
	a.Price = math.Round(utils.Percent(95, a.Price)*100) / 100
}
