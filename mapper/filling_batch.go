package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/model"
	"time"
)

func ToFillingBatchResponse(batch *model.FillingBatch, includeDetails bool) *dto.FillingBatchResponse {
	itemName := ""
	if batch.MasterItem.Id != "" {
		itemName = batch.MasterItem.Name
	}

	res := &dto.FillingBatchResponse{
		Id:          batch.Id,
		BatchNumber: batch.BatchNumber,
		ItemId:      batch.ItemId,
		ItemName:    itemName,
		GasType:     batch.GasType,
		Status:      string(batch.Status),
		CylinderQty: batch.CylinderQty,
		Notes:       batch.Notes,
		CreatedAt:   batch.CreatedAt.Format(time.RFC3339),
	}

	if includeDetails {
		for _, d := range batch.Details {
			res.Details = append(res.Details, dto.FillingBatchDetailResponse{
				Id:         d.Id,
				CylinderId: d.CylinderId,
				BarcodeSN:  d.BarcodeSN,
			})
		}
	}
	return res
}

func ToFillingBatchResponses(batches []model.FillingBatch) []dto.FillingBatchResponse {
	responses := make([]dto.FillingBatchResponse, 0, len(batches))
	for i := range batches {
		responses = append(responses, *ToFillingBatchResponse(&batches[i], false))
	}
	return responses
}
