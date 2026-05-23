package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type WorkOrderHandler struct {
	usecase usecase.WorkOrderUsecase
}

func NewWorkOrderHandler(usecase usecase.WorkOrderUsecase) *WorkOrderHandler {
	return &WorkOrderHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary	List work orders
//	@Tags		Maintenance
//	@Security	Bearer
//	@Param		page	query		int		false	"Page"
//	@Param		limit	query		int		false	"Limit"
//	@Param		search	query		string	false	"Search"
//	@Success	200		{object}	global.Response[dto.PaginatedWorkOrderList]
//	@Router		/maintenance/work-orders [get]
func (h *WorkOrderHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary	Get work order by id
//	@Tags		Maintenance
//	@Security	Bearer
//	@Param		id	path		string	true	"Work order ID"
//	@Success	200	{object}	global.Response[dto.WorkOrderResponse]
//	@Router		/maintenance/work-orders/{id} [get]
func (h *WorkOrderHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary	Create work order
//	@Tags		Maintenance
//	@Security	Bearer
//	@Param		request	body		dto.CreateWorkOrderRequest	true	"Request"
//	@Success	200		{object}	global.Response[dto.WorkOrderResponse]
//	@Router		/maintenance/work-orders [post]
func (h *WorkOrderHandler) Create(c fiber.Ctx) error {
	var req dto.CreateWorkOrderRequest
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

// Complete godoc
//
//	@Summary		Complete work order
//	@Description	Deduct spare part stock when completing
//	@Tags			Maintenance
//	@Security		Bearer
//	@Param			id	path		string	true	"Work order ID"
//	@Success		200	{object}	global.Response[dto.WorkOrderResponse]
//	@Router			/maintenance/work-orders/{id}/complete [post]
func (h *WorkOrderHandler) Complete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Complete(actorUserId, c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
