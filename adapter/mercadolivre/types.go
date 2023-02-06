package mercadolivre

import "time"

type Credentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type Order struct {
	ID                      uint64      `json:"id,omitempty" validate:"required"`
	DateCreated             string      `json:"date_created,omitempty" validate:"required"`
	DateClosed              string      `json:"date_closed,omitempty"`
	LastUpdated             string      `json:"last_updated,omitempty"`
	ManufacturingEndingDate interface{} `json:"manufacturing_ending_date,omitempty"`
	Comment                 interface{} `json:"comment,omitempty"`
	PackID                  int64       `json:"pack_id,omitempty"`
	PickupID                interface{} `json:"pickup_id,omitempty"`
	OrderRequest            struct {
		Return interface{} `json:"return,omitempty"`
		Change interface{} `json:"change,omitempty"`
	} `json:"order_request,omitempty"`
	Fulfilled   interface{}   `json:"fulfilled,omitempty"`
	Mediations  []interface{} `json:"mediations,omitempty"`
	TotalAmount float64       `json:"total_amount,omitempty"`
	PaidAmount  float64       `json:"paid_amount,omitempty"`
	Coupon      struct {
		ID     interface{} `json:"id,omitempty"`
		Amount float64     `json:"amount,omitempty"`
	} `json:"coupon,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	OrderItems     []struct {
		Item struct {
			ID                  string      `json:"id,omitempty" validate:"required"`
			Title               string      `json:"title,omitempty" validate:"required"`
			CategoryID          string      `json:"category_id,omitempty"`
			VariationID         int64       `json:"variation_id,omitempty"`
			SellerCustomField   interface{} `json:"seller_custom_field,omitempty"`
			VariationAttributes []struct {
				ID        string `json:"id,omitempty"`
				Name      string `json:"name,omitempty"`
				ValueID   string `json:"value_id,omitempty"`
				ValueName string `json:"value_name,omitempty"`
			} `json:"variation_attributes,omitempty"`
			Warranty    string      `json:"warranty,omitempty"`
			Condition   string      `json:"condition,omitempty"`
			SellerSku   string      `json:"seller_sku,omitempty" validate:"required"`
			GlobalPrice interface{} `json:"global_price,omitempty"`
			NetWeight   interface{} `json:"net_weight,omitempty"`
		} `json:"item,omitempty"`
		Quantity          int `json:"quantity,omitempty" validate:"required"`
		RequestedQuantity struct {
			Value   int    `json:"value,omitempty"`
			Measure string `json:"measure,omitempty"`
		} `json:"requested_quantity,omitempty"`
		PickedQuantity    interface{} `json:"picked_quantity,omitempty"`
		UnitPrice         float64     `json:"unit_price,omitempty"`
		FullUnitPrice     float64     `json:"full_unit_price,omitempty"`
		CurrencyID        string      `json:"currency_id,omitempty"`
		ManufacturingDays interface{} `json:"manufacturing_days,omitempty"`
		SaleFee           float64     `json:"sale_fee,omitempty"`
		ListingTypeID     string      `json:"listing_type_id,omitempty"`
	} `json:"order_items,omitempty"`
	CurrencyID string `json:"currency_id,omitempty"`
	Payments   []struct {
		ID        int64 `json:"id,omitempty"`
		OrderID   int64 `json:"order_id,omitempty"`
		PayerID   int   `json:"payer_id,omitempty"`
		Collector struct {
			ID int `json:"id,omitempty"`
		} `json:"collector,omitempty"`
		CardID               interface{} `json:"card_id,omitempty"`
		SiteID               string      `json:"site_id,omitempty"`
		Reason               string      `json:"reason,omitempty"`
		PaymentMethodID      string      `json:"payment_method_id,omitempty"`
		CurrencyID           string      `json:"currency_id,omitempty"`
		Installments         int         `json:"installments,omitempty"`
		IssuerID             interface{} `json:"issuer_id,omitempty"`
		AtmTransferReference struct {
			CompanyID     interface{} `json:"company_id,omitempty"`
			TransactionID interface{} `json:"transaction_id,omitempty"`
		} `json:"atm_transfer_reference,omitempty"`
		CouponID                  interface{} `json:"coupon_id,omitempty"`
		ActivationURI             interface{} `json:"activation_uri,omitempty"`
		OperationType             string      `json:"operation_type,omitempty"`
		PaymentType               string      `json:"payment_type,omitempty"`
		AvailableActions          []string    `json:"available_actions,omitempty"`
		Status                    string      `json:"status,omitempty"`
		StatusCode                interface{} `json:"status_code,omitempty"`
		StatusDetail              string      `json:"status_detail,omitempty"`
		TransactionAmount         float64     `json:"transaction_amount,omitempty"`
		TransactionAmountRefunded float64     `json:"transaction_amount_refunded,omitempty"`
		TaxesAmount               float64     `json:"taxes_amount,omitempty"`
		ShippingCost              float64     `json:"shipping_cost,omitempty"`
		CouponAmount              float64     `json:"coupon_amount,omitempty"`
		OverpaidAmount            float64     `json:"overpaid_amount,omitempty"`
		TotalPaidAmount           float64     `json:"total_paid_amount,omitempty"`
		InstallmentAmount         interface{} `json:"installment_amount,omitempty"`
		DeferredPeriod            interface{} `json:"deferred_period,omitempty"`
		DateApproved              string      `json:"date_approved,omitempty"`
		AuthorizationCode         interface{} `json:"authorization_code,omitempty"`
		TransactionOrderID        interface{} `json:"transaction_order_id,omitempty"`
		DateCreated               string      `json:"date_created,omitempty"`
		DateLastModified          string      `json:"date_last_modified,omitempty"`
	} `json:"payments,omitempty"`
	Shipping struct {
		ID int64 `json:"id,omitempty"`
	} `json:"shipping,omitempty"`
	Status       string      `json:"status,omitempty" validate:"required"`
	StatusDetail interface{} `json:"status_detail,omitempty"`
	Tags         []string    `json:"tags,omitempty"`
	Feedback     struct {
		Buyer  interface{} `json:"buyer,omitempty"`
		Seller interface{} `json:"seller,omitempty"`
	} `json:"feedback,omitempty"`
	Context struct {
		Channel string        `json:"channel,omitempty"`
		Site    string        `json:"site,omitempty"`
		Flows   []interface{} `json:"flows,omitempty"`
	} `json:"context,omitempty"`
	Buyer struct {
		ID        int    `json:"id,omitempty"`
		Nickname  string `json:"nickname,omitempty"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	} `json:"buyer,omitempty"`
	Seller struct {
		ID int `json:"id,omitempty"`
	} `json:"seller,omitempty"`
	Taxes struct {
		Amount     float64     `json:"amount,omitempty"`
		CurrencyID interface{} `json:"currency_id,omitempty"`
		ID         interface{} `json:"id,omitempty"`
	} `json:"taxes,omitempty"`
}

type QueryAnnouncementViaSku struct {
	SellerID string      `json:"seller_id,omitempty"`
	Query    interface{} `json:"query,omitempty"`
	Paging   interface{} `json:"paging,omitempty"`
	Results  []string    `json:"results,omitempty"`
	Orders   []struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"orders,omitempty"`
	AvailableOrders interface{} `json:"available_orders,omitempty"`
}

type AnnouncementsMultiGet []struct {
	Code int `json:"code,omitempty"`
	Body struct {
		ID                string      `json:"id,omitempty"`
		Message           string      `json:"message,omitempty"`
		Error             string      `json:"error,omitempty"`
		SiteID            string      `json:"site_id,omitempty"`
		Title             string      `json:"title,omitempty"`
		Subtitle          interface{} `json:"subtitle,omitempty"`
		SellerID          int         `json:"seller_id,omitempty"`
		CategoryID        string      `json:"category_id,omitempty"`
		OfficialStoreID   interface{} `json:"official_store_id,omitempty"`
		Price             float64     `json:"price,omitempty"`
		BasePrice         float64     `json:"base_price,omitempty"`
		OriginalPrice     float64     `json:"original_price,omitempty"`
		InventoryID       interface{} `json:"inventory_id,omitempty"`
		CurrencyID        string      `json:"currency_id,omitempty"`
		InitialQuantity   int         `json:"initial_quantity,omitempty"`
		AvailableQuantity int         `json:"available_quantity,omitempty"`
		SoldQuantity      int         `json:"sold_quantity,omitempty"`
		SaleTerms         []struct {
			ID          string      `json:"id,omitempty"`
			Name        string      `json:"name,omitempty"`
			ValueID     interface{} `json:"value_id,omitempty"`
			ValueName   string      `json:"value_name,omitempty"`
			ValueStruct struct {
				Number int    `json:"number,omitempty"`
				Unit   string `json:"unit,omitempty"`
			} `json:"value_struct,omitempty"`
			Values []struct {
				ID     interface{} `json:"id,omitempty"`
				Name   string      `json:"name,omitempty"`
				Struct struct {
					Number int    `json:"number,omitempty"`
					Unit   string `json:"unit,omitempty"`
				} `json:"struct,omitempty"`
			} `json:"values,omitempty"`
			ValueType string `json:"value_type,omitempty"`
		} `json:"sale_terms,omitempty"`
		BuyingMode      string    `json:"buying_mode,omitempty"`
		ListingTypeID   string    `json:"listing_type_id,omitempty"`
		StartTime       time.Time `json:"start_time,omitempty"`
		StopTime        time.Time `json:"stop_time,omitempty"`
		EndTime         time.Time `json:"end_time,omitempty"`
		ExpirationTime  time.Time `json:"expiration_time,omitempty"`
		Condition       string    `json:"condition,omitempty"`
		Permalink       string    `json:"permalink,omitempty"`
		ThumbnailID     string    `json:"thumbnail_id,omitempty"`
		Thumbnail       string    `json:"thumbnail,omitempty"`
		SecureThumbnail string    `json:"secure_thumbnail,omitempty"`
		Pictures        []struct {
			ID        string `json:"id,omitempty"`
			URL       string `json:"url,omitempty"`
			SecureURL string `json:"secure_url,omitempty"`
			Size      string `json:"size,omitempty"`
			MaxSize   string `json:"max_size,omitempty"`
			Quality   string `json:"quality,omitempty"`
		} `json:"pictures,omitempty"`
		VideoID                      interface{}   `json:"video_id,omitempty"`
		Descriptions                 []interface{} `json:"descriptions,omitempty"`
		AcceptsMercadopago           bool          `json:"accepts_mercadopago,omitempty"`
		NonMercadoPagoPaymentMethods []interface{} `json:"non_mercado_pago_payment_methods,omitempty"`
		Shipping                     struct {
			Mode         string        `json:"mode,omitempty"`
			Methods      []interface{} `json:"methods,omitempty"`
			Tags         []string      `json:"tags,omitempty"`
			Dimensions   interface{}   `json:"dimensions,omitempty"`
			LocalPickUp  bool          `json:"local_pick_up,omitempty"`
			FreeShipping bool          `json:"free_shipping,omitempty"`
			LogisticType string        `json:"logistic_type,omitempty"`
			StorePickUp  bool          `json:"store_pick_up,omitempty"`
		} `json:"shipping,omitempty"`
		InternationalDeliveryMode string `json:"international_delivery_mode,omitempty"`
		SellerAddress             struct {
			Comment     string `json:"comment,omitempty"`
			AddressLine string `json:"address_line,omitempty"`
			ZipCode     string `json:"zip_code,omitempty"`
			City        struct {
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			} `json:"city,omitempty"`
			State struct {
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			} `json:"state,omitempty"`
			Country struct {
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			} `json:"country,omitempty"`
			SearchLocation struct {
				City struct {
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				} `json:"city,omitempty"`
				State struct {
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				} `json:"state,omitempty"`
			} `json:"search_location,omitempty"`
			Latitude  float64 `json:"latitude,omitempty"`
			Longitude float64 `json:"longitude,omitempty"`
			ID        int     `json:"id,omitempty"`
		} `json:"seller_address,omitempty"`
		SellerContact interface{} `json:"seller_contact,omitempty"`
		Location      struct {
		} `json:"location,omitempty"`
		Geolocation struct {
			Latitude  float64 `json:"latitude,omitempty"`
			Longitude float64 `json:"longitude,omitempty"`
		} `json:"geolocation,omitempty"`
		CoverageAreas []interface{} `json:"coverage_areas,omitempty"`
		Attributes    []struct {
			ID          string      `json:"id,omitempty"`
			Name        string      `json:"name,omitempty"`
			ValueID     string      `json:"value_id,omitempty"`
			ValueName   string      `json:"value_name,omitempty"`
			ValueStruct interface{} `json:"value_struct,omitempty"`
			Values      []struct {
				ID     string      `json:"id,omitempty"`
				Name   string      `json:"name,omitempty"`
				Struct interface{} `json:"struct,omitempty"`
			} `json:"values,omitempty"`
			AttributeGroupID   string `json:"attribute_group_id,omitempty"`
			AttributeGroupName string `json:"attribute_group_name,omitempty"`
			ValueType          string `json:"value_type,omitempty"`
		} `json:"attributes,omitempty"`
		Warnings            []interface{} `json:"warnings,omitempty"`
		ListingSource       string        `json:"listing_source,omitempty"`
		Variations          []interface{} `json:"variations,omitempty"`
		Status              string        `json:"status,omitempty"`
		SubStatus           []interface{} `json:"sub_status,omitempty"`
		Tags                []string      `json:"tags,omitempty"`
		Warranty            string        `json:"warranty,omitempty"`
		CatalogProductID    interface{}   `json:"catalog_product_id,omitempty"`
		DomainID            string        `json:"domain_id,omitempty"`
		SellerCustomField   interface{}   `json:"seller_custom_field,omitempty"`
		ParentItemID        interface{}   `json:"parent_item_id,omitempty"`
		DifferentialPricing interface{}   `json:"differential_pricing,omitempty"`
		DealIds             []interface{} `json:"deal_ids,omitempty"`
		AutomaticRelist     bool          `json:"automatic_relist,omitempty"`
		DateCreated         time.Time     `json:"date_created,omitempty"`
		LastUpdated         time.Time     `json:"last_updated,omitempty"`
		Health              float64       `json:"health,omitempty"`
		CatalogListing      bool          `json:"catalog_listing,omitempty"`
		ItemRelations       []interface{} `json:"item_relations,omitempty"`
		Channels            []string      `json:"channels,omitempty"`
	} `json:"body,omitempty"`
}

type Announcement struct {
	ID                string      `json:"id,omitempty"`
	Message           string      `json:"message,omitempty"`
	Error             string      `json:"error,omitempty"`
	SiteID            string      `json:"site_id,omitempty"`
	Title             string      `json:"title,omitempty"`
	Subtitle          interface{} `json:"subtitle,omitempty"`
	SellerID          int         `json:"seller_id,omitempty"`
	CategoryID        string      `json:"category_id,omitempty"`
	OfficialStoreID   interface{} `json:"official_store_id,omitempty"`
	Price             float64     `json:"price,omitempty"`
	BasePrice         float64     `json:"base_price,omitempty"`
	OriginalPrice     float64     `json:"original_price,omitempty"`
	InventoryID       interface{} `json:"inventory_id,omitempty"`
	CurrencyID        string      `json:"currency_id,omitempty"`
	InitialQuantity   int         `json:"initial_quantity,omitempty"`
	AvailableQuantity int         `json:"available_quantity,omitempty"`
	SoldQuantity      int         `json:"sold_quantity,omitempty"`
	SaleTerms         []struct {
		ID          string      `json:"id,omitempty"`
		Name        string      `json:"name,omitempty"`
		ValueID     interface{} `json:"value_id,omitempty"`
		ValueName   string      `json:"value_name,omitempty"`
		ValueStruct struct {
			Number int    `json:"number,omitempty"`
			Unit   string `json:"unit,omitempty"`
		} `json:"value_struct,omitempty"`
		Values []struct {
			ID     interface{} `json:"id,omitempty"`
			Name   string      `json:"name,omitempty"`
			Struct struct {
				Number int    `json:"number,omitempty"`
				Unit   string `json:"unit,omitempty"`
			} `json:"struct,omitempty"`
		} `json:"values,omitempty"`
		ValueType string `json:"value_type,omitempty"`
	} `json:"sale_terms,omitempty"`
	BuyingMode      string    `json:"buying_mode,omitempty"`
	ListingTypeID   string    `json:"listing_type_id,omitempty"`
	StartTime       time.Time `json:"start_time,omitempty"`
	StopTime        time.Time `json:"stop_time,omitempty"`
	EndTime         time.Time `json:"end_time,omitempty"`
	ExpirationTime  time.Time `json:"expiration_time,omitempty"`
	Condition       string    `json:"condition,omitempty"`
	Permalink       string    `json:"permalink,omitempty"`
	ThumbnailID     string    `json:"thumbnail_id,omitempty"`
	Thumbnail       string    `json:"thumbnail,omitempty"`
	SecureThumbnail string    `json:"secure_thumbnail,omitempty"`
	Pictures        []struct {
		ID        string `json:"id,omitempty"`
		URL       string `json:"url,omitempty"`
		SecureURL string `json:"secure_url,omitempty"`
		Size      string `json:"size,omitempty"`
		MaxSize   string `json:"max_size,omitempty"`
		Quality   string `json:"quality,omitempty"`
	} `json:"pictures,omitempty"`
	VideoID                      interface{}   `json:"video_id,omitempty"`
	Descriptions                 []interface{} `json:"descriptions,omitempty"`
	AcceptsMercadopago           bool          `json:"accepts_mercadopago,omitempty"`
	NonMercadoPagoPaymentMethods []interface{} `json:"non_mercado_pago_payment_methods,omitempty"`
	Shipping                     struct {
		Mode         string        `json:"mode,omitempty"`
		Methods      []interface{} `json:"methods,omitempty"`
		Tags         []string      `json:"tags,omitempty"`
		Dimensions   interface{}   `json:"dimensions,omitempty"`
		LocalPickUp  bool          `json:"local_pick_up,omitempty"`
		FreeShipping bool          `json:"free_shipping,omitempty"`
		LogisticType string        `json:"logistic_type,omitempty"`
		StorePickUp  bool          `json:"store_pick_up,omitempty"`
	} `json:"shipping,omitempty"`
	InternationalDeliveryMode string `json:"international_delivery_mode,omitempty"`
	SellerAddress             struct {
		Comment     string `json:"comment,omitempty"`
		AddressLine string `json:"address_line,omitempty"`
		ZipCode     string `json:"zip_code,omitempty"`
		City        struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"city,omitempty"`
		State struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"state,omitempty"`
		Country struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"country,omitempty"`
		SearchLocation struct {
			City struct {
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			} `json:"city,omitempty"`
			State struct {
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			} `json:"state,omitempty"`
		} `json:"search_location,omitempty"`
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
		ID        int     `json:"id,omitempty"`
	} `json:"seller_address,omitempty"`
	SellerContact interface{} `json:"seller_contact,omitempty"`
	Location      struct {
	} `json:"location,omitempty"`
	Geolocation struct {
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	} `json:"geolocation,omitempty"`
	CoverageAreas []interface{} `json:"coverage_areas,omitempty"`
	Attributes    []struct {
		ID          string      `json:"id,omitempty"`
		Name        string      `json:"name,omitempty"`
		ValueID     string      `json:"value_id,omitempty"`
		ValueName   string      `json:"value_name,omitempty"`
		ValueStruct interface{} `json:"value_struct,omitempty"`
		Values      []struct {
			ID     string      `json:"id,omitempty"`
			Name   string      `json:"name,omitempty"`
			Struct interface{} `json:"struct,omitempty"`
		} `json:"values,omitempty"`
		AttributeGroupID   string `json:"attribute_group_id,omitempty"`
		AttributeGroupName string `json:"attribute_group_name,omitempty"`
		ValueType          string `json:"value_type,omitempty"`
	} `json:"attributes,omitempty"`
	Warnings            []interface{} `json:"warnings,omitempty"`
	ListingSource       string        `json:"listing_source,omitempty"`
	Variations          []interface{} `json:"variations,omitempty"`
	Status              string        `json:"status,omitempty"`
	SubStatus           []interface{} `json:"sub_status,omitempty"`
	Tags                []string      `json:"tags,omitempty"`
	Warranty            string        `json:"warranty,omitempty"`
	CatalogProductID    interface{}   `json:"catalog_product_id,omitempty"`
	DomainID            string        `json:"domain_id,omitempty"`
	SellerCustomField   interface{}   `json:"seller_custom_field,omitempty"`
	ParentItemID        interface{}   `json:"parent_item_id,omitempty"`
	DifferentialPricing interface{}   `json:"differential_pricing,omitempty"`
	DealIds             []interface{} `json:"deal_ids,omitempty"`
	AutomaticRelist     bool          `json:"automatic_relist,omitempty"`
	DateCreated         time.Time     `json:"date_created,omitempty"`
	LastUpdated         time.Time     `json:"last_updated,omitempty"`
	Health              float64       `json:"health,omitempty"`
	CatalogListing      bool          `json:"catalog_listing,omitempty"`
	ItemRelations       []interface{} `json:"item_relations,omitempty"`
	Channels            []string      `json:"channels,omitempty"`
}

type Description struct {
	Text        string        `json:"text,omitempty"`
	PlainText   string        `json:"plain_text,omitempty"`
	LastUpdated time.Time     `json:"last_updated,omitempty"`
	DateCreated time.Time     `json:"date_created,omitempty"`
	Message     string        `json:"message,omitempty"`
	Error       string        `json:"error,omitempty"`
	Status      int           `json:"status,omitempty"`
	Cause       []interface{} `json:"cause,omitempty"`
	Snapshot    struct {
		URL    string `json:"url,omitempty"`
		Width  int    `json:"width,omitempty"`
		Height int    `json:"height,omitempty"`
		Status string `json:"status,omitempty"`
	} `json:"snapshot,omitempty"`
}

type MeliError struct {
	Message string        `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
	Status  int           `json:"status,omitempty"`
	Cause   []interface{} `json:"cause,omitempty"`
}
