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
	FromDate    string             `json:"from_date"`
	ToDate      string             `json:"to_date"`
	SampleCount int                `json:"sample_count"`
	AverageDays float64            `json:"average_days"`
	Samples     []TurnaroundSample `json:"samples"`
}

type TurnaroundSample struct {
	BarcodeSN   string  `json:"barcode_sn"`
	Days        float64 `json:"days"`
	StartedAt   string  `json:"started_at"`
	CompletedAt string  `json:"completed_at"`
}

type VirtualWarehouseCylinder struct {
	BarcodeSN      string `json:"barcode_sn"`
	DaysAtCustomer int    `json:"days_at_customer"`
	MaxDays        int    `json:"max_days"`
	IsOverdue      bool   `json:"is_overdue"`
}

type VirtualWarehouseCustomer struct {
	CustomerId       string                     `json:"customer_id"`
	CustomerCode     string                     `json:"customer_code"`
	CustomerName     string                     `json:"customer_name"`
	OutstandingCount int                        `json:"outstanding_count"`
	OverdueCount     int                        `json:"overdue_count"`
	Cylinders        []VirtualWarehouseCylinder `json:"cylinders"`
}

type VirtualWarehouseResponse struct {
	Customers []VirtualWarehouseCustomer `json:"customers"`
}
