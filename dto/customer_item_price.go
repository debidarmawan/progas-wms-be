package dto

type CreateCustomerItemPriceRequest struct {
	MasterItemId  string  `json:"master_item_id" validate:"required"`
	Price         float64 `json:"price" validate:"gt=0"`
	EffectiveFrom string  `json:"effective_from" validate:"required"`
	EffectiveTo   *string `json:"effective_to"`
}

type UpdateCustomerItemPriceRequest struct {
	Price         float64 `json:"price" validate:"gt=0"`
	EffectiveFrom string  `json:"effective_from" validate:"required"`
	EffectiveTo   *string `json:"effective_to"`
}

type CustomerItemPriceResponse struct {
	Id            string  `json:"id"`
	CustomerId    string  `json:"customer_id"`
	MasterItemId  string  `json:"master_item_id"`
	ItemName      string  `json:"item_name"`
	ItemSKU       string  `json:"item_sku"`
	Price         float64 `json:"price"`
	EffectiveFrom string  `json:"effective_from"`
	EffectiveTo   *string `json:"effective_to,omitempty"`
	IsActive      bool    `json:"is_active"`
	IsFallback    bool    `json:"is_fallback"`
}
