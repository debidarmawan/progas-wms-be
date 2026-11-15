package dto

type CreateCustomerRequest struct {
	Name               string  `json:"name" validate:"required"`
	Pic                string  `json:"pic"`
	Npwp               string  `json:"npwp"`
	Fax                string  `json:"fax"`
	Email              string  `json:"email"`
	Phone              string  `json:"phone"`
	Address            string  `json:"address"`
	PostalCode         string  `json:"postal_code" validate:"max=16"`
	CylinderQuotaLimit int     `json:"cylinder_quota_limit" validate:"gte=0"`
	CreditLimit        float64 `json:"credit_limit" validate:"gte=0"`
	TermOfPaymentId    string  `json:"term_of_payment_id"`
	SalesPicName       string  `json:"sales_pic_name"`
	SalesPicPhone      string  `json:"sales_pic_phone"`
	SalesPicEmail      string  `json:"sales_pic_email"`
}

type UpdateCustomerRequest struct {
	Name               string  `json:"name" validate:"required"`
	Pic                string  `json:"pic"`
	Npwp               string  `json:"npwp"`
	Fax                string  `json:"fax"`
	Email              string  `json:"email"`
	Phone              string  `json:"phone"`
	Address            string  `json:"address"`
	PostalCode         string  `json:"postal_code" validate:"max=16"`
	CylinderQuotaLimit int     `json:"cylinder_quota_limit" validate:"gte=0"`
	CreditLimit        float64 `json:"credit_limit" validate:"gte=0"`
	TermOfPaymentId    string  `json:"term_of_payment_id"`
	SalesPicName       string  `json:"sales_pic_name"`
	SalesPicPhone      string  `json:"sales_pic_phone"`
	SalesPicEmail      string  `json:"sales_pic_email"`
	IsActive           bool    `json:"is_active"`
}

type CustomerResponse struct {
	Id                 string  `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Pic                string  `json:"pic"`
	Npwp               string  `json:"npwp"`
	Fax                string  `json:"fax"`
	Email              string  `json:"email"`
	Phone              string  `json:"phone"`
	Address            string  `json:"address"`
	PostalCode         string  `json:"postal_code"`
	CylinderQuotaLimit int     `json:"cylinder_quota_limit"`
	CreditLimit        float64 `json:"credit_limit"`
	TermOfPaymentId    *string `json:"term_of_payment_id"`
	TermOfPaymentCode  string  `json:"term_of_payment_code,omitempty"`
	TermOfPaymentName  string  `json:"term_of_payment_name,omitempty"`
	PaymentTermDays    int     `json:"payment_term_days"`
	SalesPicName       string  `json:"sales_pic_name"`
	SalesPicPhone      string  `json:"sales_pic_phone"`
	SalesPicEmail      string  `json:"sales_pic_email"`
	OutstandingCount   int     `json:"outstanding_count"`
	IsActive           bool    `json:"is_active"`
}
