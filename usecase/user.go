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
	txManager        helper.TxManager
	userRepository   repository.UserRepository
	roleRepository   repository.RoleRepository
	driverRepository repository.DriverRepository
	auditLogRepo     repository.AuditLogRepository
}

func NewUserUsecase(
	txManager helper.TxManager,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	driverRepository repository.DriverRepository,
	auditLogRepo repository.AuditLogRepository,
) UserUsecase {
	return &userUsecase{
		txManager:        txManager,
		userRepository:   userRepository,
		roleRepository:   roleRepository,
		driverRepository: driverRepository,
		auditLogRepo:     auditLogRepo,
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

	driverId, err := u.resolveDriverAssignment(request.RoleId, request.DriverId, "")
	if err != nil {
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
		DriverId: driverId,
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

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUserCreate, constant.AuditObjectUser, user.Id, map[string]any{
		"email":     user.Email,
		"name":      user.Name,
		"role_id":   user.RoleId,
		"driver_id": user.DriverId,
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

	driverId, err := u.resolveDriverAssignment(request.RoleId, request.DriverId, id)
	if err != nil {
		return err
	}

	if actorUserId == id && !request.IsActive {
		return global.BadRequestError("cannot deactivate your own account")
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Phone = request.Phone
	user.RoleId = request.RoleId
	user.DriverId = driverId
	user.Driver = nil
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
		"role_id":   user.RoleId,
		"driver_id": user.DriverId,
		"is_active": user.IsActive,
	})

	return nil
}

func (u *userUsecase) DeleteUser(actorUserId, id string) global.ErrorResponse {
	if actorUserId == id {
		return global.BadRequestError("cannot delete your own account")
	}

	user, err := u.userRepository.FindById(id)
	if err != nil {
		return err
	}
	driverId := user.DriverId

	tx := u.txManager.New()
	defer tx.CheckPanic()

	user.DriverId = nil
	user.Driver = nil
	if err := u.userRepository.Update(tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if err := u.userRepository.Delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return global.InternalServerError(err)
	}

	_ = u.auditLogRepo.Log(actorUserId, constant.AuditUserDelete, constant.AuditObjectUser, id, map[string]any{
		"driver_id": driverId,
	})

	return nil
}

func (u *userUsecase) resolveDriverAssignment(roleId, driverId, excludeUserId string) (*string, global.ErrorResponse) {
	role, err := u.roleRepository.FindById(roleId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return nil, global.BadRequestError("invalid role")
		}
		return nil, err
	}

	if role.Name != constant.RoleDriver {
		if driverId != "" {
			return nil, global.BadRequestError("driver_id is only allowed for the Driver role")
		}
		return nil, nil
	}

	if driverId == "" {
		return nil, global.BadRequestError("driver_id is required for the Driver role")
	}

	driver, err := u.driverRepository.FindById(driverId)
	if err != nil {
		if err.GetCode() == fiber.StatusNotFound {
			return nil, global.BadRequestError("invalid driver")
		}
		return nil, err
	}
	if !driver.IsActive && excludeUserId == "" {
		return nil, global.BadRequestError("driver is not active")
	}
	if !driver.IsActive {
		currentUser, err := u.userRepository.FindById(excludeUserId)
		if err != nil {
			return nil, err
		}
		if currentUser.DriverId == nil || *currentUser.DriverId != driverId {
			return nil, global.BadRequestError("driver is not active")
		}
	}

	assignedUser, err := u.userRepository.FindByDriverIdExceptId(driverId, excludeUserId)
	if err != nil && err.GetCode() != fiber.StatusNotFound {
		return nil, err
	}
	if assignedUser != nil {
		return nil, global.BadRequestError("driver is already assigned to a user")
	}

	return &driverId, nil
}
