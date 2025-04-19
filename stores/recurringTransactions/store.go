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
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	res, err := ctx.SQL.ExecContext(ctx, createTransaction, recurringTransaction.UserID, recurringTransaction.Account.ID,
		recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.Category, recurringTransaction.Description,
		recurringTransaction.Frequency, recurringTransaction.CustomDays, recurringTransaction.StartDate, recurringTransaction.EndDate,
		recurringTransaction.LastRun, recurringTransaction.NextRun, createdAt)
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
		startDate            time.Time
		endDate              time.Time
		lastRun              time.Time
		nextRun              time.Time
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
	recurringTransaction.StartDate = startDate.Format("2006-01-02T15:04:05.000Z")
	recurringTransaction.EndDate = endDate.Format("2006-01-02T15:04:05.000Z")
	recurringTransaction.LastRun = lastRun.Format("2006-01-02T15:04:05.000Z")
	recurringTransaction.NextRun = nextRun.Format("2006-01-02T15:04:05.000Z")

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
			startDate            time.Time
			endDate              time.Time
			lastRun              time.Time
			nextRun              time.Time
		)

		err = rows.Scan(&recurringTransaction.ID, &recurringTransaction.UserID,
			&recurringTransaction.Account.ID, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Category, &recurringTransaction.Description,
			&recurringTransaction.Frequency, &recurringTransaction.CustomDays, &startDate, &endDate, &lastRun, &nextRun, &createdAt, &deletedAt, &recurringTransaction.Account.Name)
		if err != nil {
			return nil, err
		}

		recurringTransaction.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")
		recurringTransaction.StartDate = startDate.Format("2006-01-02T15:04:05.000Z")
		recurringTransaction.EndDate = endDate.Format("2006-01-02T15:04:05.000Z")
		recurringTransaction.LastRun = lastRun.Format("2006-01-02T15:04:05.000Z")
		recurringTransaction.NextRun = nextRun.Format("2006-01-02T15:04:05.000Z")

		if deletedAt.Valid {
			recurringTransaction.DeletedAt = deletedAt.String
		}

		allRecurringTransactions = append(allRecurringTransactions, &recurringTransaction)
	}

	return allRecurringTransactions, nil
}

func (s *recurringTransactionStore) Update(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) error {
	_, err := ctx.SQL.ExecContext(ctx, updateTransaction, recurringTransaction.Account.ID, recurringTransaction.Amount,
		recurringTransaction.Type, recurringTransaction.Category, recurringTransaction.Description, recurringTransaction.Frequency,
		recurringTransaction.CustomDays, recurringTransaction.StartDate, recurringTransaction.EndDate, recurringTransaction.LastRun,
		recurringTransaction.NextRun, recurringTransaction.ID)
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
