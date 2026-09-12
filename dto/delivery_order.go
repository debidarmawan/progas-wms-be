package dto

type IssueDeliveryOrderRequest struct {
	CustomerId   string   `json:"customer_id" validate:"required"`
	FleetId      string   `json:"fleet_id" validate:"required"`
	SalesOrderId string   `json:"sales_order_id"`
	Barcodes     []string `json:"barcodes" validate:"required,min=1,dive,required"`
	Notes        string   `json:"notes"`
}

type DeliveryOrderDetailResponse struct {
	Id           string  `json:"id"`
	CylinderId   string  `json:"cylinder_id"`
	BarcodeSN    string  `json:"barcode_sn"`
	WeightKg     float64 `json:"weight_kg"`
	MasterItemId string  `json:"master_item_id"`
	UnitPrice    float64 `json:"unit_price"`
	PriceSource  string  `json:"price_source"`
	LineTotal    float64 `json:"line_total"`
}

type DeliveryOrderResponse struct {
	Id            string                        `json:"id"`
	DONumber      string                        `json:"do_number"`
	SalesOrderId  string                        `json:"sales_order_id,omitempty"`
	CustomerId    string                        `json:"customer_id"`
	CustomerName  string                        `json:"customer_name"`
	FleetId       string                        `json:"fleet_id"`
	PlateNumber   string                        `json:"plate_number"`
	Status        string                        `json:"status"`
	TotalWeightKg float64                       `json:"total_weight_kg"`
	CylinderQty   int                           `json:"cylinder_qty"`
	Notes         string                        `json:"notes"`
	CreatedAt     string                        `json:"created_at"`
	Details       []DeliveryOrderDetailResponse `json:"details,omitempty"`
}

type PaginatedDeliveryOrderList struct {
	Items []DeliveryOrderResponse `json:"items"`
	Meta  PaginationMeta          `json:"meta"`
}
