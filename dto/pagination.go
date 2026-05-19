package dto

// ListQuery is used for paginated list endpoints with optional search.
type ListQuery struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Items []T            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// Concrete list types for Swagger (swag does not support nested generics).
type PaginatedRoleList struct {
	Items []RoleResponse `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginatedMasterItemList struct {
	Items []MasterItemResponse `json:"items"`
	Meta  PaginationMeta       `json:"meta"`
}

type PaginatedCylinderList struct {
	Items []CylinderResponse `json:"items"`
	Meta  PaginationMeta     `json:"meta"`
}

type PaginatedCustomerList struct {
	Items []CustomerResponse `json:"items"`
	Meta  PaginationMeta     `json:"meta"`
}
