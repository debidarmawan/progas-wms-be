package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"time"
)

type HydrotestUsecase interface {
	FindDue(withinDays int) (*dto.HydrotestDueResponse, global.ErrorResponse)
	Record(actorUserId, cylinderId string, req *dto.RecordHydrotestRequest) global.ErrorResponse
}

type hydrotestUsecase struct {
	txManager    helper.TxManager
	cylinderRepo repository.CylinderRepository
	ledgerRepo   repository.CylinderLedgerRepository
	auditLogRepo repository.AuditLogRepository
}

func NewHydrotestUsecase(
	txManager helper.TxManager,
	cylinderRepo repository.CylinderRepository,
	ledgerRepo repository.CylinderLedgerRepository,
	auditLogRepo repository.AuditLogRepository,
) HydrotestUsecase {
	return &hydrotestUsecase{
		txManager:    txManager,
		cylinderRepo: cylinderRepo,
		ledgerRepo:   ledgerRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *hydrotestUsecase) FindDue(withinDays int) (*dto.HydrotestDueResponse, global.ErrorResponse) {
	if withinDays <= 0 {
		withinDays = 30
	}
	cylinders, err := u.cylinderRepo.FindHydrotestDue(withinDays)
	if err != nil {
		return nil, err
	}

	items := make([]dto.HydrotestDueCylinder, 0, len(cylinders))
	for _, cyl := range cylinders {
		expiry := helper.HydrotestExpiryDate(cyl.LastHydrotestDate)
		items = append(items, dto.HydrotestDueCylinder{
			Id:                cyl.Id,
			BarcodeSN:         cyl.BarcodeSN,
			Status:            string(cyl.Status),
			LastHydrotestDate: cyl.LastHydrotestDate.Format(time.RFC3339),
			ExpiryDate:        expiry.Format(time.RFC3339),
			IsExpired:         helper.IsHydrotestExpired(cyl.LastHydrotestDate),
		})
	}

	return &dto.HydrotestDueResponse{
		DueWithinDays: withinDays,
		Items:         items,
	}, nil
}

func (u *hydrotestUsecase) Record(actorUserId, cylinderId string, req *dto.RecordHydrotestRequest) global.ErrorResponse {
	hydrotestDate, parseErr := time.Parse(time.RFC3339, req.LastHydrotestDate)
	if parseErr != nil {
		hydrotestDate, parseErr = time.Parse("2006-01-02", req.LastHydrotestDate)
	}
	if parseErr != nil {
		return global.BadRequestError("invalid last_hydrotest_date format (use YYYY-MM-DD or RFC3339)")
	}
	if !helper.ValidateHydrotestDate(hydrotestDate) {
		return global.BadRequestError("last hydrotest date is invalid or expired")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	cylinder, err := u.cylinderRepo.FindByIdForUpdate(tx, cylinderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	prevStatus := cylinder.Status
	cylinder.LastHydrotestDate = hydrotestDate

	ledgerCyl := *cylinder
	ledgerCyl.Status = prevStatus
	repository.LogCylinderStatusChanges(u.ledgerRepo, tx, []model.Cylinder{ledgerCyl}, enum.CylinderStatusMaintenance, constant.LedgerActionHydrotest, constant.AuditObjectCylinder, cylinder.Id)

	cylinder.Status = enum.CylinderStatusMaintenance
	if err := u.cylinderRepo.Update(tx, cylinder); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditHydrotestRecord, constant.AuditObjectCylinder, cylinder.Id, map[string]any{
		"barcode_sn":          cylinder.BarcodeSN,
		"last_hydrotest_date": hydrotestDate.Format(time.RFC3339),
		"notes":               req.Notes,
	})

	return nil
}
