package dto

type LowStockSparepartAlert struct {
	ItemId    string `json:"item_id"`
	ItemName  string `json:"item_name"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	MinStock  int    `json:"min_stock"`
}

type CustomerQuotaAlert struct {
	CustomerId       string `json:"customer_id"`
	CustomerCode     string `json:"customer_code"`
	CustomerName     string `json:"customer_name"`
	OutstandingCount int    `json:"outstanding_count"`
	QuotaLimit       int    `json:"quota_limit"`
}

type DashboardSummaryResponse struct {
	CylindersByStatus        map[string]int          `json:"cylinders_by_status"`
	LowStockSpareparts       []LowStockSparepartAlert `json:"low_stock_spareparts"`
	TotalOutstandingCylinders int                    `json:"total_outstanding_cylinders"`
	CustomersOverQuota       []CustomerQuotaAlert    `json:"customers_over_quota"`
	HydrotestExpiredCount    int                     `json:"hydrotest_expired_count"`
	HydrotestDueSoonCount    int                     `json:"hydrotest_due_soon_count"`
}
