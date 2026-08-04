package dto

type CreateDriverRequest struct {
	Name          string `json:"name" validate:"required"`
	Phone         string `json:"phone"`
	LicenseNumber string `json:"license_number"`
}

type UpdateDriverRequest struct {
	Name          string `json:"name" validate:"required"`
	Phone         string `json:"phone"`
	LicenseNumber string `json:"license_number"`
	IsActive      bool   `json:"is_active"`
}

type DriverResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	LicenseNumber string `json:"license_number"`
	IsActive      bool   `json:"is_active"`
}

type PaginatedDriverList struct {
	Items []DriverResponse `json:"items"`
	Meta  PaginationMeta   `json:"meta"`
}
