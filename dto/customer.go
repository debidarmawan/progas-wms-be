package dto

type CreateCustomerRequest struct {
	Name               string `json:"name" validate:"required"`
	Pic                string `json:"pic"`
	Npwp               string `json:"npwp"`
	Fax                string `json:"fax"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	CylinderQuotaLimit int    `json:"cylinder_quota_limit" validate:"gte=0"`
}

type UpdateCustomerRequest struct {
	Name               string `json:"name" validate:"required"`
	Pic                string `json:"pic"`
	Npwp               string `json:"npwp"`
	Fax                string `json:"fax"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	CylinderQuotaLimit int    `json:"cylinder_quota_limit" validate:"gte=0"`
	IsActive           bool   `json:"is_active"`
}

type CustomerResponse struct {
	Id                 string `json:"id"`
	Code               string `json:"code"`
	Name               string `json:"name"`
	Pic                string `json:"pic"`
	Npwp               string `json:"npwp"`
	Fax                string `json:"fax"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	CylinderQuotaLimit int    `json:"cylinder_quota_limit"`
	OutstandingCount   int    `json:"outstanding_count"`
	IsActive           bool   `json:"is_active"`
}
