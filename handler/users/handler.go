package users

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
	"strconv"
	"strings"
)

type usersHandler struct {
	userSvc services.User
}

func New(userSvc services.User) handler.User {
	return &usersHandler{userSvc: userSvc}
}

func (h *usersHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var user *models.User

	err := ctx.Bind(&user)
	if err != nil {
		return nil, err
	}

	err = h.userSvc.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return "user created successfully", nil
}

func (h *usersHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	user, err := h.userSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *usersHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	users, err := h.userSvc.GetAll(ctx, &filters.User{})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (h *usersHandler) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	var user *models.User

	err = ctx.Bind(&user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	err = h.userSvc.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return "user updated successfully", nil
}

func (h *usersHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	err = h.userSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "user deleted successfully", nil
}
