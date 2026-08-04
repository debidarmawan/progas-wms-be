package dto

type CreateFleetRequest struct {
	PlateNumber string  `json:"plate_number" validate:"required"`
	MaxWeightKg float64 `json:"max_weight_kg" validate:"required,gt=0"`
}

type UpdateFleetRequest struct {
	MaxWeightKg float64 `json:"max_weight_kg" validate:"required,gt=0"`
	IsActive    bool    `json:"is_active"`
}

type FleetResponse struct {
	Id          string  `json:"id"`
	PlateNumber string  `json:"plate_number"`
	MaxWeightKg float64 `json:"max_weight_kg"`
	IsActive    bool    `json:"is_active"`
}

type PaginatedFleetList struct {
	Items []FleetResponse `json:"items"`
	Meta  PaginationMeta  `json:"meta"`
}
