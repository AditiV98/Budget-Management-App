package accounts

import (
	"moneyManagement/filters"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
	"strconv"
	"strings"

	"gofr.dev/pkg/gofr"
)

type accounts struct {
	accountSvc services.Account
}

func New(accountSvc services.Account) handler.Account {
	return &accounts{accountSvc: accountSvc}
}

func (h *accounts) Create(ctx *gofr.Context) (interface{}, error) {
	var account *models.Account

	err := ctx.Bind(&account)
	if err != nil {
		return nil, err
	}

	newAccount, err := h.accountSvc.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (h *accounts) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	account, err := h.accountSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (h *accounts) GetAll(ctx *gofr.Context) (interface{}, error) {
	accounts, err := h.accountSvc.GetAll(ctx,&filters.Account{})
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (h *accounts) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	var account *models.Account

	err = ctx.Bind(&account)
	if err != nil {
		return nil, err
	}

	account.ID = id

	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	updatedAccount, err := h.accountSvc.Update(ctx, account, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (h *accounts) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	err = h.accountSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "account deleted successfully", nil
}
