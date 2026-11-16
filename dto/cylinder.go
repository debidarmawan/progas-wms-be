package dto

// Manual location input only accepts real places/parties; IN_TRANSIT, MAINTENANCE and
// UNKNOWN are set by the flow usecases, never by hand.
type CreateCylinderRequest struct {
	BarcodeSN         string `json:"barcode_sn"` // optional: auto-generated when empty
	ItemId            string `json:"item_id" validate:"required"`
	OwnershipType     string `json:"ownership_type" validate:"required"`
	OwnerId           string `json:"owner_id"`
	LastHydrotestDate string `json:"last_hydrotest_date" validate:"required"`
	Remarks           string `json:"remarks"`
	LocationType      string `json:"location_type" validate:"omitempty,oneof=WAREHOUSE FILLING_STATION WORKSHOP CUSTOMER VENDOR"`
	LocationId        string `json:"location_id"`
	// Pointer: nil inherits MasterItem.MaxDaysAtCustomer, 0 means "no limit".
	MaxDaysAtCustomer *int `json:"max_days_at_customer" validate:"omitempty,gte=0,lte=3650"`
}

type UpdateCylinderRequest struct {
	BarcodeSN         string `json:"barcode_sn" validate:"required"`
	ItemId            string `json:"item_id" validate:"required"`
	OwnershipType     string `json:"ownership_type" validate:"required"`
	OwnerId           string `json:"owner_id"`
	LastHydrotestDate string `json:"last_hydrotest_date" validate:"required"`
	Remarks           string `json:"remarks"`
	LocationType      string `json:"location_type" validate:"omitempty,oneof=WAREHOUSE FILLING_STATION WORKSHOP CUSTOMER VENDOR"`
	LocationId        string `json:"location_id"`
	// Pointer: nil keeps the current value, 0 means "no limit".
	MaxDaysAtCustomer *int `json:"max_days_at_customer" validate:"omitempty,gte=0,lte=3650"`
}

// UpdateCylinderMaxDaysRequest allows changing the aging limit for any cylinder that is
// still in the lifecycle, including IN_TRANSIT/OUTSTANDING which PUT /cylinders/:id rejects.
type UpdateCylinderMaxDaysRequest struct {
	MaxDaysAtCustomer *int `json:"max_days_at_customer" validate:"required,gte=0,lte=3650"`
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
	LocationType      string `json:"location_type"`
	LocationId        string `json:"location_id,omitempty"`
	LocationName      string `json:"location_name,omitempty"`
	MaxDaysAtCustomer int    `json:"max_days_at_customer"`
	OutstandingSince  string `json:"outstanding_since,omitempty"`
	DaysAtCustomer    int    `json:"days_at_customer"`
	IsOverdue         bool   `json:"is_overdue"`
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
