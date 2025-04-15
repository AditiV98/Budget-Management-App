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

func (s *savingsSvc) Create(ctx *gofr.Context, savings *models.Savings) (*models.Savings, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = s.savingsStore.Create(ctx, savings, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	newSaving, err := s.GetByID(ctx, savings.ID)
	if err != nil {
		return nil, err
	}

	return newSaving, nil
}

func (s *savingsSvc) CreateWithTx(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) (*models.Savings, error) {
	err := s.savingsStore.Create(ctx, savings, tx)
	if err != nil {
		return nil, err
	}

	newSaving, err := s.GetByID(ctx, savings.ID)
	if err != nil {
		return nil, err
	}

	return newSaving, nil
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

func (s *savingsSvc) Update(ctx *gofr.Context, savings *models.Savings) (*models.Savings, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = s.savingsStore.Update(ctx, savings, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedSaving, err := s.GetByID(ctx, savings.ID)
	if err != nil {
		return nil, err
	}

	return updatedSaving, nil
}

func (s *savingsSvc) UpdateWithTx(ctx *gofr.Context, savings *models.Savings, IsTransactionID bool, tx *sql.Tx) (*models.Savings, error) {
	if IsTransactionID {
		err := s.savingsStore.UpdateWIthTransactionID(ctx, savings, tx)
		if err != nil {
			return nil, err
		}
	} else {
		err := s.savingsStore.Update(ctx, savings, tx)
		if err != nil {
			return nil, err
		}
	}

	updatedSaving, err := s.GetByID(ctx, savings.ID)
	if err != nil {
		return nil, err
	}

	return updatedSaving, nil
}

func (s *savingsSvc) Delete(ctx *gofr.Context, id int) error {
	err := s.savingsStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
