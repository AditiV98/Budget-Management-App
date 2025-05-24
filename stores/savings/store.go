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

	res, err := tx.ExecContext(ctx, createSavings, savings.UserID, savings.TransactionID,
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
		savings         models.Savings
		maturityDate    sql.NullString
		deletedAt       sql.NullString
		createdAt       time.Time
		withdrawnAmount sql.NullFloat64
	)

	err := ctx.SQL.QueryRowContext(ctx, getByIDSavings, id).Scan(&savings.ID, &savings.UserID, &savings.TransactionID,
		&savings.Category, &savings.Amount, &savings.CurrentValue, &savings.StartDate, &maturityDate,
		&createdAt, &deletedAt, &savings.Status, &withdrawnAmount)
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

	if withdrawnAmount.Valid {
		savings.WithdrawnAmount = withdrawnAmount.Float64
	}

	return &savings, nil
}

func (s *savingsStore) GetByTransactionID(ctx *gofr.Context, id int) (*models.Savings, error) {
	var (
		savings         models.Savings
		maturityDate    sql.NullString
		deletedAt       sql.NullString
		createdAt       time.Time
		withdrawnAmount sql.NullFloat64
	)

	err := ctx.SQL.QueryRowContext(ctx, getByTransactionIDSavings, id).Scan(&savings.ID, &savings.UserID, &savings.TransactionID,
		&savings.Category, &savings.Amount, &savings.CurrentValue, &savings.StartDate, &maturityDate,
		&createdAt, &deletedAt, &withdrawnAmount)
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

	if withdrawnAmount.Valid {
		savings.WithdrawnAmount = withdrawnAmount.Float64
	}

	return &savings, nil
}

func (s *savingsStore) GetAll(ctx *gofr.Context, f *filters.Savings) ([]*models.Savings, error) {
	var allSavings []*models.Savings

	clause, val := f.WhereClause()

	query := getAllSavings + clause + " ORDER BY start_date"

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
			savings         models.Savings
			maturityDate    sql.NullString
			deletedAt       sql.NullString
			createdAt       time.Time
			withdrawnAmount sql.NullFloat64
		)

		err = rows.Scan(&savings.ID, &savings.UserID, &savings.TransactionID, &savings.Category,
			&savings.Amount, &savings.CurrentValue, &savings.StartDate, &maturityDate,
			&createdAt, &deletedAt, &savings.Status, &withdrawnAmount)
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

		if withdrawnAmount.Valid {
			savings.WithdrawnAmount = withdrawnAmount.Float64
		}

		allSavings = append(allSavings, &savings)
	}

	return allSavings, nil
}

func (s *savingsStore) Update(ctx *gofr.Context, savings *models.Savings, tx *datasourceSQL.Tx) error {
	var maturityDate interface{}

	if savings.MaturityDate == "" {
		maturityDate = nil
	} else {
		maturityDate = savings.MaturityDate
	}

	_, err := tx.ExecContext(ctx, updateSavings, savings.CurrentValue, maturityDate, savings.Status, savings.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsStore) UpdateWIthTransactionID(ctx *gofr.Context, savings *models.Savings, tx *datasourceSQL.Tx) error {
	var startDate, maturityDate interface{}
	var withdrawnAmount sql.NullFloat64

	if savings.WithdrawnAmount != 0 {
		withdrawnAmount = sql.NullFloat64{Valid: true, Float64: savings.WithdrawnAmount}
	}

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

	_, err := tx.ExecContext(ctx, updateSavingsWithTransactionID, savings.Category, savings.Amount, savings.CurrentValue,
		startDate, maturityDate, withdrawnAmount, savings.TransactionID)
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

func (s *savingsStore) DeleteWithTx(ctx *gofr.Context, txnID int, tx *datasourceSQL.Tx) error {
	deletedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := tx.ExecContext(ctx, deleteSavingsByTransactionID, deletedAt, txnID)
	if err != nil {
		return err
	}

	return nil
}
