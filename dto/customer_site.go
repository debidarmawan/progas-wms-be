package dto

type CreateCustomerSiteRequest struct {
	Name       string `json:"name" validate:"required,max=128"`
	Pic        string `json:"pic" validate:"max=128"`
	Phone      string `json:"phone" validate:"max=32"`
	Address    string `json:"address" validate:"max=255"`
	PostalCode string `json:"postal_code" validate:"max=16"`
	Notes      string `json:"notes" validate:"max=255"`
	IsPrimary  bool   `json:"is_primary"`
	IsActive   *bool  `json:"is_active"`
}

type UpdateCustomerSiteRequest struct {
	Name       string `json:"name" validate:"required,max=128"`
	Pic        string `json:"pic" validate:"max=128"`
	Phone      string `json:"phone" validate:"max=32"`
	Address    string `json:"address" validate:"max=255"`
	PostalCode string `json:"postal_code" validate:"max=16"`
	Notes      string `json:"notes" validate:"max=255"`
	IsPrimary  bool   `json:"is_primary"`
	IsActive   bool   `json:"is_active"`
}

type CustomerSiteResponse struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	Name       string `json:"name"`
	Pic        string `json:"pic"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	Notes      string `json:"notes"`
	IsPrimary  bool   `json:"is_primary"`
	IsActive   bool   `json:"is_active"`
}

type CustomerSiteListResponse struct {
	Items []CustomerSiteResponse `json:"items"`
}
