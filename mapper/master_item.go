package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
)

func ToMasterItemResponse(item *model.MasterItem, stockQty *int) *dto.MasterItemResponse {
	res := &dto.MasterItemResponse{
		Id:            item.Id,
		Name:          item.Name,
		SKU:           item.SKU,
		GasType:       item.GasType,
		IsSerialized:  item.IsSerialized,
		EmptyWeightKg: item.EmptyWeightKg,
		GasWeightKg:   item.GasWeightKg,
		MinStockAlert: item.MinStockAlert,
		StockQuantity: stockQty,
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
