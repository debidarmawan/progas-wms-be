package dto

type CreateTermOfPaymentRequest struct {
	Code     string `json:"code" validate:"required,max=20"`
	Name     string `json:"name" validate:"required,max=100"`
	Days     int    `json:"days" validate:"gte=0,lte=365"`
	IsActive *bool  `json:"is_active"`
}

type UpdateTermOfPaymentRequest struct {
	Code     string `json:"code" validate:"required,max=20"`
	Name     string `json:"name" validate:"required,max=100"`
	Days     int    `json:"days" validate:"gte=0,lte=365"`
	IsActive *bool  `json:"is_active"`
}

type TermOfPaymentResponse struct {
	Id       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Days     int    `json:"days"`
	IsActive bool   `json:"is_active"`
}

type PaginatedTermOfPaymentList struct {
	Items []TermOfPaymentResponse `json:"items"`
	Meta  PaginationMeta          `json:"meta"`
}
