package dto

type CreateCylinderRequest struct {
	BarcodeSN         string `json:"barcode_sn" validate:"required"`
	ItemId            string `json:"item_id" validate:"required"`
	OwnershipType     string `json:"ownership_type" validate:"required"`
	OwnerId           string `json:"owner_id"`
	LastHydrotestDate string `json:"last_hydrotest_date" validate:"required"`
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
}
