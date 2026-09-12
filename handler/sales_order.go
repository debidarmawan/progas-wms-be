package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type SalesOrderHandler struct {
	usecase usecase.SalesOrderUsecase
}

func NewSalesOrderHandler(usecase usecase.SalesOrderUsecase) *SalesOrderHandler {
	return &SalesOrderHandler{usecase: usecase}
}

func (h *SalesOrderHandler) FindAll(c fiber.Ctx) error {
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

func (h *SalesOrderHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

func (h *SalesOrderHandler) Create(c fiber.Ctx) error {
	var req dto.CreateSalesOrderRequest
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

func (h *SalesOrderHandler) Confirm(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Confirm(actorUserId, c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

func (h *SalesOrderHandler) Cancel(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Cancel(actorUserId, c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
