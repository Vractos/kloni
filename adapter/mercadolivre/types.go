package mercadolivre

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
	TotalAmount int           `json:"total_amount,omitempty"`
	PaidAmount  int           `json:"paid_amount,omitempty"`
	Coupon      struct {
		ID     interface{} `json:"id,omitempty"`
		Amount int         `json:"amount,omitempty"`
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
		UnitPrice         int         `json:"unit_price,omitempty"`
		FullUnitPrice     int         `json:"full_unit_price,omitempty"`
		CurrencyID        string      `json:"currency_id,omitempty"`
		ManufacturingDays interface{} `json:"manufacturing_days,omitempty"`
		SaleFee           int         `json:"sale_fee,omitempty"`
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
		TransactionAmount         int         `json:"transaction_amount,omitempty"`
		TransactionAmountRefunded int         `json:"transaction_amount_refunded,omitempty"`
		TaxesAmount               int         `json:"taxes_amount,omitempty"`
		ShippingCost              int         `json:"shipping_cost,omitempty"`
		CouponAmount              int         `json:"coupon_amount,omitempty"`
		OverpaidAmount            int         `json:"overpaid_amount,omitempty"`
		TotalPaidAmount           int         `json:"total_paid_amount,omitempty"`
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
		Amount     interface{} `json:"amount,omitempty"`
		CurrencyID interface{} `json:"currency_id,omitempty"`
		ID         interface{} `json:"id,omitempty"`
	} `json:"taxes,omitempty"`
}
