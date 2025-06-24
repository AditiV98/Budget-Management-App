package configs

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
)

type configsSvc struct {
	configsStore stores.Configs
}

func New(configsStore stores.Configs) services.Configs {
	return &configsSvc{
		configsStore: configsStore,
	}
}

func (s *configsSvc) Create(ctx *gofr.Context) error {
	userID, _ := ctx.Value("userID").(int)

	err := s.configsStore.Create(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *configsSvc) Update(ctx *gofr.Context, config *models.Config) (*models.Config, error) {
	userID, _ := ctx.Value("userID").(int)

	currentConfigs, err := s.configsStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if currentConfigs.IsAutoRead != config.IsAutoRead {
		currentConfigs.IsAutoRead = config.IsAutoRead

		err = s.configsStore.Update(ctx, currentConfigs)
		if err != nil {
			return nil, err
		}
	}

	return currentConfigs, nil
}

func (s *configsSvc) Get(ctx *gofr.Context) (*models.Config, error) {
	userID, _ := ctx.Value("userID").(int)

	configs, err := s.configsStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
