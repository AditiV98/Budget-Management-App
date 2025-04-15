package savings

import (
	"errors"
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
	var saving *models.Savings

	err := ctx.Bind(&saving)
	if err != nil {
		return nil, errors.New("bind error")
	}

	newSaving, err := h.savingsSvc.Create(ctx, saving)
	if err != nil {
		return nil, err
	}

	return newSaving, nil
}

func (h *savings) GetAll(ctx *gofr.Context) (interface{}, error) {
	saving, err := h.savingsSvc.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return saving, nil
}

func (h *savings) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	saving, err := h.savingsSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return saving, nil
}

func (h *savings) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	var saving *models.Savings

	err = ctx.Bind(&saving)
	if err != nil {
		return nil, errors.New("bind error")
	}

	saving.ID = id

	updatedSaving, err := h.savingsSvc.Update(ctx, saving)
	if err != nil {
		return nil, err
	}

	return updatedSaving, nil
}

func (h *savings) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	err = h.savingsSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "saving deleted successfully", nil
}
