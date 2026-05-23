package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToFleetResponse(f *model.FleetVehicle) *dto.FleetResponse {
	return &dto.FleetResponse{
		Id:          f.Id,
		PlateNumber: f.PlateNumber,
		DriverName:  f.DriverName,
		MaxWeightKg: f.MaxWeightKg,
		IsActive:    f.IsActive,
	}
}

func ToFleetResponses(fleets []model.FleetVehicle) []dto.FleetResponse {
	responses := make([]dto.FleetResponse, 0, len(fleets))
	for i := range fleets {
		responses = append(responses, *ToFleetResponse(&fleets[i]))
	}
	return responses
}
