package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type CustomerPOHandler struct {
	usecase usecase.CustomerPOUsecase
}

func NewCustomerPOHandler(usecase usecase.CustomerPOUsecase) *CustomerPOHandler {
	return &CustomerPOHandler{usecase: usecase}
}

func (h *CustomerPOHandler) FindAll(c fiber.Ctx) error {
	var query dto.ListQuery
	if err := helper.ValidateQuery(c, &query); err != nil {
		return err.ToResponse(c)
	}
	res, err := h.usecase.FindAll(&query)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

func (h *CustomerPOHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

func (h *CustomerPOHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCustomerPORequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Create(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

func (h *CustomerPOHandler) Confirm(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Confirm(actorUserId, c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
