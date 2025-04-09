package savings

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
)

type savingsSvc struct {
	savingsStore stores.Savings
}

func New(savingsStore stores.Savings) services.Savings {
	return &savingsSvc{
		savingsStore: savingsStore,
	}
}

func (s *savingsSvc) Create(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) error {
	err := s.savingsStore.Create(ctx, savings, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsSvc) GetByID(ctx *gofr.Context, id int) (*models.Savings, error) {
	savings, err := s.savingsStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return savings, nil
}

func (s *savingsSvc) GetByTransactionID(ctx *gofr.Context, id int) (*models.Savings, error) {
	savings, err := s.savingsStore.GetByTransactionID(ctx, id)
	if err != nil {
		return nil, err
	}

	return savings, nil
}

func (s *savingsSvc) GetAll(ctx *gofr.Context) ([]*models.Savings, error) {
	allSavings, err := s.savingsStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return allSavings, nil
}

func (s *savingsSvc) Update(ctx *gofr.Context, savings *models.Savings, IsTransactionID bool, tx *sql.Tx) error {
	if IsTransactionID {
		err := s.savingsStore.UpdateWIthTransactionID(ctx, savings, tx)
		if err != nil {
			return err
		}
	} else {
		err := s.savingsStore.Update(ctx, savings, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *savingsSvc) Delete(ctx *gofr.Context, id int) error {
	err := s.savingsStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
