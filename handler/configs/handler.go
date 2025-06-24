package configs

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
)

type configsHandler struct {
	configsSvc services.Configs
}

func New(configsSvc services.Configs) handler.Configs {
	return &configsHandler{configsSvc: configsSvc}
}

func (h configsHandler) Create(ctx *gofr.Context) (interface{}, error) {
	err := h.configsSvc.Create(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h configsHandler) Get(ctx *gofr.Context) (interface{}, error) {
	config, err := h.configsSvc.Get(ctx)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (h configsHandler) Update(ctx *gofr.Context) (interface{}, error) {
	var config models.Config

	err := ctx.Bind(&config)
	if err != nil {
		return nil, errors.New("bind error")
	}

	updatedConfig, err := h.configsSvc.Update(ctx, &config)
	if err != nil {
		return nil, err
	}

	return updatedConfig, nil
}
