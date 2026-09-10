package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type CustomerItemPriceHandler struct {
	usecase usecase.CustomerItemPriceUsecase
}

func NewCustomerItemPriceHandler(usecase usecase.CustomerItemPriceUsecase) *CustomerItemPriceHandler {
	return &CustomerItemPriceHandler{usecase: usecase}
}

// FindAllByCustomer godoc
//
//	@Summary	Find customer item prices
//	@Tags		Customer Pricing
//	@Security	Bearer
//	@Param		customerId	path		string	true	"Customer ID"
//	@Param		page		query		int		false	"Page number"
//	@Param		limit		query		int		false	"Items per page"
//	@Param		search		query		string	false	"Search by item name or SKU"
//	@Success	200			{object}	global.Response[dto.PaginatedCustomerItemPriceList]
//	@Router		/customers/{customerId}/pricing [get]
func (h *CustomerItemPriceHandler) FindAllByCustomer(c fiber.Ctx) error {
	var query dto.ListQuery
	if err := helper.ValidateQuery(c, &query); err != nil {
		return err.ToResponse(c)
	}
	res, err := h.usecase.FindAllByCustomer(c.Params("customerId"), &query)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// FindActive godoc
//
//	@Summary	Resolve customer item price
//	@Tags		Customer Pricing
//	@Security	Bearer
//	@Param		customerId		path		string	true	"Customer ID"
//	@Param		masterItemId	path		string	true	"Master item ID"
//	@Success	200				{object}	global.Response[dto.CustomerItemPriceResponse]
//	@Router		/customers/{customerId}/items/{masterItemId}/pricing [get]
func (h *CustomerItemPriceHandler) FindActive(c fiber.Ctx) error {
	res, err := h.usecase.FindActive(c.Params("customerId"), c.Params("masterItemId"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary	Create customer item price
//	@Tags		Customer Pricing
//	@Security	Bearer
//	@Param		customerId	path		string								true	"Customer ID"
//	@Param		request		body		dto.CreateCustomerItemPriceRequest	true	"Customer item price request"
//	@Success	200			{object}	global.Response[dto.Message]
//	@Router		/customers/{customerId}/pricing [post]
func (h *CustomerItemPriceHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCustomerItemPriceRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, c.Params("customerId"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer item price created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary	Update future customer item price
//	@Tags		Customer Pricing
//	@Security	Bearer
//	@Param		customerId	path		string								true	"Customer ID"
//	@Param		id			path		string								true	"Price ID"
//	@Param		request		body		dto.UpdateCustomerItemPriceRequest	true	"Customer item price request"
//	@Success	200			{object}	global.Response[dto.Message]
//	@Router		/customers/{customerId}/pricing/{id} [put]
func (h *CustomerItemPriceHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateCustomerItemPriceRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("customerId"), c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer item price updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary	Delete customer item price
//	@Tags		Customer Pricing
//	@Security	Bearer
//	@Param		customerId	path		string	true	"Customer ID"
//	@Param		id			path		string	true	"Price ID"
//	@Success	200			{object}	global.Response[dto.Message]
//	@Router		/customers/{customerId}/pricing/{id} [delete]
func (h *CustomerItemPriceHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("customerId"), c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer item price deleted successfully", fiber.StatusOK, c)
}
