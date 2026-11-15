package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type TermOfPaymentHandler struct {
	usecase usecase.TermOfPaymentUsecase
}

func NewTermOfPaymentHandler(usecase usecase.TermOfPaymentUsecase) *TermOfPaymentHandler {
	return &TermOfPaymentHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all terms of payment
//	@Description	Find all payment terms with pagination and search
//	@Tags			Term Of Payment
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by name or code"
//	@Success		200		{object}	global.Response[dto.PaginatedTermOfPaymentList]
//	@Router			/terms-of-payment [get]
func (h *TermOfPaymentHandler) FindAll(c fiber.Ctx) error {
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

// FindById godoc
//
//	@Summary		Find term of payment by id
//	@Description	Find a payment term by id
//	@Tags			Term Of Payment
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Term of payment ID"
//	@Success		200	{object}	global.Response[dto.TermOfPaymentResponse]
//	@Router			/terms-of-payment/{id} [get]
func (h *TermOfPaymentHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create term of payment
//	@Description	Create a payment term (days = 0 means cash/COD)
//	@Tags			Term Of Payment
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateTermOfPaymentRequest	true	"Create term of payment request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/terms-of-payment [post]
func (h *TermOfPaymentHandler) Create(c fiber.Ctx) error {
	var req dto.CreateTermOfPaymentRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Term of payment created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update term of payment
//	@Description	Update a payment term by id
//	@Tags			Term Of Payment
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string							true	"Term of payment ID"
//	@Param			request	body		dto.UpdateTermOfPaymentRequest	true	"Update term of payment request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/terms-of-payment/{id} [put]
func (h *TermOfPaymentHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateTermOfPaymentRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Term of payment updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete term of payment
//	@Description	Delete a payment term by id
//	@Tags			Term Of Payment
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Term of payment ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/terms-of-payment/{id} [delete]
func (h *TermOfPaymentHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Term of payment deleted successfully", fiber.StatusOK, c)
}
