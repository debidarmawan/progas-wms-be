package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type DriverHandler struct {
	usecase usecase.DriverUsecase
}

func NewDriverHandler(usecase usecase.DriverUsecase) *DriverHandler {
	return &DriverHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		List drivers
//	@Description	List drivers with pagination and search
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by name, phone, or license number"
//	@Success		200		{object}	global.Response[dto.PaginatedDriverList]
//	@Router			/logistics/drivers [get]
func (h *DriverHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Get driver by id
//	@Description	Get driver detail
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Driver ID"
//	@Success		200	{object}	global.Response[dto.DriverResponse]
//	@Router			/logistics/drivers/{id} [get]
func (h *DriverHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create driver
//	@Description	Register a delivery driver
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateDriverRequest	true	"Create driver request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/logistics/drivers [post]
func (h *DriverHandler) Create(c fiber.Ctx) error {
	var req dto.CreateDriverRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Driver created", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update driver
//	@Description	Update driver info or active status
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string					true	"Driver ID"
//	@Param			request	body		dto.UpdateDriverRequest	true	"Update driver request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/logistics/drivers/{id} [put]
func (h *DriverHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateDriverRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Driver updated", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete driver
//	@Description	Soft-delete a driver
//	@Tags			Logistics
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Driver ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/logistics/drivers/{id} [delete]
func (h *DriverHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("Driver deleted", fiber.StatusOK, c)
}
