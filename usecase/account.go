package usecase

import (
	"strings"

	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/mapper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type AccountUsecase interface {
	FindAll(query *dto.AccountListQuery) (*dto.PaginatedResponse[dto.AccountResponse], global.ErrorResponse)
	FindById(id string) (*dto.AccountResponse, global.ErrorResponse)
	Create(actorUserId string, req *dto.CreateAccountRequest) global.ErrorResponse
	Update(actorUserId, id string, req *dto.UpdateAccountRequest) global.ErrorResponse
	Delete(actorUserId, id string) global.ErrorResponse
}

type accountUsecase struct {
	txManager        helper.TxManager
	accountRepo      repository.AccountRepository
	journalEntryRepo repository.JournalEntryRepository
	auditLogRepo     repository.AuditLogRepository
}

func NewAccountUsecase(
	txManager helper.TxManager,
	accountRepo repository.AccountRepository,
	journalEntryRepo repository.JournalEntryRepository,
	auditLogRepo repository.AuditLogRepository,
) AccountUsecase {
	return &accountUsecase{
		txManager:        txManager,
		accountRepo:      accountRepo,
		journalEntryRepo: journalEntryRepo,
		auditLogRepo:     auditLogRepo,
	}
}

func (u *accountUsecase) FindAll(query *dto.AccountListQuery) (*dto.PaginatedResponse[dto.AccountResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(&query.ListQuery)
	search := helper.NormalizeSearch(query.Search)
	accountType := strings.ToUpper(strings.TrimSpace(query.Type))
	if accountType != "" && !enum.AccountType(accountType).IsValid() {
		return nil, global.BadRequestError("invalid account type")
	}

	accounts, total, err := u.accountRepo.FindAll(page, limit, search, accountType, query.IsActive)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.AccountResponse]{
		Items: mapper.ToAccountResponses(accounts),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *accountUsecase) FindById(id string) (*dto.AccountResponse, global.ErrorResponse) {
	account, err := u.accountRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToAccountResponse(account), nil
}

func (u *accountUsecase) Create(actorUserId string, req *dto.CreateAccountRequest) global.ErrorResponse {
	code := normalizeAccountCode(req.Code)
	name := strings.TrimSpace(req.Name)
	accountType, err := parseAccountType(req.Type)
	if err != nil {
		return err
	}
	if code == "" {
		return global.BadRequestError("code is required")
	}
	if name == "" {
		return global.BadRequestError("name is required")
	}
	if err := u.ensureCodeAvailable(code, ""); err != nil {
		return err
	}
	parent, err := u.resolveParent(req.ParentId, accountType, "")
	if err != nil {
		return err
	}

	account := &model.Account{
		Code:     code,
		Name:     name,
		Type:     accountType,
		IsActive: req.IsActive == nil || *req.IsActive,
	}
	if parent != nil {
		account.ParentId = parent.Id
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.accountRepo.Create(tx, account); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditAccountCreate, constant.AuditObjectAccount, account.Id, map[string]any{
		"code": account.Code,
		"type": string(account.Type),
	})

	return nil
}

func (u *accountUsecase) Update(actorUserId, id string, req *dto.UpdateAccountRequest) global.ErrorResponse {
	account, err := u.accountRepo.FindById(id)
	if err != nil {
		return err
	}

	code := normalizeAccountCode(req.Code)
	name := strings.TrimSpace(req.Name)
	accountType, err := parseAccountType(req.Type)
	if err != nil {
		return err
	}
	if code == "" {
		return global.BadRequestError("code is required")
	}
	if name == "" {
		return global.BadRequestError("name is required")
	}
	if err := u.ensureCodeAvailable(code, account.Id); err != nil {
		return err
	}
	parent, err := u.resolveParent(req.ParentId, accountType, account.Id)
	if err != nil {
		return err
	}

	account.Code = code
	account.Name = name
	account.Type = accountType
	account.ParentId = ""
	if parent != nil {
		account.ParentId = parent.Id
	}
	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.accountRepo.Update(tx, account); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditAccountUpdate, constant.AuditObjectAccount, account.Id, map[string]any{
		"code":      account.Code,
		"is_active": account.IsActive,
	})

	return nil
}

func (u *accountUsecase) Delete(actorUserId, id string) global.ErrorResponse {
	if _, err := u.accountRepo.FindById(id); err != nil {
		return err
	}

	children, err := u.accountRepo.CountByParentId(id)
	if err != nil {
		return err
	}
	if children > 0 {
		return global.BadRequestError("account has child accounts and cannot be deleted")
	}

	lines, err := u.journalEntryRepo.CountLinesByAccountId(id)
	if err != nil {
		return err
	}
	if lines > 0 {
		return global.BadRequestError("account is used in journal entries and cannot be deleted")
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.accountRepo.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditAccountDelete, constant.AuditObjectAccount, id, nil)

	return nil
}

// resolveParent enforces CoA hierarchy rules: the parent must exist, must not
// be the account itself, and must share its type.
func (u *accountUsecase) resolveParent(parentId string, accountType enum.AccountType, selfId string) (*model.Account, global.ErrorResponse) {
	parentId = strings.TrimSpace(parentId)
	if parentId == "" {
		return nil, nil
	}
	parent, err := u.accountRepo.FindById(parentId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return nil, global.BadRequestError("parent account not found")
		}
		return nil, err
	}
	if selfId != "" && parent.Id == selfId {
		return nil, global.BadRequestError("parent account cannot be the account itself")
	}
	if parent.Type != accountType {
		return nil, global.BadRequestError("parent account type must match account type")
	}
	return parent, nil
}

func (u *accountUsecase) ensureCodeAvailable(code, currentId string) global.ErrorResponse {
	existing, err := u.accountRepo.FindByCode(code)
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

func parseAccountType(value string) (enum.AccountType, global.ErrorResponse) {
	accountType := enum.AccountType(strings.ToUpper(strings.TrimSpace(value)))
	if !accountType.IsValid() {
		return "", global.BadRequestError("invalid account type")
	}
	return accountType, nil
}

func normalizeAccountCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
