package dto

type InvoiceResponse struct {
	Id              string            `json:"id"`
	InvoiceNumber   string            `json:"invoice_number"`
	DeliveryOrderId string            `json:"delivery_order_id"`
	DONumber        string            `json:"do_number,omitempty"`
	CustomerId      string            `json:"customer_id"`
	CustomerName    string            `json:"customer_name"`
	InvoiceDate     string            `json:"invoice_date"`
	DueDate         string            `json:"due_date"`
	Status          string            `json:"status"`
	TotalAmount     float64           `json:"total_amount"`
	PaidAmount      float64           `json:"paid_amount"`
	Notes           string            `json:"notes,omitempty"`
	CreatedAt       string            `json:"created_at"`
	Payments        []PaymentResponse `json:"payments,omitempty"`
}

type PaymentResponse struct {
	Id          string  `json:"id"`
	Amount      float64 `json:"amount"`
	PaidAt      string  `json:"paid_at"`
	Method      string  `json:"method"`
	ReferenceNo string  `json:"reference_no,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type RecordPaymentRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Method      string  `json:"method" validate:"required"`
	ReferenceNo string  `json:"reference_no"`
}
