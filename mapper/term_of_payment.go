package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToTermOfPaymentResponse(top *model.TermOfPayment) *dto.TermOfPaymentResponse {
	return &dto.TermOfPaymentResponse{
		Id:       top.Id,
		Code:     top.Code,
		Name:     top.Name,
		Days:     top.Days,
		IsActive: top.IsActive,
	}
}

func ToTermOfPaymentResponses(tops []model.TermOfPayment) []dto.TermOfPaymentResponse {
	responses := make([]dto.TermOfPaymentResponse, 0, len(tops))
	for i := range tops {
		responses = append(responses, *ToTermOfPaymentResponse(&tops[i]))
	}
	return responses
}
