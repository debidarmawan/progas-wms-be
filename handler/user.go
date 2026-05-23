package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase}
}

// FindAll godoc
//
//	@Summary		List users
//	@Description	List users with pagination and search (Superadmin & Manager only)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by name, email, or phone"
//	@Success		200		{object}	global.Response[dto.PaginatedUserList]
//	@Router			/users [get]
func (h *UserHandler) FindAll(c fiber.Ctx) error {
	var query dto.ListQuery
	if err := helper.ValidateQuery(c, &query); err != nil {
		return err.ToResponse(c)
	}
	res, err := h.userUsecase.FindAll(&query)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// FindById godoc
//
//	@Summary		Get user by id
//	@Description	Get user detail (Superadmin & Manager only)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	global.Response[dto.UserListResponse]
//	@Router			/users/{id} [get]
func (h *UserHandler) FindById(c fiber.Ctx) error {
	res, err := h.userUsecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create user (Superadmin & Manager only)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateUserRequest	true	"Create user request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.userUsecase.CreateUser(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("User created successfully", fiber.StatusOK, c)
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update user profile, role, active status, or password (Superadmin & Manager only)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string					true	"User ID"
//	@Param			request	body		dto.UpdateUserRequest	true	"Update user request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	var req dto.UpdateUserRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.userUsecase.UpdateUser(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("User updated successfully", fiber.StatusOK, c)
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Soft-delete user (Superadmin & Manager only)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.userUsecase.DeleteUser(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("User deleted successfully", fiber.StatusOK, c)
}
