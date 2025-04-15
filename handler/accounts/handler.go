package accounts

import (
	"errors"
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
		return nil, errors.New("bind error")
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
		return nil, errors.New("invalid id")
	}

	account, err := h.accountSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (h *accounts) GetAll(ctx *gofr.Context) (interface{}, error) {
	account, err := h.accountSvc.GetAll(ctx, &filters.Account{})
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (h *accounts) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	var account *models.Account

	err = ctx.Bind(&account)
	if err != nil {
		return nil, errors.New("bind error")
	}

	account.ID = id

	updatedAccount, err := h.accountSvc.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (h *accounts) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	err = h.accountSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "account deleted successfully", nil
}
