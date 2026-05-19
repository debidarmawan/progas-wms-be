package dto

type SubmitFillingBatchRequest struct {
	ItemId   string   `json:"item_id" validate:"required"`
	Barcodes []string `json:"barcodes" validate:"required,min=1,dive,required"`
	Notes    string   `json:"notes"`
}

type FillingBatchDetailResponse struct {
	Id         string `json:"id"`
	CylinderId string `json:"cylinder_id"`
	BarcodeSN  string `json:"barcode_sn"`
}

type FillingBatchResponse struct {
	Id          string                     `json:"id"`
	BatchNumber string                     `json:"batch_number"`
	ItemId      string                     `json:"item_id"`
	ItemName    string                     `json:"item_name"`
	GasType     string                     `json:"gas_type"`
	Status      string                     `json:"status"`
	CylinderQty int                        `json:"cylinder_qty"`
	Notes       string                     `json:"notes"`
	CreatedAt   string                     `json:"created_at"`
	Details     []FillingBatchDetailResponse `json:"details,omitempty"`
}

type PaginatedFillingBatchList struct {
	Items []FillingBatchResponse `json:"items"`
	Meta  PaginationMeta         `json:"meta"`
}
