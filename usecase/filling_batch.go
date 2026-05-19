package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
)

type FillingBatchUsecase interface {
	Submit(actorUserId string, req *dto.SubmitFillingBatchRequest) (*dto.FillingBatchResponse, global.ErrorResponse)
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.FillingBatchResponse], global.ErrorResponse)
	FindById(id string) (*dto.FillingBatchResponse, global.ErrorResponse)
}

type fillingBatchUsecase struct {
	txManager       helper.TxManager
	fillingBatchRepo repository.FillingBatchRepository
	cylinderRepo    repository.CylinderRepository
	masterItemRepo  repository.MasterItemRepository
	auditLogRepo    repository.AuditLogRepository
}

func NewFillingBatchUsecase(
	txManager helper.TxManager,
	fillingBatchRepo repository.FillingBatchRepository,
	cylinderRepo repository.CylinderRepository,
	masterItemRepo repository.MasterItemRepository,
	auditLogRepo repository.AuditLogRepository,
) FillingBatchUsecase {
	return &fillingBatchUsecase{
		txManager:        txManager,
		fillingBatchRepo: fillingBatchRepo,
		cylinderRepo:     cylinderRepo,
		masterItemRepo:   masterItemRepo,
		auditLogRepo:     auditLogRepo,
	}
}

func (u *fillingBatchUsecase) Submit(actorUserId string, req *dto.SubmitFillingBatchRequest) (*dto.FillingBatchResponse, global.ErrorResponse) {
	if err := helper.ValidateBarcodeList(req.Barcodes); err != nil {
		return nil, global.BadRequestError(err.Error())
	}

	item, err := u.masterItemRepo.FindById(req.ItemId)
	if err != nil {
		return nil, err
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	cylinders, err := u.cylinderRepo.FindByBarcodes(tx, req.Barcodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if validationErr := helper.ValidateFillingBatchCylinders(cylinders, item); validationErr != nil {
		tx.Rollback()
		return nil, global.BadRequestError(validationErr.Error())
	}

	batch := &model.FillingBatch{
		BaseModel:   model.BaseModel{CreatedBy: actorUserId},
		BatchNumber: helper.GenerateFillingBatchNumber(),
		ItemId:      item.Id,
		GasType:     item.GasType,
		Status:      enum.FillingBatchStatusCompleted,
		CylinderQty: len(cylinders),
		Notes:       req.Notes,
	}

	if err := u.fillingBatchRepo.Create(tx, batch); err != nil {
		tx.Rollback()
		return nil, err
	}

	details := make([]model.FillingBatchDetail, 0, len(cylinders))
	cylinderIds := make([]string, 0, len(cylinders))
	for _, cyl := range cylinders {
		details = append(details, model.FillingBatchDetail{
			FillingBatchId: batch.Id,
			CylinderId:     cyl.Id,
			BarcodeSN:      cyl.BarcodeSN,
		})
		cylinderIds = append(cylinderIds, cyl.Id)
	}

	if err := u.fillingBatchRepo.CreateDetails(tx, details); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := u.cylinderRepo.UpdateStatusByIds(tx, cylinderIds, enum.CylinderStatusReady); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditFillingBatchSubmit, constant.AuditObjectFillingBatch, batch.Id, map[string]any{
		"batch_number": batch.BatchNumber,
		"item_id":      batch.ItemId,
		"gas_type":     batch.GasType,
		"cylinder_qty": batch.CylinderQty,
	})

	batch.Details = details
	return mapper.ToFillingBatchResponse(batch, true), nil
}

func (u *fillingBatchUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.FillingBatchResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	batches, total, err := u.fillingBatchRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.FillingBatchResponse]{
		Items: mapper.ToFillingBatchResponses(batches),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *fillingBatchUsecase) FindById(id string) (*dto.FillingBatchResponse, global.ErrorResponse) {
	batch, err := u.fillingBatchRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToFillingBatchResponse(batch, true), nil
}
