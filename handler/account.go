package handler

import (
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/usecase"

	"github.com/gofiber/fiber/v3"
)

type AccountHandler struct {
	usecase usecase.AccountUsecase
}

func NewAccountHandler(usecase usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{usecase: usecase}
}

// FindAll godoc
//
//	@Summary		Find all accounts
//	@Description	List Chart of Accounts with pagination, search and filters
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			page		query		int		false	"Page number (default 1)"
//	@Param			limit		query		int		false	"Items per page (default 10, max 100)"
//	@Param			search		query		string	false	"Search by code or name"
//	@Param			type		query		string	false	"ASSET | LIABILITY | EQUITY | REVENUE | EXPENSE"
//	@Param			is_active	query		bool	false	"Filter by active flag"
//	@Success		200			{object}	global.Response[dto.PaginatedAccountList]
//	@Router			/accounts [get]
func (h *AccountHandler) FindAll(c fiber.Ctx) error {
	var query dto.AccountListQuery
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
//	@Summary		Find account by id
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Account ID"
//	@Success		200	{object}	global.Response[dto.AccountResponse]
//	@Router			/accounts/{id} [get]
func (h *AccountHandler) FindById(c fiber.Ctx) error {
	res, err := h.usecase.FindById(c.Params("id"))
	if err != nil {
		return err.ToResponse(c)
	}
	return global.CreateResponse(res, fiber.StatusOK, c)
}

// Create godoc
//
//	@Summary		Create account
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body		dto.CreateAccountRequest	true	"Create account request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/accounts [post]
func (h *AccountHandler) Create(c fiber.Ctx) error {
	var req dto.CreateAccountRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Create(actorUserId, &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("account created successfully", fiber.StatusOK, c)
}

// Update godoc
//
//	@Summary		Update account
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		string						true	"Account ID"
//	@Param			request	body		dto.UpdateAccountRequest	true	"Update account request"
//	@Success		200		{object}	global.Response[dto.Message]
//	@Router			/accounts/{id} [put]
func (h *AccountHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateAccountRequest
	if err := helper.ValidateBody(c, &req); err != nil {
		return err.ToResponse(c)
	}
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Update(actorUserId, c.Params("id"), &req); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("account updated successfully", fiber.StatusOK, c)
}

// Delete godoc
//
//	@Summary		Delete account
//	@Description	Delete an account (blocked when it has children or is used in journal entries)
//	@Tags			Accounting
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		string	true	"Account ID"
//	@Success		200	{object}	global.Response[dto.Message]
//	@Router			/accounts/{id} [delete]
func (h *AccountHandler) Delete(c fiber.Ctx) error {
	actorUserId, _ := c.Locals("user_id").(string)
	if err := h.usecase.Delete(actorUserId, c.Params("id")); err != nil {
		return err.ToResponse(c)
	}
	return global.CreateMessageResponse("account deleted successfully", fiber.StatusOK, c)
}
