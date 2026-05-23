package dto

type ProcessExchangeRequest struct {
	CustomerId    string   `json:"customer_id" validate:"required"`
	OutBarcodes   []string `json:"out_barcodes" validate:"required,min=1,dive,required"`
	InBarcodes    []string `json:"in_barcodes" validate:"required,min=1,dive,required"`
	ForceApprove  bool     `json:"force_approve"`
}

type ExchangeResponse struct {
	CustomerId         string   `json:"customer_id"`
	OutCount           int      `json:"out_count"`
	InCount            int      `json:"in_count"`
	OutstandingBefore  int      `json:"outstanding_before"`
	OutstandingAfter   int      `json:"outstanding_after"`
	OutstandingDelta   int      `json:"outstanding_delta"`
	OutBarcodes        []string `json:"out_barcodes"`
	InBarcodes         []string `json:"in_barcodes"`
	CrossCustomerAlerts []string `json:"cross_customer_alerts,omitempty"`
}
