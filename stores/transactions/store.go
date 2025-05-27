package transactions

import (
	"database/sql"
	"encoding/json"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
	datasourceSQL "gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type transactionStore struct{}

func New() stores.Transactions {
	return &transactionStore{}
}

func (s *transactionStore) Create(ctx *gofr.Context, transaction *models.Transaction, tx *datasourceSQL.Tx) error {
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	var withdrawFrom sql.NullInt64

	if transaction.WithdrawFrom != 0 {
		withdrawFrom = sql.NullInt64{Valid: true, Int64: transaction.WithdrawFrom}
	}

	jsonData, _ := json.Marshal(transaction.MetaData)

	res, err := tx.ExecContext(ctx, createTransaction, transaction.UserID, transaction.Account.ID, transaction.Amount,
		transaction.Type, transaction.Category, transaction.Description, transaction.TransactionDate, createdAt,
		withdrawFrom, string(jsonData))
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return err
	}

	transaction.ID = int(id)

	return nil
}

func (s *transactionStore) GetByID(ctx *gofr.Context, id, userID int) (*models.Transaction, error) {
	var (
		transaction     models.Transaction
		deletedAt       sql.NullString
		createdAt       time.Time
		transactionDate time.Time
		withdrawFrom    sql.NullInt64
		metaData        sql.NullString
	)

	err := ctx.SQL.QueryRowContext(ctx, getByIDTransactions, id, userID).Scan(&transaction.ID, &transaction.UserID,
		&transaction.Account.ID, &transaction.Amount, &transaction.Type, &transaction.Category, &transaction.Description,
		&transactionDate, &createdAt, &deletedAt, &withdrawFrom, &metaData, &transaction.Account.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	transaction.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")
	transaction.TransactionDate = transactionDate.Format("2006-01-02T15:04:05.000Z")

	if deletedAt.Valid {
		transaction.DeletedAt = deletedAt.String
	}

	if withdrawFrom.Valid {
		transaction.WithdrawFrom = withdrawFrom.Int64
	}

	if metaData.Valid {
		_ = json.Unmarshal([]byte(metaData.String), &transaction.MetaData)
	}

	return &transaction, nil
}

func (s *transactionStore) GetAll(ctx *gofr.Context, f *filters.Transactions) ([]*models.Transaction, error) {
	var allTransactions []*models.Transaction

	clause, val := f.WhereClause()

	query := getAllTransactions + clause + " ORDER BY transaction_date"

	if f.SortBy == "DESC" {
		query += " DESC"
	}

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
			transaction         models.Transaction
			deletedAt           sql.NullString
			createdAt           time.Time
			transactionDate     time.Time
			withdrawFrom        sql.NullInt64
			metaData            sql.NullString
			savingID            sql.NullInt64
			status              sql.NullString
			savingTransactionID sql.NullInt64
			savingAmount        sql.NullFloat64
			savingCurrentValue  sql.NullFloat64
		)

		err = rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Account.ID, &transaction.Amount, &transaction.Type,
			&transaction.Category, &transaction.Description, &transactionDate, &createdAt, &deletedAt, &withdrawFrom,
			&metaData, &transaction.Account.Name, &savingID, &status, &savingTransactionID, &savingAmount, &savingCurrentValue)
		if err != nil {
			return nil, err
		}

		transaction.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")
		transaction.TransactionDate = transactionDate.Format("2006-01-02T15:04:05.000Z")

		if deletedAt.Valid {
			transaction.DeletedAt = deletedAt.String
		}

		if withdrawFrom.Valid {
			transaction.WithdrawFrom = withdrawFrom.Int64
		}

		if metaData.Valid {
			_ = json.Unmarshal([]byte(metaData.String), &transaction.MetaData)
		}

		if status.Valid {
			transaction.Saving.Status = status.String
		}

		if savingID.Valid {
			transaction.Saving.ID = int(savingID.Int64)
		}

		if savingTransactionID.Valid {
			transaction.Saving.TransactionID = int(savingTransactionID.Int64)
		}

		if savingAmount.Valid {
			transaction.Saving.Amount = savingAmount.Float64
		}

		if savingCurrentValue.Valid {
			transaction.Saving.CurrentValue = savingCurrentValue.Float64
		}

		allTransactions = append(allTransactions, &transaction)
	}

	return allTransactions, nil
}

func (s *transactionStore) Update(ctx *gofr.Context, transaction *models.Transaction, tx *datasourceSQL.Tx) error {
	jsonData, _ := json.Marshal(transaction.MetaData)

	_, err := tx.ExecContext(ctx, updateTransaction, transaction.Account.ID, transaction.Amount, transaction.Type,
		transaction.Category, transaction.Description, transaction.TransactionDate, string(jsonData), transaction.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *transactionStore) Delete(ctx *gofr.Context, id int, tx *datasourceSQL.Tx) error {
	deletedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := tx.ExecContext(ctx, deleteTransaction, deletedAt, id)
	if err != nil {
		return err
	}

	return nil
}
