package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToUomResponse(u *model.Uom) *dto.UomResponse {
	return &dto.UomResponse{
		Id:   u.Id,
		Code: u.Code,
		Name: u.Name,
	}
}

func ToUomResponses(uoms []model.Uom) []dto.UomResponse {
	responses := make([]dto.UomResponse, 0, len(uoms))
	for i := range uoms {
		responses = append(responses, *ToUomResponse(&uoms[i]))
	}
	return responses
}
