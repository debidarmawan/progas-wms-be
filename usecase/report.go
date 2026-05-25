package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/repository"
	"time"
)

type ReportUsecase interface {
	StockLedger(barcode string) (*dto.StockLedgerReportResponse, global.ErrorResponse)
	Turnaround(from, to string) (*dto.TurnaroundReportResponse, global.ErrorResponse)
}

type reportUsecase struct {
	ledgerRepo   repository.CylinderLedgerRepository
	cylinderRepo repository.CylinderRepository
}

func NewReportUsecase(
	ledgerRepo repository.CylinderLedgerRepository,
	cylinderRepo repository.CylinderRepository,
) ReportUsecase {
	return &reportUsecase{
		ledgerRepo:   ledgerRepo,
		cylinderRepo: cylinderRepo,
	}
}

func (u *reportUsecase) StockLedger(barcode string) (*dto.StockLedgerReportResponse, global.ErrorResponse) {
	if barcode == "" {
		return nil, global.BadRequestError("barcode is required")
	}
	if _, err := u.cylinderRepo.FindByBarcode(barcode); err != nil {
		return nil, err
	}

	entries, err := u.ledgerRepo.FindByBarcode(barcode)
	if err != nil {
		return nil, err
	}

	return &dto.StockLedgerReportResponse{
		BarcodeSN: barcode,
		Entries:   mapper.ToCylinderLedgerEntries(entries),
	}, nil
}

func (u *reportUsecase) Turnaround(from, to string) (*dto.TurnaroundReportResponse, global.ErrorResponse) {
	fromTime, err := parseReportDate(from, time.Now().AddDate(0, -1, 0))
	if err != nil {
		return nil, err
	}
	toTime, err := parseReportDate(to, time.Now())
	if err != nil {
		return nil, err
	}
	if toTime.Before(fromTime) {
		return nil, global.BadRequestError("to date must be after from date")
	}

	entries, err := u.ledgerRepo.FindByDateRange(fromTime, toTime)
	if err != nil {
		return nil, err
	}

	samples := helper.ComputeTurnaroundFromLedger(entries)
	dtoSamples := mapper.ToTurnaroundSamples(samples)

	return &dto.TurnaroundReportResponse{
		FromDate:    fromTime.Format("2006-01-02"),
		ToDate:      toTime.Format("2006-01-02"),
		SampleCount: len(dtoSamples),
		AverageDays: helper.AverageTurnaroundDays(samples),
		Samples:     dtoSamples,
	}, nil
}

func parseReportDate(value string, fallback time.Time) (time.Time, global.ErrorResponse) {
	if value == "" {
		return fallback, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, global.BadRequestError("invalid date format, use YYYY-MM-DD")
	}
	return t, nil
}
