package dto

type CylinderLedgerEntryResponse struct {
	Id            string `json:"id"`
	BarcodeSN     string `json:"barcode_sn"`
	FromStatus    string `json:"from_status"`
	ToStatus      string `json:"to_status"`
	Action        string `json:"action"`
	ReferenceType string `json:"reference_type"`
	ReferenceId   string `json:"reference_id"`
	CreatedAt     string `json:"created_at"`
}

type StockLedgerReportResponse struct {
	BarcodeSN string                        `json:"barcode_sn"`
	Entries   []CylinderLedgerEntryResponse `json:"entries"`
}

type TurnaroundReportResponse struct {
	FromDate       string             `json:"from_date"`
	ToDate         string             `json:"to_date"`
	SampleCount    int                `json:"sample_count"`
	AverageDays    float64            `json:"average_days"`
	Samples        []TurnaroundSample `json:"samples"`
}

type TurnaroundSample struct {
	BarcodeSN   string  `json:"barcode_sn"`
	Days        float64 `json:"days"`
	StartedAt   string  `json:"started_at"`
	CompletedAt string  `json:"completed_at"`
}

type VirtualWarehouseCustomer struct {
	CustomerId       string   `json:"customer_id"`
	CustomerCode     string   `json:"customer_code"`
	CustomerName     string   `json:"customer_name"`
	OutstandingCount int      `json:"outstanding_count"`
	CylinderBarcodes []string `json:"cylinder_barcodes"`
}

type VirtualWarehouseResponse struct {
	Customers []VirtualWarehouseCustomer `json:"customers"`
}
