package dto

type CreateCustomerPORequest struct {
	PONumber    string                  `json:"po_number" validate:"required,max=50"`
	CustomerId  string                  `json:"customer_id" validate:"required"`
	PODate      string                  `json:"po_date" validate:"required"`
	ValidUntil  string                  `json:"valid_until"`
	DocumentURL string                  `json:"document_url" validate:"omitempty,url,max=255"`
	Lines       []CustomerPOLineRequest `json:"lines" validate:"required,min=1,dive"`
}

type CustomerPOLineRequest struct {
	MasterItemId string `json:"master_item_id" validate:"required"`
	Quantity     int    `json:"quantity" validate:"required,min=1"`
}

type CustomerPOLineResponse struct {
	Id           string `json:"id"`
	MasterItemId string `json:"master_item_id"`
	SKU          string `json:"sku"`
	ItemName     string `json:"item_name"`
	Quantity     int    `json:"quantity"`
}

type CustomerPOResponse struct {
	Id           string                   `json:"id"`
	PONumber     string                   `json:"po_number"`
	CustomerId   string                   `json:"customer_id"`
	CustomerName string                   `json:"customer_name"`
	PODate       string                   `json:"po_date"`
	ValidUntil   string                   `json:"valid_until,omitempty"`
	DocumentURL  string                   `json:"document_url,omitempty"`
	Status       string                   `json:"status"`
	CreatedAt    string                   `json:"created_at"`
	Lines        []CustomerPOLineResponse `json:"lines,omitempty"`
}
