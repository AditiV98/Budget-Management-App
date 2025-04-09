package savingsSource

import (
	"database/sql"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type savingsSourceStore struct{}

func New() stores.SavingsSource {
	return &savingsSourceStore{}
}

func (s *savingsSourceStore) Create(ctx *gofr.Context, savingsSource *models.SavingsSources) error {
	res, err := ctx.SQL.ExecContext(ctx, createSavingsSource, savingsSource.SavingID, savingsSource.TransactionID, savingsSource.Amount)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return err
	}

	savingsSource.ID = int(id)

	return nil
}

func (s *savingsSourceStore) GetByID(ctx *gofr.Context, id int) (*models.SavingsSources, error) {
	var savingsSource models.SavingsSources

	err := ctx.SQL.QueryRowContext(ctx, getByIDSavingsSource, id).Scan(&savingsSource.ID, &savingsSource.SavingID, &savingsSource.TransactionID, &savingsSource.Amount, &savingsSource.CreatedAt, &savingsSource.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	return &savingsSource, nil
}

func (s *savingsSourceStore) Update(ctx *gofr.Context, savingsSource *models.SavingsSources) error {
	_, err := ctx.SQL.ExecContext(ctx, updateSavingsSource, savingsSource.SavingID, savingsSource.TransactionID, savingsSource.Amount, savingsSource.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsSourceStore) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.SQL.ExecContext(ctx, deleteSavingsSource, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
