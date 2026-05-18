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
	group.Post("/logout", h.Logout)
}

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

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	// Optional: extract user ID from locals if we need to do something in DB
	// userId := c.Locals("user_id").(string)

	err := h.authUsecase.Logout("")
	if err != nil {
		return err.ToResponse(c)
	}

	return global.MessageResponse("Logout successful", fiber.StatusOK, c)
}
