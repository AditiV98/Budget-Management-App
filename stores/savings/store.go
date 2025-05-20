package savings

import (
	"database/sql"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
	datasourceSQL "gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type savingsStore struct{}

func New() stores.Savings {
	return &savingsStore{}
}

func (s *savingsStore) Create(ctx *gofr.Context, savings *models.Savings, tx *datasourceSQL.Tx) error {
	var startDate, maturityDate interface{}

	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	if savings.StartDate == "" {
		startDate = nil
	} else {
		startDate = savings.StartDate
	}

	if savings.MaturityDate == "" {
		maturityDate = nil
	} else {
		maturityDate = savings.MaturityDate
	}

	res, err := tx.ExecContext(ctx, createSavings, savings.UserID, savings.TransactionID, savings.Type,
		savings.Category, savings.Amount, savings.CurrentValue, startDate, maturityDate, createdAt, savings.Status)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return err
	}

	savings.ID = int(id)

	return nil
}

func (s *savingsStore) GetByID(ctx *gofr.Context, id int) (*models.Savings, error) {
	var (
		savings      models.Savings
		maturityDate sql.NullString
		deletedAt    sql.NullString
		createdAt    time.Time
	)

	err := ctx.SQL.QueryRowContext(ctx, getByIDSavings, id).Scan(&savings.ID, &savings.UserID, &savings.TransactionID,
		&savings.Type, &savings.Category, &savings.Amount, &savings.CurrentValue, &savings.StartDate, &maturityDate,
		&createdAt, &deletedAt, &savings.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	savings.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

	if deletedAt.Valid {
		savings.DeletedAt = deletedAt.String
	}

	if maturityDate.Valid {
		savings.MaturityDate = maturityDate.String
	}

	return &savings, nil
}

func (s *savingsStore) GetByTransactionID(ctx *gofr.Context, id int) (*models.Savings, error) {
	var (
		savings      models.Savings
		maturityDate sql.NullString
		deletedAt    sql.NullString
		createdAt    time.Time
	)

	err := ctx.SQL.QueryRowContext(ctx, getByTransactionIDSavings, id).Scan(&savings.ID, &savings.UserID, &savings.TransactionID,
		&savings.Type, &savings.Category, &savings.Amount, &savings.CurrentValue, &savings.StartDate, &maturityDate,
		&createdAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	savings.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

	if deletedAt.Valid {
		savings.DeletedAt = deletedAt.String
	}

	if maturityDate.Valid {
		savings.MaturityDate = maturityDate.String
	}

	return &savings, nil
}

func (s *savingsStore) GetAll(ctx *gofr.Context, f *filters.Savings) ([]*models.Savings, error) {
	var allSavings []*models.Savings

	clause, val := f.WhereClause()

	query := getAllSavings + clause

	rows, err := ctx.SQL.QueryContext(ctx, query, val...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var (
			savings      models.Savings
			maturityDate sql.NullString
			deletedAt    sql.NullString
			createdAt    time.Time
		)

		err = rows.Scan(&savings.ID, &savings.UserID, &savings.TransactionID, &savings.Type, &savings.Category,
			&savings.Amount, &savings.CurrentValue, &savings.StartDate, &maturityDate,
			&createdAt, &deletedAt, &savings.Status)
		if err != nil {
			return nil, err
		}

		savings.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

		if deletedAt.Valid {
			savings.DeletedAt = deletedAt.String
		}

		if maturityDate.Valid {
			savings.MaturityDate = maturityDate.String
		}

		allSavings = append(allSavings, &savings)
	}

	return allSavings, nil
}

func (s *savingsStore) Update(ctx *gofr.Context, savings *models.Savings, tx *datasourceSQL.Tx) error {
	var startDate, maturityDate interface{}

	if savings.StartDate == "" {
		startDate = nil
	} else {
		startDate = savings.StartDate
	}

	if savings.MaturityDate == "" {
		maturityDate = nil
	} else {
		maturityDate = savings.MaturityDate
	}

	_, err := tx.ExecContext(ctx, updateSavings, savings.Type, savings.Category, savings.Amount, savings.CurrentValue,
		startDate, maturityDate, savings.Status, savings.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsStore) UpdateWIthTransactionID(ctx *gofr.Context, savings *models.Savings, tx *datasourceSQL.Tx) error {
	var startDate, maturityDate interface{}

	if savings.StartDate == "" {
		startDate = nil
	} else {
		startDate = savings.StartDate
	}

	if savings.MaturityDate == "" {
		maturityDate = nil
	} else {
		maturityDate = savings.MaturityDate
	}

	_, err := tx.ExecContext(ctx, updateSavingsWithTransactionID, savings.Type, savings.Category, savings.Amount, savings.CurrentValue,
		startDate, maturityDate, savings.TransactionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsStore) Delete(ctx *gofr.Context, id int) error {
	deletedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := ctx.SQL.ExecContext(ctx, deleteSavings, deletedAt, id)
	if err != nil {
		return err
	}

	return nil
}
