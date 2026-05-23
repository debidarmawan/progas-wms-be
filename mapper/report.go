package mapper

import (
	"progas-wms-be/dto"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"time"
)

func ToCylinderLedgerEntries(entries []model.CylinderLedger) []dto.CylinderLedgerEntryResponse {
	responses := make([]dto.CylinderLedgerEntryResponse, 0, len(entries))
	for _, e := range entries {
		responses = append(responses, dto.CylinderLedgerEntryResponse{
			Id:            e.Id,
			BarcodeSN:     e.BarcodeSN,
			FromStatus:    string(e.FromStatus),
			ToStatus:      string(e.ToStatus),
			Action:        e.Action,
			ReferenceType: e.ReferenceType,
			ReferenceId:   e.ReferenceId,
			CreatedAt:     e.CreatedAt.Format(time.RFC3339),
		})
	}
	return responses
}

func ToTurnaroundSamples(samples []helper.TurnaroundSample) []dto.TurnaroundSample {
	result := make([]dto.TurnaroundSample, 0, len(samples))
	for _, s := range samples {
		result = append(result, dto.TurnaroundSample{
			BarcodeSN:   s.BarcodeSN,
			Days:        s.Days,
			StartedAt:   s.StartedAt,
			CompletedAt: s.CompletedAt,
		})
	}
	return result
}
