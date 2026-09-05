package dto

type CreateMasterItemRequest struct {
	Name              string  `json:"name" validate:"required"`
	SKU               string  `json:"sku" validate:"required"`
	ItemType          string  `json:"item_type" validate:"required"`
	GasType           string  `json:"gas_type"`
	HnaPrice          float64 `json:"hna_price" validate:"gte=0"`
	IsSerialized      bool    `json:"is_serialized"`
	EmptyWeightKg     float64 `json:"empty_weight_kg" validate:"gte=0"`
	GasWeightKg       float64 `json:"gas_weight_kg" validate:"gte=0"`
	MinStockAlert     int     `json:"min_stock_alert" validate:"gte=0"`
	MaxDaysAtCustomer int     `json:"max_days_at_customer" validate:"gte=0"`
}

type BulkCreateMasterItemRequest struct {
	Items []CreateMasterItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateMasterItemRequest struct {
	Name              string  `json:"name" validate:"required"`
	ItemType          string  `json:"item_type" validate:"required"`
	GasType           string  `json:"gas_type"`
	HnaPrice          float64 `json:"hna_price" validate:"gte=0"`
	EmptyWeightKg     float64 `json:"empty_weight_kg" validate:"gte=0"`
	GasWeightKg       float64 `json:"gas_weight_kg" validate:"gte=0"`
	MinStockAlert     int     `json:"min_stock_alert" validate:"gte=0"`
	MaxDaysAtCustomer int     `json:"max_days_at_customer" validate:"gte=0"`
}

type MasterItemResponse struct {
	Id                string  `json:"id"`
	Name              string  `json:"name"`
	SKU               string  `json:"sku"`
	ItemType          string  `json:"item_type"`
	GasType           string  `json:"gas_type"`
	HnaPrice          float64 `json:"hna_price"`
	IsSerialized      bool    `json:"is_serialized"`
	EmptyWeightKg     float64 `json:"empty_weight_kg"`
	GasWeightKg       float64 `json:"gas_weight_kg"`
	MinStockAlert     int     `json:"min_stock_alert"`
	MaxDaysAtCustomer int     `json:"max_days_at_customer"`
	StockQuantity     *int    `json:"stock_quantity,omitempty"`
}
