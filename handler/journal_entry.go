package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type JournalEntryHandler struct {
	usecase usecase.JournalEntryUsecase
}

func NewJournalEntryHandler(usecase usecase.JournalEntryUsecase) *JournalEntryHandler {
	return &JournalEntryHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all journal entries
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10, max 100)"
//	@Param			search	query		string	false	"Search by entry number or description"
//	@Success		200		{object}	global.Response[dto.PaginatedJournalEntryList]
//	@Router			/journals [get]
func (h *JournalEntryHandler) FindAll(c fiber.Ctx) error {
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
//	@Summary		Find journal entry by id
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Journal entry ID"
//	@Success		200	{object}	global.Response[dto.JournalEntryResponse]
//	@Router			/journals/{id} [get]
func (h *JournalEntryHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create journal entry
//	@Description	Post a balanced double-entry journal (total debit must equal total credit)
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateJournalEntryRequest	true	"Create journal entry request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/journals [post]
func (h *JournalEntryHandler) Create(c fiber.Ctx) error {
	var req dto.CreateJournalEntryRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("journal entry created successfully", fiber.StatusOK, c)
}
