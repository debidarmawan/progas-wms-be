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

func (h *UserHandler) Routes(group fiber.Router) {
	group.Post("/users", h.CreateUser)
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create user
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
	err := h.userUsecase.CreateUser(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}

	return global.CreateMessageResponse("User created successfully", fiber.StatusOK, c)
}
