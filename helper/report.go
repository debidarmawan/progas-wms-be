package helper

import (
	"progas-wms-be/enum"
	"progas-wms-be/model"
	"time"
)

type TurnaroundSample struct {
	BarcodeSN  string  `json:"barcode_sn"`
	Days       float64 `json:"days"`
	StartedAt  string  `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

func ComputeTurnaroundFromLedger(entries []model.CylinderLedger) []TurnaroundSample {
	byBarcode := make(map[string][]model.CylinderLedger)
	for _, e := range entries {
		byBarcode[e.BarcodeSN] = append(byBarcode[e.BarcodeSN], e)
	}

	var samples []TurnaroundSample
	for barcode, logs := range byBarcode {
		var emptyAt, readyAt *time.Time
		for _, log := range logs {
			t := log.CreatedAt
			if log.ToStatus == enum.CylinderStatusEmpty && emptyAt == nil {
				emptyAt = &t
			}
			if log.ToStatus == enum.CylinderStatusReady {
				readyAt = &t
			}
		}
		if emptyAt != nil && readyAt != nil && readyAt.After(*emptyAt) {
			days := readyAt.Sub(*emptyAt).Hours() / 24
			samples = append(samples, TurnaroundSample{
				BarcodeSN:   barcode,
				Days:        days,
				StartedAt:   emptyAt.Format(time.RFC3339),
				CompletedAt: readyAt.Format(time.RFC3339),
			})
		}
	}
	return samples
}

func AverageTurnaroundDays(samples []TurnaroundSample) float64 {
	if len(samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		sum += s.Days
	}
	return sum / float64(len(samples))
}
