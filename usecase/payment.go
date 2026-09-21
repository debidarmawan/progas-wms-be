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
	"time"
)

type PaymentUsecase interface {
	Record(actorUserId, invoiceId string, req *dto.RecordPaymentRequest) (*dto.InvoiceResponse, global.ErrorResponse)
}

type paymentUsecase struct {
	txManager    helper.TxManager
	paymentRepo  repository.PaymentRepository
	invoiceRepo  repository.InvoiceRepository
	customerRepo repository.CustomerRepository
	auditLogRepo repository.AuditLogRepository
}

func NewPaymentUsecase(
	txManager helper.TxManager,
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.InvoiceRepository,
	customerRepo repository.CustomerRepository,
	auditLogRepo repository.AuditLogRepository,
) PaymentUsecase {
	return &paymentUsecase{
		txManager:    txManager,
		paymentRepo:  paymentRepo,
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *paymentUsecase) Record(actorUserId, invoiceId string, req *dto.RecordPaymentRequest) (*dto.InvoiceResponse, global.ErrorResponse) {
	invoice, err := u.invoiceRepo.FindById(invoiceId)
	if err != nil {
		return nil, err
	}
	if invoice.Status == enum.InvoiceStatusPaid || invoice.Status == enum.InvoiceStatusCancelled {
		return nil, global.BadRequestError("invoice is already paid or cancelled")
	}

	newPaidAmount := invoice.PaidAmount + req.Amount
	if newPaidAmount > invoice.TotalAmount {
		return nil, global.BadRequestError("payment amount exceeds remaining invoice balance")
	}
	status := helper.DetermineInvoiceStatus(newPaidAmount, invoice.TotalAmount)

	tx := u.txManager.New()
	defer tx.CheckPanic()

	payment := &model.Payment{
		BaseModel:   model.BaseModel{CreatedBy: actorUserId},
		InvoiceId:   invoice.Id,
		Amount:      req.Amount,
		PaidAt:      time.Now(),
		Method:      req.Method,
		ReferenceNo: req.ReferenceNo,
	}
	if err := u.paymentRepo.Create(tx, payment); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := u.invoiceRepo.AdjustPaidAmount(tx, invoice.Id, req.Amount, string(status)); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := u.customerRepo.AdjustOutstandingBalance(tx, invoice.CustomerId, -req.Amount); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditPaymentRecord, constant.AuditObjectPayment, payment.Id, map[string]any{
		"invoice_id": invoice.Id,
		"amount":     req.Amount,
	})

	updated, err := u.invoiceRepo.FindById(invoice.Id)
	if err != nil {
		return nil, err
	}
	return mapper.ToInvoiceResponse(updated, true), nil
}
