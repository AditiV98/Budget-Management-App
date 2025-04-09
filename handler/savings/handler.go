package savings

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
	"strconv"
	"strings"
)

type savings struct {
	savingsSvc services.Savings
}

func New(savingsSvc services.Savings) handler.Savings {
	return &savings{savingsSvc: savingsSvc}
}

func (h *savings) Create(ctx *gofr.Context) (interface{}, error) {
	var savings *models.Savings

	err := ctx.Bind(&savings)
	if err != nil {
		return nil, err
	}

	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = h.savingsSvc.Create(ctx, savings, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *savings) GetAll(ctx *gofr.Context) (interface{}, error) {
	savings, err := h.savingsSvc.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return savings, nil
}

func (h *savings) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	savings, err := h.savingsSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return savings, nil
}

func (h *savings) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	var savings *models.Savings

	err = ctx.Bind(&savings)
	if err != nil {
		return nil, err
	}

	savings.ID = id

	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = h.savingsSvc.Update(ctx, savings, false, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *savings) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	err = h.savingsSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
