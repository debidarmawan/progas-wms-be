package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type CustomerSiteHandler struct {
	usecase usecase.CustomerSiteUsecase
}

func NewCustomerSiteHandler(usecase usecase.CustomerSiteUsecase) *CustomerSiteHandler {
	return &CustomerSiteHandler{usecase: usecase}
}

// FindByCustomer godoc
//
//	@Summary		List customer sites
//	@Description	List every site/branch address of a customer
//	@Tags			Customer Site
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			customerId	path		string	true	"Customer ID"
//	@Success		200			{object}	global.Response[dto.CustomerSiteListResponse]
//	@Router			/customers/{customerId}/sites [get]
func (h *CustomerSiteHandler) FindByCustomer(c fiber.Ctx) error {
	res, err := h.usecase.FindByCustomer(c.Params("customerId"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create customer site
//	@Description	Create a site/branch address used as delivery ship-to
//	@Tags			Customer Site
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			customerId	path		string						true	"Customer ID"
//	@Param			request		body		dto.CreateCustomerSiteRequest	true	"Create customer site request"
//	@Success		200			{object}	global.Response[dto.Message]
//	@Router			/customers/{customerId}/sites [post]
func (h *CustomerSiteHandler) Create(c fiber.Ctx) error {
	var req dto.CreateCustomerSiteRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, c.Params("customerId"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer site created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update customer site
//	@Description	Update a site/branch address of a customer
//	@Tags			Customer Site
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			customerId	path		string						true	"Customer ID"
//	@Param			id			path		string						true	"Customer site ID"
//	@Param			request		body		dto.UpdateCustomerSiteRequest	true	"Update customer site request"
//	@Success		200			{object}	global.Response[dto.Message]
//	@Router			/customers/{customerId}/sites/{id} [put]
func (h *CustomerSiteHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateCustomerSiteRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("customerId"), c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer site updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete customer site
//	@Description	Soft delete a site/branch address of a customer (existing delivery orders keep their ship-to snapshot)
//	@Tags			Customer Site
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			customerId	path		string	true	"Customer ID"
//	@Param			id			path		string	true	"Customer site ID"
//	@Success		200			{object}	global.Response[dto.Message]
//	@Router			/customers/{customerId}/sites/{id} [delete]
func (h *CustomerSiteHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("customerId"), c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Customer site deleted successfully", fiber.StatusOK, c)
}
