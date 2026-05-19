package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type CustomerHandler struct {
	usecase usecase.CustomerUsecase
}

func NewCustomerHandler(usecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all customers
//	@Description	Find all customers with pagination and search
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by code, name, or phone"
//	@Success		200		{object}	global.Response[dto.PaginatedCustomerList]
//	@Router			/customers [get]
func (h *CustomerHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Find customer by id
//	@Description	Find customer by id
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Customer ID"
//	@Success		200	{object}	global.Response[dto.CustomerResponse]
//	@Router			/customers/{id} [get]
func (h *CustomerHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create customer
//	@Description	Create customer with cylinder quota limit
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateCustomerRequest	true	"Create customer request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/customers [post]
func (h *CustomerHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCustomerRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update customer
//	@Description	Update customer by id
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string						true	"Customer ID"
//	@Param			request	body		dto.UpdateCustomerRequest	true	"Update customer request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/customers/{id} [put]
func (h *CustomerHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateCustomerRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer updated successfully", fiber.StatusOK, c)
}
