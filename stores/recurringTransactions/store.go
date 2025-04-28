package recurringTransactions

import (
	"database/sql"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type recurringTransactionStore struct{}

func New() stores.RecurringTransactions {
	return &recurringTransactionStore{}
}

func (s *recurringTransactionStore) Create(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) error {
	var nextRun interface{}

	if recurringTransaction.NextRun == "" {
		nextRun = nil
	} else {
		nextRun = recurringTransaction.NextRun
	}

	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	res, err := ctx.SQL.ExecContext(ctx, createTransaction, recurringTransaction.UserID, recurringTransaction.Account.ID,
		recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.Category, recurringTransaction.Description,
		recurringTransaction.Frequency, recurringTransaction.CustomDays, recurringTransaction.StartDate, recurringTransaction.EndDate,
		nextRun, createdAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return err
	}

	recurringTransaction.ID = int(id)

	return nil
}

func (s *recurringTransactionStore) GetByID(ctx *gofr.Context, id, userID int) (*models.RecurringTransaction, error) {
	var (
		recurringTransaction models.RecurringTransaction
		deletedAt            sql.NullString
		createdAt            time.Time
		startDate            sql.NullTime
		endDate              sql.NullTime
		lastRun              sql.NullTime
		nextRun              sql.NullTime
	)

	err := ctx.SQL.QueryRowContext(ctx, getByIDTransactions, id, userID).Scan(&recurringTransaction.ID, &recurringTransaction.UserID,
		&recurringTransaction.Account.ID, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Category, &recurringTransaction.Description,
		&recurringTransaction.Frequency, &recurringTransaction.CustomDays, &startDate, &endDate, &lastRun, &nextRun, &createdAt, &deletedAt, &recurringTransaction.Account.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	recurringTransaction.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

	if startDate.Valid {
		recurringTransaction.StartDate = startDate.Time.Format("2006-01-02T15:04:05.000Z")
	}

	if endDate.Valid {
		recurringTransaction.EndDate = endDate.Time.Format("2006-01-02T15:04:05.000Z")
	}

	if nextRun.Valid {
		recurringTransaction.NextRun = nextRun.Time.Format("2006-01-02T15:04:05.000Z")
	}

	if lastRun.Valid {
		recurringTransaction.LastRun = lastRun.Time.Format("2006-01-02T15:04:05.000Z")
	}

	if deletedAt.Valid {
		recurringTransaction.DeletedAt = deletedAt.String
	}

	return &recurringTransaction, nil
}

func (s *recurringTransactionStore) GetAll(ctx *gofr.Context, f *filters.RecurringTransactions) ([]*models.RecurringTransaction, error) {
	var allRecurringTransactions []*models.RecurringTransaction

	clause, val := f.WhereClause()

	query := getAllTransactions + clause + " ORDER BY next_run"

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
			recurringTransaction models.RecurringTransaction
			deletedAt            sql.NullString
			createdAt            time.Time
			startDate            sql.NullTime
			endDate              sql.NullTime
			lastRun              sql.NullTime
			nextRun              sql.NullTime
		)

		err = rows.Scan(&recurringTransaction.ID, &recurringTransaction.UserID,
			&recurringTransaction.Account.ID, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Category, &recurringTransaction.Description,
			&recurringTransaction.Frequency, &recurringTransaction.CustomDays, &startDate, &endDate, &lastRun, &nextRun, &createdAt, &deletedAt, &recurringTransaction.Account.Name)
		if err != nil {
			return nil, err
		}

		recurringTransaction.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")
		if startDate.Valid {
			recurringTransaction.StartDate = startDate.Time.Format("2006-01-02T15:04:05.000Z")
		}

		if endDate.Valid {
			recurringTransaction.EndDate = endDate.Time.Format("2006-01-02T15:04:05.000Z")
		}

		if nextRun.Valid {
			recurringTransaction.NextRun = nextRun.Time.Format("2006-01-02T15:04:05.000Z")
		}

		if lastRun.Valid {
			recurringTransaction.LastRun = lastRun.Time.Format("2006-01-02T15:04:05.000Z")
		}

		if deletedAt.Valid {
			recurringTransaction.DeletedAt = deletedAt.String
		}

		allRecurringTransactions = append(allRecurringTransactions, &recurringTransaction)
	}

	return allRecurringTransactions, nil
}

func (s *recurringTransactionStore) Update(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) error {
	var lastRun, nextRun interface{}
	if recurringTransaction.LastRun == "" {
		lastRun = nil
	} else {
		lastRun = recurringTransaction.LastRun
	}

	if recurringTransaction.NextRun == "" {
		nextRun = nil
	} else {
		nextRun = recurringTransaction.NextRun
	}

	_, err := ctx.SQL.ExecContext(ctx, updateTransaction, recurringTransaction.Account.ID, recurringTransaction.Amount,
		recurringTransaction.Type, recurringTransaction.Category, recurringTransaction.Description, recurringTransaction.Frequency,
		recurringTransaction.CustomDays, recurringTransaction.StartDate, recurringTransaction.EndDate, lastRun,
		nextRun, recurringTransaction.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *recurringTransactionStore) Delete(ctx *gofr.Context, id int) error {
	deletedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := ctx.SQL.ExecContext(ctx, deleteTransaction, deletedAt, id)
	if err != nil {
		return err
	}

	return nil
}
