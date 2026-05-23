package dto

type RecordHydrotestRequest struct {
	LastHydrotestDate string `json:"last_hydrotest_date" validate:"required"`
	Notes             string `json:"notes"`
}

type HydrotestDueCylinder struct {
	Id                string `json:"id"`
	BarcodeSN         string `json:"barcode_sn"`
	Status            string `json:"status"`
	LastHydrotestDate string `json:"last_hydrotest_date"`
	ExpiryDate        string `json:"expiry_date"`
	IsExpired         bool   `json:"is_expired"`
}

type HydrotestDueResponse struct {
	DueWithinDays int                    `json:"due_within_days"`
	Items         []HydrotestDueCylinder `json:"items"`
}
