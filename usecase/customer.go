package usecase

import (
	"fmt"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
	"strings"

	"github.com/google/uuid"
)

type CustomerUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerResponse], global.ErrorResponse)
	FindById(id string) (*dto.CustomerResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateCustomerRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateCustomerRequest) global.ErrorResponse
}

type customerUsecase struct {
	txManager         helper.TxManager
	customerRepo      repository.CustomerRepository
	termOfPaymentRepo repository.TermOfPaymentRepository
	auditLogRepo      repository.AuditLogRepository
}

func NewCustomerUsecase(
	txManager helper.TxManager,
	customerRepo repository.CustomerRepository,
	termOfPaymentRepo repository.TermOfPaymentRepository,
	auditLogRepo repository.AuditLogRepository,
) CustomerUsecase {
	return &customerUsecase{
		txManager:         txManager,
		customerRepo:      customerRepo,
		termOfPaymentRepo: termOfPaymentRepo,
		auditLogRepo:      auditLogRepo,
	}
}

// resolveTermOfPayment maps the optional TOP id to its FK; an empty id means "no TOP".
// Invoice due date still uses Customer.PaymentTermDays, so nothing is derived here.
func (u *customerUsecase) resolveTermOfPayment(id string) (*string, global.ErrorResponse) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return nil, nil
	}
	top, err := u.termOfPaymentRepo.FindById(trimmed)
	if err != nil {
		return nil, err
	}
	if !top.IsActive {
		return nil, global.BadRequestError("Term of payment is inactive")
	}
	return &top.Id, nil
}

func (u *customerUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	customers, total, err := u.customerRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.CustomerResponse]{
		Items: mapper.ToCustomerResponses(customers),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *customerUsecase) FindById(id string) (*dto.CustomerResponse, global.ErrorResponse) {
	customer, err := u.customerRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToCustomerResponse(customer), nil
}

func generateCustomerCode() string {
	return fmt.Sprintf("CUST-%s", strings.ToUpper(uuid.New().String()[:8]))
}

func (u *customerUsecase) Create(actorUserId string, req *dto.CreateCustomerRequest) global.ErrorResponse {
	termOfPaymentId, errResp := u.resolveTermOfPayment(req.TermOfPaymentId)
	if errResp != nil {
		return errResp
	}

	customer := &model.Customer{
		Code:               generateCustomerCode(),
		Name:               req.Name,
		Pic:                req.Pic,
		Npwp:               req.Npwp,
		Fax:                req.Fax,
		Email:              req.Email,
		Phone:              req.Phone,
		Address:            req.Address,
		PostalCode:         req.PostalCode,
		CylinderQuotaLimit: req.CylinderQuotaLimit,
		CreditLimit:        req.CreditLimit,
		TermOfPaymentId:    termOfPaymentId,
		SalesPicName:       req.SalesPicName,
		SalesPicPhone:      req.SalesPicPhone,
		SalesPicEmail:      req.SalesPicEmail,
		OutstandingCount:   0,
		IsActive:           true,
		PaymentTermDays:    30,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.customerRepo.Create(tx, customer); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerCreate, constant.AuditObjectCustomer, customer.Id, map[string]any{
		"code": customer.Code,
		"name": customer.Name,
	})

	return nil
}

func (u *customerUsecase) Update(actorUserId, id string, req *dto.UpdateCustomerRequest) global.ErrorResponse {
	customer, err := u.customerRepo.FindById(id)
	if err != nil {
		return err
	}

	termOfPaymentId, errResp := u.resolveTermOfPayment(req.TermOfPaymentId)
	if errResp != nil {
		return errResp
	}

	customer.Name = req.Name
	customer.Pic = req.Pic
	customer.Npwp = req.Npwp
	customer.Fax = req.Fax
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.Address = req.Address
	customer.PostalCode = req.PostalCode
	customer.CylinderQuotaLimit = req.CylinderQuotaLimit
	customer.CreditLimit = req.CreditLimit
	customer.TermOfPaymentId = termOfPaymentId
	customer.SalesPicName = req.SalesPicName
	customer.SalesPicPhone = req.SalesPicPhone
	customer.SalesPicEmail = req.SalesPicEmail
	customer.IsActive = req.IsActive

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.customerRepo.Update(tx, customer); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerUpdate, constant.AuditObjectCustomer, customer.Id, map[string]any{
		"code": customer.Code,
	})

	return nil
}
