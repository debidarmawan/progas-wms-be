package dto

type CreateSalesOrderRequest struct {
	CustomerPOId string                  `json:"customer_po_id"`
	CustomerId   string                  `json:"customer_id" validate:"required_without=CustomerPOId"`
	Lines        []SalesOrderLineRequest `json:"lines" validate:"required,min=1,dive"`
}

type SalesOrderLineRequest struct {
	MasterItemId string `json:"master_item_id" validate:"required"`
	QtyOrdered   int    `json:"qty_ordered" validate:"required,min=1"`
}

type SalesOrderLineResponse struct {
	Id           string `json:"id"`
	MasterItemId string `json:"master_item_id"`
	SKU          string `json:"sku"`
	ItemName     string `json:"item_name"`
	QtyOrdered   int    `json:"qty_ordered"`
	QtyDelivered int    `json:"qty_delivered"`
}

type SalesOrderResponse struct {
	Id           string                   `json:"id"`
	SONumber     string                   `json:"so_number"`
	CustomerPOId string                   `json:"customer_po_id,omitempty"`
	CustomerId   string                   `json:"customer_id"`
	CustomerName string                   `json:"customer_name"`
	Status       string                   `json:"status"`
	CreatedAt    string                   `json:"created_at"`
	Lines        []SalesOrderLineResponse `json:"lines,omitempty"`
}
