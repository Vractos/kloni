package entity

type Announcement struct {
	ID          string
	Title       string
	Price       float32
	Sku         string
	Quantity    int
	Description string
	IsNew       bool
	Images      []string
	Category    string
}

// func NewAnnouncement()
