package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"time"
)

func ToCylinderResponse(cylinder *model.Cylinder) *dto.CylinderResponse {
	itemName := ""
	gasType := ""
	if cylinder.MasterItem.Id != "" {
		itemName = cylinder.MasterItem.Name
		gasType = cylinder.MasterItem.GasType
	}
	ownerId := ""
	if cylinder.OwnerId != nil {
		ownerId = *cylinder.OwnerId
	}
	labelPrintedAt := ""
	if cylinder.LabelPrintedAt != nil {
		labelPrintedAt = cylinder.LabelPrintedAt.Format(time.RFC3339)
	}
	locationId := ""
	if cylinder.LocationId != nil {
		locationId = *cylinder.LocationId
	}
	outstandingSince := ""
	if cylinder.OutstandingSince != nil {
		outstandingSince = cylinder.OutstandingSince.Format(time.RFC3339)
	}
	return &dto.CylinderResponse{
		Id:                cylinder.Id,
		BarcodeSN:         cylinder.BarcodeSN,
		ItemId:            cylinder.ItemId,
		ItemName:          itemName,
		GasType:           gasType,
		OwnershipType:     string(cylinder.OwnershipType),
		OwnerId:           ownerId,
		Status:            string(cylinder.Status),
		LastHydrotestDate: cylinder.LastHydrotestDate.Format(time.RFC3339),
		Remarks:           cylinder.Remarks,
		LabelPrintedAt:    labelPrintedAt,
		LocationType:      string(cylinder.LocationType),
		LocationId:        locationId,
		LocationName:      cylinder.LocationName,
		MaxDaysAtCustomer: cylinder.MaxDaysAtCustomer,
		OutstandingSince:  outstandingSince,
		DaysAtCustomer:    helper.DaysAtCustomer(cylinder.OutstandingSince),
		IsOverdue:         helper.IsOverdueAtCustomer(cylinder.OutstandingSince, cylinder.MaxDaysAtCustomer),
	}
}

func ToCylinderResponses(cylinders []model.Cylinder) []dto.CylinderResponse {
	responses := make([]dto.CylinderResponse, 0, len(cylinders))
	for i := range cylinders {
		responses = append(responses, *ToCylinderResponse(&cylinders[i]))
	}
	return responses
}
