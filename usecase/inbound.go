package usecase

import (
	"fmt"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/repository"
	"strings"
)

type InboundUsecase interface {
	EmptyReceive(actorUserId string, req *dto.BarcodeListRequest) (*dto.BarcodeOperationResponse, global.ErrorResponse)
	PreFillQC(actorUserId string, req *dto.BarcodeListRequest) (*dto.BarcodeOperationResponse, global.ErrorResponse)
}

type inboundUsecase struct {
	txManager    helper.TxManager
	cylinderRepo repository.CylinderRepository
	auditLogRepo repository.AuditLogRepository
}

func NewInboundUsecase(
	txManager helper.TxManager,
	cylinderRepo repository.CylinderRepository,
	auditLogRepo repository.AuditLogRepository,
) InboundUsecase {
	return &inboundUsecase{
		txManager:    txManager,
		cylinderRepo: cylinderRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *inboundUsecase) EmptyReceive(actorUserId string, req *dto.BarcodeListRequest) (*dto.BarcodeOperationResponse, global.ErrorResponse) {
	if err := helper.ValidateBarcodeList(req.Barcodes); err != nil {
		return nil, global.BadRequestError(err.Error())
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	cylinders, err := u.cylinderRepo.FindByBarcodes(tx, req.Barcodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var invalid []string
	var ids []string
	var processed []string
	for _, cyl := range cylinders {
		if !helper.CanReceiveAsEmpty(cyl.Status) {
			invalid = append(invalid, fmt.Sprintf("%s (%s)", cyl.BarcodeSN, cyl.Status))
			continue
		}
		ids = append(ids, cyl.Id)
		processed = append(processed, cyl.BarcodeSN)
	}

	if len(invalid) > 0 {
		tx.Rollback()
		return nil, global.BadRequestError(fmt.Sprintf("invalid cylinder status for empty receive: %s", strings.Join(invalid, ", ")))
	}

	if err := u.cylinderRepo.UpdateStatusByIds(tx, ids, enum.CylinderStatusEmpty); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditEmptyReceive, constant.AuditObjectCylinder, "", map[string]any{
		"barcodes": processed,
		"count":    len(processed),
	})

	return &dto.BarcodeOperationResponse{
		ProcessedCount: len(processed),
		Barcodes:       processed,
	}, nil
}

func (u *inboundUsecase) PreFillQC(actorUserId string, req *dto.BarcodeListRequest) (*dto.BarcodeOperationResponse, global.ErrorResponse) {
	if err := helper.ValidateBarcodeList(req.Barcodes); err != nil {
		return nil, global.BadRequestError(err.Error())
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	cylinders, err := u.cylinderRepo.FindByBarcodes(tx, req.Barcodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var invalid []string
	var ids []string
	var processed []string
	for _, cyl := range cylinders {
		if !helper.CanPreFillQC(cyl.Status) {
			invalid = append(invalid, fmt.Sprintf("%s (%s)", cyl.BarcodeSN, cyl.Status))
			continue
		}
		ids = append(ids, cyl.Id)
		processed = append(processed, cyl.BarcodeSN)
	}

	if len(invalid) > 0 {
		tx.Rollback()
		return nil, global.BadRequestError(fmt.Sprintf("invalid cylinder status for pre-fill QC: %s", strings.Join(invalid, ", ")))
	}

	if err := u.cylinderRepo.UpdateStatusByIds(tx, ids, enum.CylinderStatusReadyToFill); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditPreFillQC, constant.AuditObjectCylinder, "", map[string]any{
		"barcodes": processed,
		"count":    len(processed),
	})

	return &dto.BarcodeOperationResponse{
		ProcessedCount: len(processed),
		Barcodes:       processed,
	}, nil
}
