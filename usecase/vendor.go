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

type VendorUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.VendorResponse], global.ErrorResponse)
	FindById(id string) (*dto.VendorDetailResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateVendorRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateVendorRequest) global.ErrorResponse
	Delete(actorUserId, id string) global.ErrorResponse
}

type vendorUsecase struct {
	txManager    helper.TxManager
	vendorRepo   repository.VendorRepository
	cylinderRepo repository.CylinderRepository
	auditLogRepo repository.AuditLogRepository
}

func NewVendorUsecase(
	txManager helper.TxManager,
	vendorRepo repository.VendorRepository,
	cylinderRepo repository.CylinderRepository,
	auditLogRepo repository.AuditLogRepository,
) VendorUsecase {
	return &vendorUsecase{
		txManager:    txManager,
		vendorRepo:   vendorRepo,
		cylinderRepo: cylinderRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *vendorUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.VendorResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	vendors, total, err := u.vendorRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(vendors))
	for _, v := range vendors {
		ids = append(ids, v.Id)
	}
	counts, err := u.cylinderRepo.CountByVendorIds(ids)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedResponse[dto.VendorResponse]{
		Items: mapper.ToVendorResponses(vendors, counts),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *vendorUsecase) FindById(id string) (*dto.VendorDetailResponse, global.ErrorResponse) {
	vendor, err := u.vendorRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	count, err := u.cylinderRepo.CountByVendorId(id)
	if err != nil {
		return nil, err
	}

	cylinders, err := u.cylinderRepo.FindByVendorId(id)
	if err != nil {
		return nil, err
	}

	base := mapper.ToVendorResponse(vendor, int(count))
	return &dto.VendorDetailResponse{
		VendorResponse:    *base,
		CylindersByStatus: mapper.CylindersByStatus(cylinders),
		Cylinders:         mapper.ToVendorCylinderSummaries(cylinders),
	}, nil
}

func (u *vendorUsecase) Create(actorUserId string, req *dto.CreateVendorRequest) global.ErrorResponse {
	existing, err := u.vendorRepo.FindByCode(req.Code)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existing != nil {
		return global.BadRequestError("vendor code already exists")
	}

	startDate, parseErr := mapper.ParseOptionalDate(req.ContractStartDate)
	if parseErr != nil {
		return global.BadRequestError("invalid contract_start_date format (use YYYY-MM-DD)")
	}
	endDate, parseErr := mapper.ParseOptionalDate(req.ContractEndDate)
	if parseErr != nil {
		return global.BadRequestError("invalid contract_end_date format (use YYYY-MM-DD)")
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return global.BadRequestError("contract_end_date must be on or after contract_start_date")
	}

	vendor := &model.Vendor{
		Code:              req.Code,
		Name:              req.Name,
		ContactPerson:     req.ContactPerson,
		Phone:             req.Phone,
		Email:             req.Email,
		Address:           req.Address,
		Notes:             req.Notes,
		ContractStartDate: startDate,
		ContractEndDate:   endDate,
		IsActive:          true,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.vendorRepo.Create(tx, vendor); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditVendorCreate, constant.AuditObjectVendor, vendor.Id, map[string]any{
		"code": vendor.Code,
		"name": vendor.Name,
	})

	return nil
}

func (u *vendorUsecase) Update(actorUserId, id string, req *dto.UpdateVendorRequest) global.ErrorResponse {
	vendor, err := u.vendorRepo.FindById(id)
	if err != nil {
		return err
	}

	startDate, parseErr := mapper.ParseOptionalDate(req.ContractStartDate)
	if parseErr != nil {
		return global.BadRequestError("invalid contract_start_date format (use YYYY-MM-DD)")
	}
	endDate, parseErr := mapper.ParseOptionalDate(req.ContractEndDate)
	if parseErr != nil {
		return global.BadRequestError("invalid contract_end_date format (use YYYY-MM-DD)")
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return global.BadRequestError("contract_end_date must be on or after contract_start_date")
	}

	vendor.Name = req.Name
	vendor.ContactPerson = req.ContactPerson
	vendor.Phone = req.Phone
	vendor.Email = req.Email
	vendor.Address = req.Address
	vendor.Notes = req.Notes
	vendor.ContractStartDate = startDate
	vendor.ContractEndDate = endDate
	vendor.IsActive = req.IsActive

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.vendorRepo.Update(tx, vendor); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditVendorUpdate, constant.AuditObjectVendor, vendor.Id, map[string]any{
		"code":      vendor.Code,
		"is_active": vendor.IsActive,
	})

	return nil
}

func (u *vendorUsecase) Delete(actorUserId, id string) global.ErrorResponse {
	if _, err := u.vendorRepo.FindById(id); err != nil {
		return err
	}

	count, err := u.cylinderRepo.CountByVendorId(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return global.BadRequestError("cannot delete vendor with registered rental cylinders")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.vendorRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditVendorDelete, constant.AuditObjectVendor, id, nil)

	return nil
}
