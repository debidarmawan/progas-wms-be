package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToCustomerResponse(customer *model.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		Id:                 customer.Id,
		Code:               customer.Code,
		Name:               customer.Name,
		Phone:              customer.Phone,
		Address:            customer.Address,
		CylinderQuotaLimit: customer.CylinderQuotaLimit,
		OutstandingCount:   customer.OutstandingCount,
		IsActive:           customer.IsActive,
	}
}

func ToCustomerResponses(customers []model.Customer) []dto.CustomerResponse {
	responses := make([]dto.CustomerResponse, 0, len(customers))
	for i := range customers {
		responses = append(responses, *ToCustomerResponse(&customers[i]))
	}
	return responses
}
