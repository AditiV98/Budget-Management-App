package dashboard

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/handler"
	"moneyManagement/services"
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

	dashboard, err := h.dashboardSvc.Get(ctx, startDate[0]+" 00:00:00", endDate[0]+" 23:59:59")
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
