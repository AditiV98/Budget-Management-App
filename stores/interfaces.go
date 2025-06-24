package stores

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
)

type User interface {
	Create(ctx *gofr.Context, user *models.User) error
	GetByID(ctx *gofr.Context, id int) (*models.User, error)
	GetAll(ctx *gofr.Context, f *filters.User) ([]*models.User, error)
	Update(ctx *gofr.Context, user *models.User) error
	Delete(ctx *gofr.Context, id int) error
}

type Account interface {
	Create(ctx *gofr.Context, account *models.Account) (int, error)
	GetByID(ctx *gofr.Context, id, userID int) (*models.Account, error)
	GetAll(ctx *gofr.Context, f *filters.Account) ([]*models.Account, error)
	Update(ctx *gofr.Context, account *models.Account, tx *sql.Tx) error
	Delete(ctx *gofr.Context, id int) error
	GetByIDForUpdate(ctx *gofr.Context, id, userID int, tx *sql.Tx) (*models.Account, error)
}

type Transactions interface {
	Create(ctx *gofr.Context, transaction *models.Transaction, tx *sql.Tx) error
	GetAll(ctx *gofr.Context, f *filters.Transactions) ([]*models.Transaction, error)
	GetByID(ctx *gofr.Context, id, userID int) (*models.Transaction, error)
	Update(ctx *gofr.Context, transaction *models.Transaction, tx *sql.Tx) error
	Delete(ctx *gofr.Context, id int, tx *sql.Tx) error
}

type Savings interface {
	Create(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) error
	GetAll(ctx *gofr.Context, f *filters.Savings) ([]*models.Savings, error)
	GetByID(ctx *gofr.Context, id int) (*models.Savings, error)
	Update(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) error
	Delete(ctx *gofr.Context, id int) error
	UpdateWIthTransactionID(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) error
	GetByTransactionID(ctx *gofr.Context, id int) (*models.Savings, error)
	DeleteWithTx(ctx *gofr.Context, txnID int, tx *sql.Tx) error
}

type SavingsSource interface {
	Create(ctx *gofr.Context, savingsSource *models.SavingsSources) error
	GetByID(ctx *gofr.Context, id int) (*models.SavingsSources, error)
	Update(ctx *gofr.Context, savingsSource *models.SavingsSources) error
	Delete(ctx *gofr.Context, id int) error
}

type RecurringTransactions interface {
	Create(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) error
	GetAll(ctx *gofr.Context, f *filters.RecurringTransactions) ([]*models.RecurringTransaction, error)
	GetByID(ctx *gofr.Context, id, userID int) (*models.RecurringTransaction, error)
	Update(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) error
	Delete(ctx *gofr.Context, id int) error
}

type Configs interface {
	Create(ctx *gofr.Context, userID int) error
	Update(ctx *gofr.Context, config *models.Config) error
	Get(ctx *gofr.Context, userID int) (*models.Config, error)
}
