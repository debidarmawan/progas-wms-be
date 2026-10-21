package dto

type CreateCylinderRequest struct {
	BarcodeSN         string `json:"barcode_sn"` // optional: auto-generated when empty
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
	LabelPrintedAt    string `json:"label_printed_at,omitempty"`
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

// ScanCylinderQuery is the menu being scanned into (pre-fill, do, exchange-in, ...).
type ScanCylinderQuery struct {
	Menu string `json:"menu" validate:"required,oneof=pre-fill filling post-fill do exchange-out exchange-in hydrotest"`
}

type ScanResolveResponse struct {
	Id            string   `json:"id"`
	BarcodeSN     string   `json:"barcode_sn"`
	ItemId        string   `json:"item_id"`
	ItemName      string   `json:"item_name"`
	GasType       string   `json:"gas_type"`
	OwnershipType string   `json:"ownership_type"`
	Status        string   `json:"status"`
	Menu          string   `json:"menu"`
	Allowed       bool     `json:"allowed"`
	Reason        string   `json:"reason,omitempty"`
	AllowedMenus  []string `json:"allowed_menus"`
}

type CylinderLabelQuery struct {
	Ids     string `json:"ids"` // comma-separated cylinder ids
	Labeled string `json:"labeled" validate:"omitempty,oneof=true false"`
	Limit   int    `json:"limit" validate:"omitempty,min=1,max=500"`
	Format  string `json:"format" validate:"omitempty,oneof=json csv"`
}

type CylinderLabelResponse struct {
	Id        string `json:"id"`
	BarcodeSN string `json:"barcode_sn"`
	ItemName  string `json:"item_name"`
	GasType   string `json:"gas_type"`
}

type MarkLabelsPrintedRequest struct {
	Ids []string `json:"ids" validate:"required,min=1,dive,required"`
}
