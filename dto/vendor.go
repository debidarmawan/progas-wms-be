package dto

type CreateVendorRequest struct {
	Code              string `json:"code" validate:"required"`
	Name              string `json:"name" validate:"required"`
	ContactPerson     string `json:"contact_person"`
	Phone             string `json:"phone"`
	Email             string `json:"email" validate:"omitempty,email"`
	Address           string `json:"address"`
	Notes             string `json:"notes"`
	ContractStartDate string `json:"contract_start_date"`
	ContractEndDate   string `json:"contract_end_date"`
}

type UpdateVendorRequest struct {
	Name              string `json:"name" validate:"required"`
	ContactPerson     string `json:"contact_person"`
	Phone             string `json:"phone"`
	Email             string `json:"email" validate:"omitempty,email"`
	Address           string `json:"address"`
	Notes             string `json:"notes"`
	ContractStartDate string `json:"contract_start_date"`
	ContractEndDate   string `json:"contract_end_date"`
	IsActive          bool   `json:"is_active"`
}

type VendorCylinderSummary struct {
	Id        string `json:"id"`
	BarcodeSN string `json:"barcode_sn"`
	Status    string `json:"status"`
	GasType   string `json:"gas_type"`
	ItemName  string `json:"item_name"`
}

type VendorResponse struct {
	Id                string `json:"id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	ContactPerson     string `json:"contact_person"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Address           string `json:"address"`
	Notes             string `json:"notes"`
	ContractStartDate string `json:"contract_start_date,omitempty"`
	ContractEndDate   string `json:"contract_end_date,omitempty"`
	IsActive          bool   `json:"is_active"`
	CylinderCount     int    `json:"cylinder_count"`
}

type VendorDetailResponse struct {
	VendorResponse
	CylindersByStatus map[string]int           `json:"cylinders_by_status"`
	Cylinders         []VendorCylinderSummary `json:"cylinders"`
}

type PaginatedVendorList struct {
	Items []VendorResponse `json:"items"`
	Meta  PaginationMeta   `json:"meta"`
}
