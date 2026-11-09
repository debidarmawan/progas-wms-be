package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToMasterItemResponse(item *model.MasterItem, stockQty *int) *dto.MasterItemResponse {
	res := &dto.MasterItemResponse{
		Id:                item.Id,
		Name:              item.Name,
		SKU:               item.SKU,
		ItemType:          item.ItemType,
		GasType:           item.GasType,
		HnaPrice:          item.HnaPrice,
		IsSerialized:      item.IsSerialized,
		EmptyWeightKg:     item.EmptyWeightKg,
		GasWeightKg:       item.GasWeightKg,
		MinStockAlert:     item.MinStockAlert,
		MaxDaysAtCustomer: item.MaxDaysAtCustomer,
		UomId:             item.UomId,
		StockQuantity:     stockQty,
	}
	if item.Uom != nil {
		res.UomCode = item.Uom.Code
		res.UomName = item.Uom.Name
	}
	return res
}

func ToMasterItemResponses(items []model.MasterItem) []dto.MasterItemResponse {
	responses := make([]dto.MasterItemResponse, 0, len(items))
	for i := range items {
		responses = append(responses, *ToMasterItemResponse(&items[i], nil))
	}
	return responses
}
