package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToCustomerPOResponse(po *model.CustomerPO, includeLines bool) *dto.CustomerPOResponse {
	res := &dto.CustomerPOResponse{
		Id:           po.Id,
		PONumber:     po.PONumber,
		CustomerId:   po.CustomerId,
		CustomerName: po.Customer.Name,
		PODate:       po.PODate.Format("2006-01-02"),
		Status:       string(po.Status),
		CreatedAt:    po.CreatedAt.Format(time.RFC3339),
		DocumentURL:  po.DocumentURL,
	}
	if po.ValidUntil != nil {
		res.ValidUntil = po.ValidUntil.Format("2006-01-02")
	}
	if includeLines {
		res.Lines = make([]dto.CustomerPOLineResponse, 0, len(po.Lines))
		for _, line := range po.Lines {
			res.Lines = append(res.Lines, dto.CustomerPOLineResponse{
				Id:           line.Id,
				MasterItemId: line.MasterItemId,
				SKU:          line.MasterItem.SKU,
				ItemName:     line.MasterItem.Name,
				Quantity:     line.Quantity,
			})
		}
	}
	return res
}

func ToCustomerPOResponses(pos []model.CustomerPO) []dto.CustomerPOResponse {
	res := make([]dto.CustomerPOResponse, 0, len(pos))
	for i := range pos {
		res = append(res, *ToCustomerPOResponse(&pos[i], false))
	}
	return res
}
