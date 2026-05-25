package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToWorkOrderResponse(wo *model.WorkOrder, includeLines bool) *dto.WorkOrderResponse {
	res := &dto.WorkOrderResponse{
		Id:          wo.Id,
		WONumber:    wo.WONumber,
		Title:       wo.Title,
		Description: wo.Description,
		Status:      string(wo.Status),
		CreatedAt:   wo.CreatedAt.Format(time.RFC3339),
	}
	if includeLines {
		for _, line := range wo.Spareparts {
			itemName, sku := "", ""
			if line.MasterItem.Id != "" {
				itemName = line.MasterItem.Name
				sku = line.MasterItem.SKU
			}
			res.Spareparts = append(res.Spareparts, dto.WorkOrderSparepartResponse{
				ItemId:   line.ItemId,
				ItemName: itemName,
				SKU:      sku,
				Quantity: line.Quantity,
			})
		}
	}
	return res
}

func ToWorkOrderResponses(orders []model.WorkOrder) []dto.WorkOrderResponse {
	responses := make([]dto.WorkOrderResponse, 0, len(orders))
	for i := range orders {
		responses = append(responses, *ToWorkOrderResponse(&orders[i], false))
	}
	return responses
}
