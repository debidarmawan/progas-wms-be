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

type UserUsecase interface {
	FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.UserListResponse], global.ErrorResponse)
	FindById(id string) (*dto.UserListResponse, global.ErrorResponse)
	CreateUser(actorUserId string, request *dto.CreateUserRequest) global.ErrorResponse
	UpdateUser(actorUserId, id string, request *dto.UpdateUserRequest) global.ErrorResponse
	DeleteUser(actorUserId, id string) global.ErrorResponse
}

type userUsecase struct {
	txManager      helper.TxManager
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
	auditLogRepo   repository.AuditLogRepository
}

func NewUserUsecase(
	txManager helper.TxManager,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	auditLogRepo repository.AuditLogRepository,
) UserUsecase {
	return &userUsecase{
		txManager:      txManager,
		userRepository: userRepository,
		roleRepository: roleRepository,
		auditLogRepo:   auditLogRepo,
	}
}

func (u *userUsecase) FindAll(query *dto.ListQuery) (*dto.PaginatedResponse[dto.UserListResponse], global.ErrorResponse) {
	page, limit, _ := helper.NormalizePagination(query)
	search := helper.NormalizeSearch(query.Search)
	users, total, err := u.userRepository.FindAll(page, limit, search)
	if err != nil {
		return nil, err
	}
	return &dto.PaginatedResponse[dto.UserListResponse]{
		Items: mapper.ToUserListResponses(users),
		Meta:  helper.BuildPaginationMeta(page, limit, total),
	}, nil
}

func (u *userUsecase) FindById(id string) (*dto.UserListResponse, global.ErrorResponse) {
	user, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToUserListResponse(user), nil
}

func (u *userUsecase) CreateUser(actorUserId string, request *dto.CreateUserRequest) global.ErrorResponse {
	existingUser, err := u.userRepository.FindByEmail(request.Email)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}
	if existingUser != nil {
		return global.BadRequestError("email already in use")
	}

	if err := u.validateRoleId(request.RoleId); err != nil {
		return err
	}

	hashedPassword, hashErr := helper.HashPassword(request.Password)
	if hashErr != nil {
		return global.InternalServerError(hashErr)
	}

	user := &model.User{
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hashedPassword,
		RoleId:   request.RoleId,
		IsActive: true,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.userRepository.Create(tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUserCreate, constant.AuditObjectUser, user.Id, map[string]string{
		"email":   user.Email,
		"name":    user.Name,
		"role_id": user.RoleId,
	})

	return nil
}

func (u *userUsecase) UpdateUser(actorUserId, id string, request *dto.UpdateUserRequest) global.ErrorResponse {
	user, err := u.userRepository.FindById(id)
	if err != nil {
		return err
	}

	duplicate, dupErr := u.userRepository.FindByEmailExceptId(request.Email, id)
	if dupErr != nil && dupErr.GetCode() != fiber.StatusNotFound {
		return dupErr
	}
	if duplicate != nil {
		return global.BadRequestError("email already in use")
	}

	if err := u.validateRoleId(request.RoleId); err != nil {
		return err
	}

	if actorUserId == id && !request.IsActive {
		return global.BadRequestError("cannot deactivate your own account")
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Phone = request.Phone
	user.RoleId = request.RoleId
	user.IsActive = request.IsActive

	if request.Password != "" {
		hashedPassword, hashErr := helper.HashPassword(request.Password)
		if hashErr != nil {
			return global.InternalServerError(hashErr)
		}
		user.Password = hashedPassword
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err = u.userRepository.Update(tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUserUpdate, constant.AuditObjectUser, user.Id, map[string]any{
		"email":     user.Email,
		"is_active": user.IsActive,
	})

	return nil
}

func (u *userUsecase) DeleteUser(actorUserId, id string) global.ErrorResponse {
	if actorUserId == id {
		return global.BadRequestError("cannot delete your own account")
	}

	if _, err := u.userRepository.FindById(id); err != nil {
		return err
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	if err := u.userRepository.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUserDelete, constant.AuditObjectUser, id, nil)

	return nil
}

func (u *userUsecase) validateRoleId(roleId string) global.ErrorResponse {
	_, err := u.roleRepository.FindById(roleId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return global.BadRequestError("invalid role")
		}
		return err
	}
	return nil
}
