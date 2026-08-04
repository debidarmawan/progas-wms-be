package dto

type CreateCylinderRequest struct {
	BarcodeSN         string `json:"barcode_sn" validate:"required"`
	ItemId            string `json:"item_id" validate:"required"`
	OwnershipType     string `json:"ownership_type" validate:"required"`
	OwnerId           string `json:"owner_id"`
	LastHydrotestDate string `json:"last_hydrotest_date" validate:"required"`
	Remarks           string `json:"remarks"`
}

type UpdateCylinderRequest struct {
	BarcodeSN         string `json:"barcode_sn" validate:"required"`
	ItemId            string `json:"item_id" validate:"required"`
	OwnershipType     string `json:"ownership_type" validate:"required"`
	OwnerId           string `json:"owner_id"`
	LastHydrotestDate string `json:"last_hydrotest_date" validate:"required"`
	Remarks           string `json:"remarks"`
}

type CylinderResponse struct {
	Id                string `json:"id"`
	BarcodeSN         string `json:"barcode_sn"`
	ItemId            string `json:"item_id"`
	ItemName          string `json:"item_name"`
	GasType           string `json:"gas_type"`
	OwnershipType     string `json:"ownership_type"`
	OwnerId           string `json:"owner_id,omitempty"`
	OwnerName         string `json:"owner_name,omitempty"`
	Status            string `json:"status"`
	LastHydrotestDate string `json:"last_hydrotest_date"`
	Remarks           string `json:"remarks,omitempty"`
}

type CylinderHistoryChange struct {
	Field string `json:"field"`
	Label string `json:"label"`
	Old   string `json:"old"`
	New   string `json:"new"`
}

type CylinderHistoryEntry struct {
	Id          string                  `json:"id"`
	Action      string                  `json:"action"`
	ActionLabel string                  `json:"action_label"`
	UserId      string                  `json:"user_id,omitempty"`
	UserName    string                  `json:"user_name,omitempty"`
	CreatedAt   string                  `json:"created_at"`
	Changes     []CylinderHistoryChange `json:"changes"`
}

type CylinderHistoryResponse struct {
	CylinderId string                 `json:"cylinder_id"`
	BarcodeSN  string                 `json:"barcode_sn"`
	Entries    []CylinderHistoryEntry `json:"entries"`
}
