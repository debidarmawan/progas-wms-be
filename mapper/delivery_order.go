package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToDeliveryOrderResponse(do *model.DeliveryOrder, includeDetails bool) *dto.DeliveryOrderResponse {
	customerName := ""
	if do.Customer.Id != "" {
		customerName = do.Customer.Name
	}
	plateNumber := ""
	if do.FleetVehicle.Id != "" {
		plateNumber = do.FleetVehicle.PlateNumber
	}

	res := &dto.DeliveryOrderResponse{
		Id:            do.Id,
		DONumber:      do.DONumber,
		CustomerId:    do.CustomerId,
		CustomerName:  customerName,
		FleetId:       do.FleetId,
		PlateNumber:   plateNumber,
		Status:        string(do.Status),
		TotalWeightKg: do.TotalWeightKg,
		CylinderQty:   do.CylinderQty,
		Notes:         do.Notes,
		CreatedAt:     do.CreatedAt.Format(time.RFC3339),
	}

	if includeDetails {
		for _, d := range do.Details {
			res.Details = append(res.Details, dto.DeliveryOrderDetailResponse{
				Id:         d.Id,
				CylinderId: d.CylinderId,
				BarcodeSN:  d.BarcodeSN,
				WeightKg:   d.WeightKg,
			})
		}
	}
	return res
}

func ToDeliveryOrderResponses(orders []model.DeliveryOrder) []dto.DeliveryOrderResponse {
	responses := make([]dto.DeliveryOrderResponse, 0, len(orders))
	for i := range orders {
		responses = append(responses, *ToDeliveryOrderResponse(&orders[i], false))
	}
	return responses
}
