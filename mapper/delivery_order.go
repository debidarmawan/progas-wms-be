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
	customerSiteId := ""
	if do.CustomerSiteId != nil {
		customerSiteId = *do.CustomerSiteId
	}
	driverId := ""
	driverName := ""
	if do.Driver != nil {
		driverId = do.Driver.Id
		driverName = do.Driver.Name
	}

	res := &dto.DeliveryOrderResponse{
		Id:               do.Id,
		DONumber:         do.DONumber,
		SalesOrderId:     do.SalesOrderId,
		CustomerId:       do.CustomerId,
		CustomerName:     customerName,
		CustomerSiteId:   customerSiteId,
		ShipToName:       do.ShipToName,
		ShipToAddress:    do.ShipToAddress,
		ShipToPostalCode: do.ShipToPostalCode,
		ShipToPhone:      do.ShipToPhone,
		FleetId:          do.FleetId,
		PlateNumber:      plateNumber,
		DriverId:         driverId,
		DriverName:       driverName,
		Status:           string(do.Status),
		TotalWeightKg:    do.TotalWeightKg,
		CylinderQty:      do.CylinderQty,
		Notes:            do.Notes,
		CreatedAt:        do.CreatedAt.Format(time.RFC3339),
	}

	if includeDetails {
		for _, d := range do.Details {
			res.Details = append(res.Details, dto.DeliveryOrderDetailResponse{
				Id:           d.Id,
				CylinderId:   d.CylinderId,
				BarcodeSN:    d.BarcodeSN,
				WeightKg:     d.WeightKg,
				MasterItemId: d.MasterItemId,
				UnitPrice:    d.UnitPrice,
				PriceSource:  d.PriceSource,
				LineTotal:    d.LineTotal,
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
