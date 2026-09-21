package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type InvoiceHandler struct {
	usecase        usecase.InvoiceUsecase
	paymentUsecase usecase.PaymentUsecase
}

func NewInvoiceHandler(usecase usecase.InvoiceUsecase, paymentUsecase usecase.PaymentUsecase) *InvoiceHandler {
	return &InvoiceHandler{usecase: usecase, paymentUsecase: paymentUsecase}
}

func (h *InvoiceHandler) FindAll(c fiber.Ctx) error {
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

func (h *InvoiceHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

func (h *InvoiceHandler) RecordPayment(c fiber.Ctx) error {
	var req dto.RecordPaymentRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.paymentUsecase.Record(actorUserId, c.Params("id"), &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
