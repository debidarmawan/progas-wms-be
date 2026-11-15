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

	"github.com/gofiber/fiber/v3"
)

const maxTermOfPaymentDays = 365

type TermOfPaymentUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.TermOfPaymentResponse], global.ErrorResponse)
	FindById(id string) (*dto.TermOfPaymentResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateTermOfPaymentRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateTermOfPaymentRequest) global.ErrorResponse
	Delete(actorUserId, id string) global.ErrorResponse
}

type termOfPaymentUsecase struct {
	txManager    helper.TxManager
	topRepo      repository.TermOfPaymentRepository
	auditLogRepo repository.AuditLogRepository
}

func NewTermOfPaymentUsecase(
	txManager helper.TxManager,
	topRepo repository.TermOfPaymentRepository,
	auditLogRepo repository.AuditLogRepository,
) TermOfPaymentUsecase {
	return &termOfPaymentUsecase{
		txManager:    txManager,
		topRepo:      topRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *termOfPaymentUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.TermOfPaymentResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	tops, total, err := u.topRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.TermOfPaymentResponse]{
		Items: mapper.ToTermOfPaymentResponses(tops),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *termOfPaymentUsecase) FindById(id string) (*dto.TermOfPaymentResponse, global.ErrorResponse) {
	top, err := u.topRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToTermOfPaymentResponse(top), nil
}

func (u *termOfPaymentUsecase) Create(actorUserId string, req *dto.CreateTermOfPaymentRequest) global.ErrorResponse {
	code, name, err := validateTermOfPayment(req.Code, req.Name, req.Days)
	if err != nil {
		return err
	}

	if err := u.ensureCodeAvailable(code, ""); err != nil {
		return err
	}

	top := &model.TermOfPayment{
		Code:     code,
		Name:     name,
		Days:     req.Days,
		IsActive: req.IsActive == nil || *req.IsActive,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.topRepo.Create(tx, top); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditTermOfPaymentCreate, constant.AuditObjectTermOfPayment, top.Id, map[string]any{
		"code": top.Code,
		"days": top.Days,
	})

	return nil
}

func (u *termOfPaymentUsecase) Update(actorUserId, id string, req *dto.UpdateTermOfPaymentRequest) global.ErrorResponse {
	top, err := u.topRepo.FindById(id)
	if err != nil {
		return err
	}

	code, name, err := validateTermOfPayment(req.Code, req.Name, req.Days)
	if err != nil {
		return err
	}

	if err := u.ensureCodeAvailable(code, top.Id); err != nil {
		return err
	}

	top.Code = code
	top.Name = name
	top.Days = req.Days
	if req.IsActive != nil {
		top.IsActive = *req.IsActive
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.topRepo.Update(tx, top); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditTermOfPaymentUpdate, constant.AuditObjectTermOfPayment, top.Id, map[string]any{
		"code":      top.Code,
		"days":      top.Days,
		"is_active": top.IsActive,
	})

	return nil
}

func (u *termOfPaymentUsecase) Delete(actorUserId, id string) global.ErrorResponse {
	if _, err := u.topRepo.FindById(id); err != nil {
		return err
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.topRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditTermOfPaymentDelete, constant.AuditObjectTermOfPayment, id, nil)

	return nil
}

func (u *termOfPaymentUsecase) ensureCodeAvailable(code, currentId string) global.ErrorResponse {
	existing, err := u.topRepo.FindByCode(code)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return nil
		}
		return err
	}
	if existing.Id != currentId {
		return global.BadRequestError("code already exists")
	}
	return nil
}

func validateTermOfPayment(rawCode, rawName string, days int) (string, string, global.ErrorResponse) {
	code := normalizeTermOfPaymentCode(rawCode)
	name := strings.TrimSpace(rawName)
	if code == "" {
		return "", "", global.BadRequestError("code is required")
	}
	if name == "" {
		return "", "", global.BadRequestError("name is required")
	}
	if days < 0 || days > maxTermOfPaymentDays {
		return "", "", global.BadRequestError("days must be between 0 and 365")
	}
	return code, name, nil
}

func normalizeTermOfPaymentCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
