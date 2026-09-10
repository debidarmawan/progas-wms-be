package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToCustomerItemPriceResponse(price *model.CustomerItemPrice, item *model.MasterItem, at time.Time) *dto.CustomerItemPriceResponse {
	response := &dto.CustomerItemPriceResponse{
		Id:            price.Id,
		CustomerId:    price.CustomerId,
		MasterItemId:  price.MasterItemId,
		Price:         price.Price,
		EffectiveFrom: price.EffectiveFrom.Format(time.RFC3339),
		IsActive:      !price.EffectiveFrom.After(at) && (price.EffectiveTo == nil || price.EffectiveTo.After(at)),
	}
	if item != nil {
		response.ItemName = item.Name
		response.ItemSKU = item.SKU
	}
	if price.EffectiveTo != nil {
		formatted := price.EffectiveTo.Format(time.RFC3339)
		response.EffectiveTo = &formatted
	}
	return response
}
