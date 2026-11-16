package dto

type CreateCylinderLocationRequest struct {
	Code    string `json:"code" validate:"required,max=32"`
	Name    string `json:"name" validate:"required,max=128"`
	Type    string `json:"type" validate:"required,oneof=WAREHOUSE FILLING_STATION WORKSHOP"`
	Address string `json:"address" validate:"max=255"`
	Notes   string `json:"notes" validate:"max=255"`
}

type UpdateCylinderLocationRequest struct {
	Code     string `json:"code" validate:"required,max=32"`
	Name     string `json:"name" validate:"required,max=128"`
	Type     string `json:"type" validate:"required,oneof=WAREHOUSE FILLING_STATION WORKSHOP"`
	Address  string `json:"address" validate:"max=255"`
	Notes    string `json:"notes" validate:"max=255"`
	IsActive *bool  `json:"is_active"`
}

type CylinderLocationResponse struct {
	Id       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address,omitempty"`
	Notes    string `json:"notes,omitempty"`
	IsActive bool   `json:"is_active"`
}

type PaginatedCylinderLocationList struct {
	Items []CylinderLocationResponse `json:"items"`
	Meta  PaginationMeta             `json:"meta"`
}
