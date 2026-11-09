package dto

type CreateAccountRequest struct {
	Code     string `json:"code" validate:"required,max=20"`
	Name     string `json:"name" validate:"required,max=100"`
	Type     string `json:"type" validate:"required,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	ParentId string `json:"parent_id" validate:"omitempty,max=36"`
	IsActive *bool  `json:"is_active"`
}

type UpdateAccountRequest struct {
	Code     string `json:"code" validate:"required,max=20"`
	Name     string `json:"name" validate:"required,max=100"`
	Type     string `json:"type" validate:"required,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	ParentId string `json:"parent_id" validate:"omitempty,max=36"`
	IsActive *bool  `json:"is_active"`
}

type AccountListQuery struct {
	ListQuery
	Type     string `query:"type"`
	IsActive *bool  `query:"is_active"`
}

type AccountResponse struct {
	Id            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	NormalBalance string `json:"normal_balance"`
	ParentId      string `json:"parent_id,omitempty"`
	ParentCode    string `json:"parent_code,omitempty"`
	IsActive      bool   `json:"is_active"`
}

type PaginatedAccountList struct {
	Items []AccountResponse `json:"items"`
	Meta  PaginationMeta    `json:"meta"`
}
