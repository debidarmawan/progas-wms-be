package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authUsecase usecase.AuthUseCase
}

func NewAuthHandler(authUsecase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUsecase}
}

func (h *AuthHandler) Routes(group fiber.Router) {
	group.Post("/login", h.Login)
	group.Post("/refresh-token", h.RefreshToken)
}

// Login godoc
//
//	@Summary		Login
//	@Description	Login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.LoginRequest	true	"Login data"
//	@Success		200		{object}	global.Response[dto.LoginResponse]
//	@Router			/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}

	res, err := h.authUsecase.Login(&req)
	if err != nil {
		return err.ToResponse(c)
	}

	return global.CreateResponse(res, fiber.StatusOK, c)
}

// RefreshToken godoc
//
//	@Summary		Refresh token
//	@Description	Refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.RefreshTokenRequest	true	"RefreshToken data"
//	@Success		200		{object}	global.Response[dto.LoginResponse]
//	@Router			/refresh-token [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}

	res, err := h.authUsecase.RefreshToken(&req)
	if err != nil {
		return err.ToResponse(c)
	}

	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logout
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(string)
	if !ok || userId == "" {
		return global.UnauthorizedError().ToResponse(c)
	}

	err := h.authUsecase.Logout(userId)
	if err != nil {
		return err.ToResponse(c)
	}

	return global.MessageResponse("Logout successful", fiber.StatusOK, c)
}

// Profile godoc
//
//	@Summary		Get profile
//	@Description	Get current user profile
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	global.Response[dto.UserResponse]
//	@Router			/profile [get]
func (h *AuthHandler) Profile(c fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(string)
	if !ok || userId == "" {
		return global.UnauthorizedError().ToResponse(c)
	}

	res, err := h.authUsecase.Profile(userId)
	if err != nil {
		return err.ToResponse(c)
	}

	return global.CreateResponse(res, fiber.StatusOK, c)
}
