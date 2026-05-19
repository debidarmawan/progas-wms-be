package usecase

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/model"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

type UserUsecase interface {
	CreateUser(request *dto.CreateUserRequest) global.ErrorResponse
}

type userUsecase struct {
	txManager      helper.TxManager
	userRepository repository.UserRepository
}

func NewUserUsecase(
	txManager helper.TxManager,
	userRepository repository.UserRepository,
) UserUsecase {
	return &userUsecase{
		txManager:      txManager,
		userRepository: userRepository,
	}
}

func (u *userUsecase) CreateUser(request *dto.CreateUserRequest) global.ErrorResponse {
	existingUser, err := u.userRepository.FindByEmail(request.Email)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return err
	}

	if existingUser != nil {
		return global.BadRequestError("User already exists")
	}

	user := &model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: helper.HashPassword(request.Password),
		RoleId:   request.RoleId,
	}

	tx := u.txManager.New()
	defer tx.CheckPanic()

	err = u.userRepository.Create(tx, user)
	if err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	return nil
}
