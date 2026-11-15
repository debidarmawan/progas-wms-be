package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToCustomerResponse(customer *model.Customer) *dto.CustomerResponse {
	res := &dto.CustomerResponse{
		Id:                 customer.Id,
		Code:               customer.Code,
		Name:               customer.Name,
		Pic:                customer.Pic,
		Npwp:               customer.Npwp,
		Fax:                customer.Fax,
		Email:              customer.Email,
		Phone:              customer.Phone,
		Address:            customer.Address,
		PostalCode:         customer.PostalCode,
		CylinderQuotaLimit: customer.CylinderQuotaLimit,
		CreditLimit:        customer.CreditLimit,
		TermOfPaymentId:    customer.TermOfPaymentId,
		PaymentTermDays:    customer.PaymentTermDays,
		SalesPicName:       customer.SalesPicName,
		SalesPicPhone:      customer.SalesPicPhone,
		SalesPicEmail:      customer.SalesPicEmail,
		OutstandingCount:   customer.OutstandingCount,
		IsActive:           customer.IsActive,
	}
	if customer.TermOfPayment != nil {
		res.TermOfPaymentCode = customer.TermOfPayment.Code
		res.TermOfPaymentName = customer.TermOfPayment.Name
	}
	return res
}

func ToCustomerResponses(customers []model.Customer) []dto.CustomerResponse {
	responses := make([]dto.CustomerResponse, 0, len(customers))
	for i := range customers {
		responses = append(responses, *ToCustomerResponse(&customers[i]))
	}
	return responses
}
