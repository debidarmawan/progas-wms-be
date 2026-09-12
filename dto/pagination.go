package dto

// ListQuery is used for paginated list endpoints with optional search and filters.
type ListQuery struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Search       string `query:"search"`
	SortBy       string `query:"sort_by"`
	SortOrder    string `query:"sort_order"`
	ItemType     string `query:"item_type"`
	GasType      string `query:"gas_type"`
	IsSerialized *bool  `query:"is_serialized"`
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

type PaginatedCustomerItemPriceList struct {
	Items []CustomerItemPriceResponse `json:"items"`
	Meta  PaginationMeta              `json:"meta"`
}
