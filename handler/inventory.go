package handler

import (
	"progas-wms-be/global"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type InventoryHandler struct {
	usecase usecase.InventoryUsecase
}

func NewInventoryHandler(usecase usecase.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{usecase: usecase}
}

// VirtualWarehouse godoc
//
//	@Summary		Virtual warehouse
//	@Description	Customers with outstanding cylinders (company inventory at customer sites)
//	@Tags			Inventory
//	@Security		Bearer
//	@Success		200	{object}	global.Response[dto.VirtualWarehouseResponse]
//	@Router			/inventory/virtual-warehouse [get]
func (h *InventoryHandler) VirtualWarehouse(c fiber.Ctx) error {
	res, err := h.usecase.VirtualWarehouse()
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
