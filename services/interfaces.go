package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
)

type User interface {
	Create(ctx *gofr.Context, user *models.User) (*models.User, error)
	GetByID(ctx *gofr.Context, id int) (*models.User, error)
	GetAll(ctx *gofr.Context, f *filters.User) ([]*models.User, error)
	Update(ctx *gofr.Context, user *models.User) (*models.User, error)
	Delete(ctx *gofr.Context, id int) error
	AuthAdaptor(ctx *gofr.Context, claims *models.GoogleClaims) error
}

type Account interface {
	Create(ctx *gofr.Context, account *models.Account) (*models.Account, error)
	GetByID(ctx *gofr.Context, id int) (*models.Account, error)
	GetAll(ctx *gofr.Context, f *filters.Account) ([]*models.Account, error)
	Update(ctx *gofr.Context, account *models.Account) (*models.Account, error)
	UpdateWithTx(ctx *gofr.Context, account *models.Account, tx *sql.Tx) (*models.Account, error)
	Delete(ctx *gofr.Context, id int) error
	GetByIDForUpdate(ctx *gofr.Context, id, userID int, tx *sql.Tx) (*models.Account, error)
}

type Transactions interface {
	Create(ctx *gofr.Context, transaction *models.Transaction) (*models.Transaction, error)
	GetAll(ctx *gofr.Context, f *filters.Transactions) ([]*models.Transaction, error)
	GetByID(ctx *gofr.Context, id int) (*models.Transaction, error)
	Update(ctx *gofr.Context, transaction *models.Transaction) (*models.Transaction, error)
	Delete(ctx *gofr.Context, id int) error
}

type Savings interface {
	Create(ctx *gofr.Context, savings *models.Savings) (*models.Savings, error)
	CreateWithTx(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) error
	GetAll(ctx *gofr.Context, f *filters.Savings) ([]*models.Savings, error)
	GetByID(ctx *gofr.Context, id int) (*models.Savings, error)
	Update(ctx *gofr.Context, savings *models.Savings) (*models.Savings, error)
	UpdateWithTx(ctx *gofr.Context, savings *models.Savings, IsTransactionID bool, tx *sql.Tx) error
	Delete(ctx *gofr.Context, id int) error
	GetByTransactionID(ctx *gofr.Context, id int) (*models.Savings, error)
	DeleteWithTx(ctx *gofr.Context, txnID int, tx *sql.Tx) error
}

type Dashboard interface {
	Get(ctx *gofr.Context, f *filters.Transactions) (models.Dashboard, error)
}

type Auth interface {
	GenerateGoogleToken(ctx *gofr.Context, code string) (map[string]interface{}, error)
	GenerateRefreshToken(claims *models.GoogleClaims) (string, error)
	GenerateAccessToken(claims *models.GoogleClaims) (string, error)
	ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error)
	VerifyGoogleIDToken(ctx context.Context, idToken string) (*models.GoogleClaims, error)
}

type Validator interface {
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
}

type RecurringTransactions interface {
	Create(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) (*models.RecurringTransaction, error)
	GetAll(ctx *gofr.Context, f *filters.RecurringTransactions) ([]*models.RecurringTransaction, error)
	GetByID(ctx *gofr.Context, id int) (*models.RecurringTransaction, error)
	Update(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) (*models.RecurringTransaction, error)
	Delete(ctx *gofr.Context, id int) error
	SkipNextRun(ctx *gofr.Context, id int) error
}
