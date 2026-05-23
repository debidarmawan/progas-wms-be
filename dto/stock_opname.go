package dto

type StockOpnameRequest struct {
	ItemId         string `json:"item_id" validate:"required"`
	ActualQuantity int    `json:"actual_quantity" validate:"gte=0"`
	Notes          string `json:"notes"`
}

type StockOpnameResponse struct {
	ItemId         string `json:"item_id"`
	ItemName       string `json:"item_name"`
	QuantityBefore int    `json:"quantity_before"`
	QuantityAfter  int    `json:"quantity_after"`
	QuantityDelta  int    `json:"quantity_delta"`
}
