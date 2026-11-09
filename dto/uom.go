package dto

type CreateUomRequest struct {
	Code string `json:"code" validate:"required,max=20"`
	Name string `json:"name" validate:"required,max=50"`
}

type UpdateUomRequest struct {
	Code string `json:"code" validate:"required,max=20"`
	Name string `json:"name" validate:"required,max=50"`
}

type UomResponse struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type PaginatedUomList struct {
	Items []UomResponse  `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
