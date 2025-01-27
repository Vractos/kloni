package common

import (
	"image"
	"time"
)

/*
###################################
###################################
---------------Models-------------
###################################
###################################
*/

type MeliCredential struct {
	AccessToken  string
	ExpiresIn    int
	UserID       string
	RefreshToken string
	UpdatedAt    time.Time
}

type OrderStatus string

const (
	Unknown           OrderStatus = ""
	Confirmed         OrderStatus = "confirmed"
	PaymentRequired   OrderStatus = "payment_required"
	PaymentInProcess  OrderStatus = "payment_in_process"
	PartiallyPaid     OrderStatus = "partially_paid"
	Paid              OrderStatus = "paid"
	PartiallyRefunded OrderStatus = "partially_refunded"
	PendingCancel     OrderStatus = "pending_cancel"
	Cancelled         OrderStatus = "cancelled"
	Invalid           OrderStatus = "invalid"
)

type OrderItem struct {
	ID          string
	Title       string
	Sku         string
	Quantity    int
	VariationID int
}

type MeliOrder struct {
	ID string
	// Can be converted to time
	DateCreated string
	Status      OrderStatus
	Items       []OrderItem
}

type MeliAnnouncement struct {
	ID            string
	Title         string
	Quantity      int
	Price         float64
	ThumbnailURL  string
	Sku           string
	Link          string
	CategoryID    string
	Status        string
	Condition     string
	ListingTypeID string
	Pictures      []string
	Description   string
	Channels      []string
	Variations    []struct {
		ID                int
		AvailableQuantity int
	}
	SaleTerms []struct {
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
	}
	Attributes []struct {
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
	}
}

type AnnouncementCompatibilityProduct struct {
	ID                 string
	DomainID           string
	CatalogProductID   string
	CatalogProductName string
	Source             string
	Universal          bool
}

/*
###################################
###################################
-------------Interfaces------------
###################################
###################################
*/

// Mercado Livre reader interface
type meliReaderStore interface {
	// Implement
}

// Mercado Livre writer interface
type meliWriterStore interface {
	RegisterCredential(code string) (*MeliCredential, error)
	RefreshCredentials(refreshToken string) (*MeliCredential, error)
}

type meliReaderOrder interface {
	FetchOrder(orderId string, accessToken string) (*MeliOrder, error)
}

type meliReaderAnnouncement interface {
	GetAnnouncementsIDsViaSKU(sku string, userId string, accessToken string) ([]string, error)
	// Max 10 IDs
	GetAnnouncements(ids []string, accessToken string) (*[]MeliAnnouncement, error)
	GetAnnouncement(id string, accessToken string) (*MeliAnnouncement, error)
	GetDescription(id string) (*string, error)
	GetProductsPictures(picturesURL []string) (pics []image.Image, err error)
	GetAnnouncementCompatibilities(id string, accessToken string) ([]AnnouncementCompatibilityProduct, error)
}

type meliWriterAnnouncement interface {
	UpdateQuantity(quantity int, announcementId, accessToken string, variationIDs ...int) error
	PublishAnnouncement(announcementJson []byte, accessToken string) (ID *string, err error)
	AddDescription(description, announcementId, accessToken string) error
	ValidateAndExchangeImages(images []*image.Image, accessToken string) (urlF []string, err error)
	AddCompatibilities(announcementId, accessToken string, compatibilities *[]AnnouncementCompatibilityProduct) error
	AddCompatibilityException(announcementId, accessToken string) error
}

type MercadoLivre interface {
	meliReaderStore
	meliReaderOrder
	meliWriterStore
	meliReaderAnnouncement
	meliWriterAnnouncement
}
