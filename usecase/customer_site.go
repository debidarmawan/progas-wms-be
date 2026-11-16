package usecase

import (
	"strings"

	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"
)

type CustomerSiteUsecase interface {
	FindByCustomer(customerId string) (*dto.CustomerSiteListResponse, global.ErrorResponse)
	Create(actorUserId, customerId string, req *dto.CreateCustomerSiteRequest) global.ErrorResponse
	Update(actorUserId, customerId, id string, req *dto.UpdateCustomerSiteRequest) global.ErrorResponse
	Delete(actorUserId, customerId, id string) global.ErrorResponse
}

type customerSiteUsecase struct {
	txManager        helper.TxManager
	customerRepo     repository.CustomerRepository
	customerSiteRepo repository.CustomerSiteRepository
	auditLogRepo     repository.AuditLogRepository
}

func NewCustomerSiteUsecase(
	txManager helper.TxManager,
	customerRepo repository.CustomerRepository,
	customerSiteRepo repository.CustomerSiteRepository,
	auditLogRepo repository.AuditLogRepository,
) CustomerSiteUsecase {
	return &customerSiteUsecase{
		txManager:        txManager,
		customerRepo:     customerRepo,
		customerSiteRepo: customerSiteRepo,
		auditLogRepo:     auditLogRepo,
	}
}

func (u *customerSiteUsecase) FindByCustomer(customerId string) (*dto.CustomerSiteListResponse, global.ErrorResponse) {
	if _, err := u.customerRepo.FindById(customerId); err != nil {
		return nil, err
	}
	sites, err := u.customerSiteRepo.FindByCustomer(customerId)
	if err != nil {
		return nil, err
	}
	return &dto.CustomerSiteListResponse{Items: mapper.ToCustomerSiteResponses(sites)}, nil
}

func (u *customerSiteUsecase) Create(actorUserId, customerId string, req *dto.CreateCustomerSiteRequest) global.ErrorResponse {
	if _, err := u.customerRepo.FindById(customerId); err != nil {
		return err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return global.BadRequestError("name is required")
	}

	existing, err := u.customerSiteRepo.FindByCustomer(customerId)
	if err != nil {
		return err
	}

	site := &model.CustomerSite{
		BaseModel:  model.BaseModel{CreatedBy: actorUserId},
		CustomerId: customerId,
		Name:       name,
		Pic:        strings.TrimSpace(req.Pic),
		Phone:      strings.TrimSpace(req.Phone),
		Address:    strings.TrimSpace(req.Address),
		PostalCode: strings.TrimSpace(req.PostalCode),
		Notes:      strings.TrimSpace(req.Notes),
		// The first site of a customer becomes primary so delivery orders have a default ship-to.
		IsPrimary: req.IsPrimary || len(existing) == 0,
		IsActive:  req.IsActive == nil || *req.IsActive,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if site.IsPrimary && len(existing) > 0 {
		if err := u.customerSiteRepo.ClearPrimary(tx, customerId, site.Id); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := u.customerSiteRepo.Create(tx, site); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerSiteCreate, constant.AuditObjectCustomerSite, site.Id, map[string]any{
		"customer_id": customerId,
		"name":        site.Name,
		"is_primary":  site.IsPrimary,
	})

	return nil
}

func (u *customerSiteUsecase) Update(actorUserId, customerId, id string, req *dto.UpdateCustomerSiteRequest) global.ErrorResponse {
	site, err := u.customerSiteRepo.FindById(id)
	if err != nil {
		return err
	}
	if site.CustomerId != customerId {
		return global.NotFoundError("Customer site not found")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return global.BadRequestError("name is required")
	}
	if !req.IsActive && (site.IsPrimary || req.IsPrimary) {
		return global.BadRequestError("primary site cannot be inactive")
	}

	site.Name = name
	site.Pic = strings.TrimSpace(req.Pic)
	site.Phone = strings.TrimSpace(req.Phone)
	site.Address = strings.TrimSpace(req.Address)
	site.PostalCode = strings.TrimSpace(req.PostalCode)
	site.Notes = strings.TrimSpace(req.Notes)
	site.IsPrimary = req.IsPrimary
	site.IsActive = req.IsActive
	site.UpdatedBy = actorUserId

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if site.IsPrimary {
		if err := u.customerSiteRepo.ClearPrimary(tx, customerId, site.Id); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := u.customerSiteRepo.Update(tx, site); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerSiteUpdate, constant.AuditObjectCustomerSite, site.Id, map[string]any{
		"customer_id": customerId,
		"name":        site.Name,
		"is_primary":  site.IsPrimary,
		"is_active":   site.IsActive,
	})

	return nil
}

func (u *customerSiteUsecase) Delete(actorUserId, customerId, id string) global.ErrorResponse {
	site, err := u.customerSiteRepo.FindById(id)
	if err != nil {
		return err
	}
	if site.CustomerId != customerId {
		return global.NotFoundError("Customer site not found")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.customerSiteRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	// Delivery orders keep their own ship-to snapshot, so deleting a site is safe.
	_ = u.auditLogRepo.Log(actorUserId, constant.AuditCustomerSiteDelete, constant.AuditObjectCustomerSite, id, map[string]any{
		"customer_id": customerId,
		"name":        site.Name,
	})

	return nil
}
