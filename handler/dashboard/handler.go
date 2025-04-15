package dashboard

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/handler"
	"moneyManagement/services"
	"strconv"
)

type dashboardHandler struct {
	dashboardSvc services.Dashboard
}

func New(dashboardSvc services.Dashboard) handler.Dashboard {
	return &dashboardHandler{dashboardSvc: dashboardSvc}
}

func (h *dashboardHandler) Get(ctx *gofr.Context) (interface{}, error) {
	startDate := ctx.Params("startDate")
	endDate := ctx.Params("endDate")
	accountID := ctx.Params("accountId")

	id, err := strconv.Atoi(accountID[0])
	if err != nil {
		return nil, errors.New("invalid id")
	}

	f := &filters.Transactions{AccountID: []int{id}, StartDate: startDate[0] + " 00:00:00", EndDate: endDate[0] + " 23:59:59"}

	dashboard, err := h.dashboardSvc.Get(ctx, f)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
