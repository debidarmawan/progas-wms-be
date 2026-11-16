package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToCylinderLocationResponse(location *model.CylinderLocation) *dto.CylinderLocationResponse {
	return &dto.CylinderLocationResponse{
		Id:       location.Id,
		Code:     location.Code,
		Name:     location.Name,
		Type:     string(location.Type),
		Address:  location.Address,
		Notes:    location.Notes,
		IsActive: location.IsActive,
	}
}

func ToCylinderLocationResponses(locations []model.CylinderLocation) []dto.CylinderLocationResponse {
	responses := make([]dto.CylinderLocationResponse, 0, len(locations))
	for i := range locations {
		responses = append(responses, *ToCylinderLocationResponse(&locations[i]))
	}
	return responses
}
