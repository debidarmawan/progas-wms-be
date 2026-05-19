package usecase

import (
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type UserUsecase interface {
	CreateUser(actorUserId string, request *dto.CreateUserRequest) global.ErrorResponse
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

func (u *userUsecase) CreateUser(actorUserId string, request *dto.CreateUserRequest) global.ErrorResponse {
	existingUser, err := u.userRepository.FindByEmail(request.Email)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}

	if existingUser != nil {
		return global.BadRequestError("User already exists")
	}

	_, err = u.roleRepository.FindById(request.RoleId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return global.BadRequestError("invalid role")
		}
		return err
	}

	hashedPassword, hashErr := helper.HashPassword(request.Password)
	if hashErr != nil {
		return global.InternalServerError(hashErr)
	}

	user := &model.User{
		Name:     request.Name,
		Email:    request.Email,
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
