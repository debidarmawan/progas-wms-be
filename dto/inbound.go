package dto

type BarcodeListRequest struct {
	Barcodes []string `json:"barcodes" validate:"required,min=1,dive,required"`
}

type BarcodeOperationResponse struct {
	ProcessedCount int      `json:"processed_count"`
	Barcodes       []string `json:"barcodes"`
}
