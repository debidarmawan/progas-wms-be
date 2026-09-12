package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToSalesOrderResponse(order *model.SalesOrder, includeLines bool) *dto.SalesOrderResponse {
	res := &dto.SalesOrderResponse{
		Id:           order.Id,
		SONumber:     order.SONumber,
		CustomerPOId: order.CustomerPOId,
		CustomerId:   order.CustomerId,
		CustomerName: order.Customer.Name,
		Status:       string(order.Status),
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
	}
	if includeLines {
		res.Lines = make([]dto.SalesOrderLineResponse, 0, len(order.Lines))
		for _, line := range order.Lines {
			res.Lines = append(res.Lines, dto.SalesOrderLineResponse{
				Id: line.Id, MasterItemId: line.MasterItemId, SKU: line.MasterItem.SKU,
				ItemName: line.MasterItem.Name, QtyOrdered: line.QtyOrdered, QtyDelivered: line.QtyDelivered,
			})
		}
	}
	return res
}

func ToSalesOrderResponses(orders []model.SalesOrder) []dto.SalesOrderResponse {
	res := make([]dto.SalesOrderResponse, 0, len(orders))
	for i := range orders {
		res = append(res, *ToSalesOrderResponse(&orders[i], false))
	}
	return res
}
