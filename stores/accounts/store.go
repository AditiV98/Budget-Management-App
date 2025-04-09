package accounts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
	datasourceSQL "gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type accountStore struct{}

func New() stores.Account {
	return &accountStore{}
}

func (s *accountStore) Create(ctx *gofr.Context, account *models.Account) (int, error) {
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	expenseCategoriesJSON, err := json.Marshal(account.ExpenseCategories)
	if err != nil {
		return 0, err
	}

	savingCategoriesJSON, err := json.Marshal(account.SavingCategories)
	if err != nil {
		return 0, err
	}

	res, err := ctx.SQL.ExecContext(ctx, createAccount, account.UserID, account.Name, account.Type, account.Balance,
		account.Status, string(expenseCategoriesJSON), string(savingCategoriesJSON), createdAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return 0, err
	}

	return int(id), nil
}

func (s *accountStore) GetByID(ctx *gofr.Context, id, userID int) (*models.Account, error) {
	var (
		account               models.Account
		createdAt             time.Time
		deletedAt             sql.NullString
		expenseCategoriesJSON string
		savingCategoriesJSON  string
	)

	err := ctx.SQL.QueryRowContext(ctx, getByIDAccount, id, userID).Scan(&account.ID, &account.UserID, &account.Name,
		&account.Type, &account.Balance, &account.Status, &expenseCategoriesJSON, &savingCategoriesJSON, &createdAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	account.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

	if deletedAt.Valid {
		account.DeletedAt = deletedAt.String
	}

	err = json.Unmarshal([]byte(expenseCategoriesJSON), &account.ExpenseCategories)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(savingCategoriesJSON), &account.SavingCategories)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountStore) GetByIDForUpdate(ctx *gofr.Context, id, userID int, tx *datasourceSQL.Tx) (*models.Account, error) {
	var (
		account               models.Account
		createdAt             time.Time
		deletedAt             sql.NullString
		expenseCategoriesJSON string
		savingCategoriesJSON  string
	)

	err := tx.QueryRowContext(ctx, getByIDAccount+" FOR UPDATE;", id, userID).Scan(&account.ID, &account.UserID, &account.Name,
		&account.Type, &account.Balance, &account.Status, &expenseCategoriesJSON, &savingCategoriesJSON, &createdAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	account.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

	if deletedAt.Valid {
		account.DeletedAt = deletedAt.String
	}

	err = json.Unmarshal([]byte(expenseCategoriesJSON), &account.ExpenseCategories)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(savingCategoriesJSON), &account.SavingCategories)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountStore) GetAll(ctx *gofr.Context, f *filters.Account) ([]*models.Account, error) {
	var allAccounts []*models.Account

	clause, val := f.WhereClause()

	q := getAllAccount + clause

	rows, err := ctx.SQL.QueryContext(ctx, q, val...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var (
			account               models.Account
			createdAt             time.Time
			deletedAt             sql.NullString
			expenseCategoriesJSON string
			savingCategoriesJSON  string
		)

		err = rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Balance, &account.Status,
			&expenseCategoriesJSON, &savingCategoriesJSON, &createdAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		account.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

		if deletedAt.Valid {
			account.DeletedAt = deletedAt.String
		}

		err = json.Unmarshal([]byte(expenseCategoriesJSON), &account.ExpenseCategories)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(savingCategoriesJSON), &account.SavingCategories)
		if err != nil {
			return nil, err
		}

		allAccounts = append(allAccounts, &account)
	}

	return allAccounts, nil
}

func (s *accountStore) Update(ctx *gofr.Context, account *models.Account, tx *datasourceSQL.Tx) error {
	expenseCategoriesJSON, err := json.Marshal(account.ExpenseCategories)
	if err != nil {
		return err
	}

	savingCategoriesJSON, err := json.Marshal(account.SavingCategories)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, updateAccount, account.Name, account.Type, account.Balance, account.Status,
		string(expenseCategoriesJSON), string(savingCategoriesJSON), account.ID, account.UserID)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errors.New("account update failed")
	}

	return nil
}

func (s *accountStore) Delete(ctx *gofr.Context, id int) error {
	deletedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := ctx.SQL.ExecContext(ctx, deleteAccount, "INACTIVE", deletedAt, id)
	if err != nil {
		return err
	}

	return nil
}
