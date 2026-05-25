package usecase

import (
	"fmt"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/repository"
)

type ExchangeUsecase interface {
	Process(actorUserId string, req *dto.ProcessExchangeRequest, canApprove bool) (*dto.ExchangeResponse, global.ErrorResponse)
}

type exchangeUsecase struct {
	txManager    helper.TxManager
	cylinderRepo repository.CylinderRepository
	ledgerRepo   repository.CylinderLedgerRepository
	customerRepo repository.CustomerRepository
	auditLogRepo repository.AuditLogRepository
}

func NewExchangeUsecase(
	txManager helper.TxManager,
	cylinderRepo repository.CylinderRepository,
	ledgerRepo repository.CylinderLedgerRepository,
	customerRepo repository.CustomerRepository,
	auditLogRepo repository.AuditLogRepository,
) ExchangeUsecase {
	return &exchangeUsecase{
		txManager:    txManager,
		cylinderRepo: cylinderRepo,
		ledgerRepo:   ledgerRepo,
		customerRepo: customerRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *exchangeUsecase) Process(actorUserId string, req *dto.ProcessExchangeRequest, canApprove bool) (*dto.ExchangeResponse, global.ErrorResponse) {
	if err := helper.ValidateBarcodeListsUnique(req.OutBarcodes, req.InBarcodes); err != nil {
		return nil, global.BadRequestError(err.Error())
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	customer, err := u.customerRepo.FindByIdForUpdate(tx, req.CustomerId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if !customer.IsActive {
		tx.Rollback()
		return nil, global.BadRequestError("customer is not active")
	}

	outCylinders, err := u.cylinderRepo.FindByBarcodes(tx, req.OutBarcodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if validationErr := helper.ValidateExchangeOutCylinders(outCylinders); validationErr != nil {
		tx.Rollback()
		return nil, global.BadRequestError(validationErr.Error())
	}

	inCylinders, err := u.cylinderRepo.FindByBarcodes(tx, req.InBarcodes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if validationErr := helper.ValidateExchangeInCylinders(inCylinders); validationErr != nil {
		tx.Rollback()
		return nil, global.BadRequestError(validationErr.Error())
	}

	outDelta := helper.CountOutstandingDelta(outCylinders)
	inDelta := helper.CountOutstandingDelta(inCylinders)
	netDelta := outDelta - inDelta

	if helper.WouldExceedQuota(customer.OutstandingCount, customer.CylinderQuotaLimit, netDelta) {
		if !req.ForceApprove {
			tx.Rollback()
			return nil, global.BadRequestError(fmt.Sprintf(
				"exchange would exceed cylinder quota (current: %d, delta: %+d, limit: %d). Set force_approve=true with approval permission",
				customer.OutstandingCount, netDelta, customer.CylinderQuotaLimit,
			))
		}
		if !canApprove {
			tx.Rollback()
			return nil, global.ForbiddenError()
		}
	}

	outIds := make([]string, 0, len(outCylinders))
	for _, cyl := range outCylinders {
		outIds = append(outIds, cyl.Id)
	}
	inIds := make([]string, 0, len(inCylinders))
	for _, cyl := range inCylinders {
		inIds = append(inIds, cyl.Id)
	}

	repository.LogCylinderStatusChanges(u.ledgerRepo, tx, outCylinders, enum.CylinderStatusOutstanding, constant.LedgerActionExchangeOut, constant.AuditObjectCustomer, customer.Id)
	if err := u.cylinderRepo.UpdateStatusByIds(tx, outIds, enum.CylinderStatusOutstanding); err != nil {
		tx.Rollback()
		return nil, err
	}
	repository.LogCylinderStatusChanges(u.ledgerRepo, tx, inCylinders, enum.CylinderStatusEmpty, constant.LedgerActionExchangeIn, constant.AuditObjectCustomer, customer.Id)
	if err := u.cylinderRepo.UpdateStatusByIds(tx, inIds, enum.CylinderStatusEmpty); err != nil {
		tx.Rollback()
		return nil, err
	}

	if netDelta != 0 {
		if err := u.customerRepo.AdjustOutstanding(tx, customer.Id, netDelta); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	alerts := append(
		helper.DetectCrossCustomerAlerts(customer.Id, outCylinders),
		helper.DetectCrossCustomerAlerts(customer.Id, inCylinders)...,
	)

	outstandingBefore := customer.OutstandingCount
	outstandingAfter := outstandingBefore + netDelta

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditExchangeComplete, constant.AuditObjectCustomer, customer.Id, map[string]any{
		"customer_id":        customer.Id,
		"out_count":          len(outCylinders),
		"in_count":           len(inCylinders),
		"outstanding_before": outstandingBefore,
		"outstanding_after":  outstandingAfter,
		"outstanding_delta":  netDelta,
		"force_approve":      req.ForceApprove,
		"cross_customer_alerts": alerts,
	})

	return &dto.ExchangeResponse{
		CustomerId:          customer.Id,
		OutCount:            len(outCylinders),
		InCount:             len(inCylinders),
		OutstandingBefore:   outstandingBefore,
		OutstandingAfter:    outstandingAfter,
		OutstandingDelta:    netDelta,
		OutBarcodes:         req.OutBarcodes,
		InBarcodes:          req.InBarcodes,
		CrossCustomerAlerts: alerts,
	}, nil
}
