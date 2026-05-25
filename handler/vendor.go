package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type VendorHandler struct {
	usecase usecase.VendorUsecase
}

func NewVendorHandler(usecase usecase.VendorUsecase) *VendorHandler {
	return &VendorHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		List vendors
//	@Description	List rental vendors with cylinder count
//	@Tags			Vendor
//	@Security		Bearer
//	@Param			page	query		int		false	"Page"
//	@Param			limit	query		int		false	"Limit"
//	@Param			search	query		string	false	"Search code, name, contact, phone, email"
//	@Success		200		{object}	global.Response[dto.PaginatedVendorList]
//	@Router			/vendors [get]
func (h *VendorHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Get vendor by id
//	@Description	Get vendor detail including rented cylinders (ownership VENDOR)
//	@Tags			Vendor
//	@Security		Bearer
//	@Param			id	path		string	true	"Vendor ID"
//	@Success		200	{object}	global.Response[dto.VendorDetailResponse]
//	@Router			/vendors/{id} [get]
func (h *VendorHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create vendor
//	@Description	Register rental cylinder vendor
//	@Tags			Vendor
//	@Security		Bearer
//	@Param			request	body		dto.CreateVendorRequest	true	"Request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/vendors [post]
func (h *VendorHandler) Create(c fiber.Ctx) error {
	var req dto.CreateVendorRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Vendor created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary	Update vendor
//	@Tags		Vendor
//	@Security	Bearer
//	@Param		id		path		string					true	"Vendor ID"
//	@Param		request	body		dto.UpdateVendorRequest	true	"Request"
//	@Success	200		{object}	global.Response[dto.Message]
//	@Router		/vendors/{id} [put]
func (h *VendorHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateVendorRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Vendor updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete vendor
//	@Description	Soft-delete vendor (only when no rental cylinders linked)
//	@Tags			Vendor
//	@Security		Bearer
//	@Param			id	path		string	true	"Vendor ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/vendors/{id} [delete]
func (h *VendorHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Vendor deleted successfully", fiber.StatusOK, c)
}
