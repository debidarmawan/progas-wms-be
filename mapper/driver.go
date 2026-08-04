package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToDriverResponse(d *model.Driver) *dto.DriverResponse {
	return &dto.DriverResponse{
		Id:            d.Id,
		Name:          d.Name,
		Phone:         d.Phone,
		LicenseNumber: d.LicenseNumber,
		IsActive:      d.IsActive,
	}
}

func ToDriverResponses(drivers []model.Driver) []dto.DriverResponse {
	responses := make([]dto.DriverResponse, 0, len(drivers))
	for i := range drivers {
		responses = append(responses, *ToDriverResponse(&drivers[i]))
	}
	return responses
}
