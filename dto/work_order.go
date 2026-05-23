package dto

type WorkOrderSparepartLine struct {
	ItemId   string `json:"item_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}

type CreateWorkOrderRequest struct {
	Title       string                   `json:"title" validate:"required"`
	Description string                   `json:"description"`
	Spareparts  []WorkOrderSparepartLine `json:"spareparts" validate:"required,min=1,dive"`
}

type WorkOrderSparepartResponse struct {
	ItemId   string `json:"item_id"`
	ItemName string `json:"item_name"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type WorkOrderResponse struct {
	Id          string                       `json:"id"`
	WONumber    string                       `json:"wo_number"`
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Status      string                       `json:"status"`
	CreatedAt   string                       `json:"created_at"`
	Spareparts  []WorkOrderSparepartResponse `json:"spareparts,omitempty"`
}

type PaginatedWorkOrderList struct {
	Items []WorkOrderResponse `json:"items"`
	Meta  PaginationMeta      `json:"meta"`
}
