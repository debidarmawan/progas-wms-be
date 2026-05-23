package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type StockOpnameHandler struct {
	usecase usecase.StockOpnameUsecase
}

func NewStockOpnameHandler(usecase usecase.StockOpnameUsecase) *StockOpnameHandler {
	return &StockOpnameHandler{usecase: usecase}
}

// Submit godoc
//
//	@Summary		Spare part stock opname
//	@Description	Set actual spare part quantity (adjusts stock to counted value)
//	@Tags			Inventory
//	@Security		Bearer
//	@Param			request	body		dto.StockOpnameRequest	true	"Request"
//	@Success		200		{object}	global.Response[dto.StockOpnameResponse]
//	@Router			/inventory/spareparts/stock-opname [post]
func (h *StockOpnameHandler) Submit(c fiber.Ctx) error {
	var req dto.StockOpnameRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Submit(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
