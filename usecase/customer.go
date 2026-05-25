package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type CustomerUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.CustomerResponse], global.ErrorResponse)
	FindById(id string) (*dto.CustomerResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateCustomerRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateCustomerRequest) global.ErrorResponse
}

type customerUsecase struct {
	txManager    helper.TxManager
	customerRepo repository.CustomerRepository
	auditLogRepo repository.AuditLogRepository
}

func NewCustomerUsecase(
	txManager helper.TxManager,
	customerRepo repository.CustomerRepository,
	auditLogRepo repository.AuditLogRepository,
) CustomerUsecase {
	return &customerUsecase{
		txManager:    txManager,
		customerRepo: customerRepo,
		auditLogRepo: auditLogRepo,
	}
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

func (u *customerUsecase) Create(actorUserId string, req *dto.CreateCustomerRequest) global.ErrorResponse {
	existing, err := u.customerRepo.FindByCode(req.Code)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existing != nil {
		return global.BadRequestError("customer code already exists")
	}

	customer := &model.Customer{
		Code:               req.Code,
		Name:               req.Name,
		Phone:              req.Phone,
		Address:            req.Address,
		CylinderQuotaLimit: req.CylinderQuotaLimit,
		OutstandingCount:   0,
		IsActive:           true,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.customerRepo.Create(tx, customer); err != nil {
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

	customer.Name = req.Name
	customer.Phone = req.Phone
	customer.Address = req.Address
	customer.CylinderQuotaLimit = req.CylinderQuotaLimit
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
