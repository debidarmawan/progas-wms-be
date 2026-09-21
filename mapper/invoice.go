package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToInvoiceResponse(invoice *model.Invoice, includePayments bool) *dto.InvoiceResponse {
	customerName := ""
	if invoice.Customer.Id != "" {
		customerName = invoice.Customer.Name
	}
	doNumber := ""
	if invoice.DeliveryOrder.Id != "" {
		doNumber = invoice.DeliveryOrder.DONumber
	}

	res := &dto.InvoiceResponse{
		Id:              invoice.Id,
		InvoiceNumber:   invoice.InvoiceNumber,
		DeliveryOrderId: invoice.DeliveryOrderId,
		DONumber:        doNumber,
		CustomerId:      invoice.CustomerId,
		CustomerName:    customerName,
		InvoiceDate:     invoice.InvoiceDate.Format(time.RFC3339),
		DueDate:         invoice.DueDate.Format(time.RFC3339),
		Status:          string(invoice.Status),
		TotalAmount:     invoice.TotalAmount,
		PaidAmount:      invoice.PaidAmount,
		Notes:           invoice.Notes,
		CreatedAt:       invoice.CreatedAt.Format(time.RFC3339),
	}

	if includePayments {
		for _, p := range invoice.Payments {
			res.Payments = append(res.Payments, dto.PaymentResponse{
				Id:          p.Id,
				Amount:      p.Amount,
				PaidAt:      p.PaidAt.Format(time.RFC3339),
				Method:      p.Method,
				ReferenceNo: p.ReferenceNo,
				CreatedAt:   p.CreatedAt.Format(time.RFC3339),
			})
		}
	}
	return res
}

func ToInvoiceResponses(invoices []model.Invoice) []dto.InvoiceResponse {
	res := make([]dto.InvoiceResponse, 0, len(invoices))
	for i := range invoices {
		res = append(res, *ToInvoiceResponse(&invoices[i], false))
	}
	return res
}
