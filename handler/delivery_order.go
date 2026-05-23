package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type DeliveryOrderHandler struct {
	usecase usecase.DeliveryOrderUsecase
}

func NewDeliveryOrderHandler(usecase usecase.DeliveryOrderUsecase) *DeliveryOrderHandler {
	return &DeliveryOrderHandler{usecase: usecase}
}

// Issue godoc
//
//	@Summary		Issue delivery order
//	@Description	Create delivery order manifest, validate weight vs fleet capacity, set cylinders to IN_TRANSIT
//	@Tags			Outbound
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.IssueDeliveryOrderRequest	true	"Issue DO request"
//	@Success		200		{object}	global.Response[dto.DeliveryOrderResponse]
//	@Router			/outbound/delivery-orders [post]
func (h *DeliveryOrderHandler) Issue(c fiber.Ctx) error {
	var req dto.IssueDeliveryOrderRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	res, err := h.usecase.Issue(actorUserId, &req)
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// FindAll godoc
//
//	@Summary		List delivery orders
//	@Description	List delivery orders with pagination and search
//	@Tags			Outbound
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by DO number or customer"
//	@Success		200		{object}	global.Response[dto.PaginatedDeliveryOrderList]
//	@Router			/outbound/delivery-orders [get]
func (h *DeliveryOrderHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Get delivery order by id
//	@Description	Get delivery order detail including manifest cylinders
//	@Tags			Outbound
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Delivery order ID"
//	@Success		200	{object}	global.Response[dto.DeliveryOrderResponse]
//	@Router			/outbound/delivery-orders/{id} [get]
func (h *DeliveryOrderHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}
