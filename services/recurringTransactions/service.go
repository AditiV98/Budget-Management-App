package recurringTransactions

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
	"time"
)

type recurringTransactionSvc struct {
	recurringTransactionStore stores.RecurringTransactions
	userSvc                   services.User
}

func New(recurringTransactionStore stores.RecurringTransactions, userSvc services.User) services.RecurringTransactions {
	return &recurringTransactionSvc{
		recurringTransactionStore: recurringTransactionStore,
		userSvc:                   userSvc,
	}
}

func (s *recurringTransactionSvc) Create(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) (*models.RecurringTransaction, error) {
	userID, _ := ctx.Value("userID").(int)

	recurringTransaction.UserID = userID
	recurringTransaction.StartDate, _ = convertToMySQLDate(recurringTransaction.StartDate)
	recurringTransaction.EndDate, _ = convertToMySQLDate(recurringTransaction.EndDate)
	recurringTransaction.NextRun, _ = convertToMySQLDate(recurringTransaction.NextRun)

	err := s.recurringTransactionStore.Create(ctx, recurringTransaction)
	if err != nil {
		return nil, err
	}

	newTransaction, err := s.GetByID(ctx, recurringTransaction.ID)
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (s *recurringTransactionSvc) GetByID(ctx *gofr.Context, id int) (*models.RecurringTransaction, error) {
	userID, _ := ctx.Value("userID").(int)

	transaction, err := s.recurringTransactionStore.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *recurringTransactionSvc) GetAll(ctx *gofr.Context, f *filters.RecurringTransactions) ([]*models.RecurringTransaction, error) {
	userID, _ := ctx.Value("userID").(int)

	f.UserID = userID

	accounts, err := s.recurringTransactionStore.GetAll(ctx, f)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *recurringTransactionSvc) Update(ctx *gofr.Context, recurringTransaction *models.RecurringTransaction) (*models.RecurringTransaction, error) {
	userID, _ := ctx.Value("userID").(int)

	recurringTransaction.UserID = userID
	recurringTransaction.StartDate, _ = convertToMySQLDate(recurringTransaction.StartDate)
	recurringTransaction.EndDate, _ = convertToMySQLDate(recurringTransaction.EndDate)
	recurringTransaction.NextRun, _ = convertToMySQLDate(recurringTransaction.NextRun)
	recurringTransaction.LastRun, _ = convertToMySQLDate(recurringTransaction.LastRun)

	err := s.recurringTransactionStore.Update(ctx, recurringTransaction)
	if err != nil {
		return nil, err
	}

	updatedTransaction, err := s.GetByID(ctx, recurringTransaction.ID)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (s *recurringTransactionSvc) Delete(ctx *gofr.Context, id int) error {
	userID, _ := ctx.Value("userID").(int)

	_, err := s.recurringTransactionStore.GetByID(ctx, id, userID)
	if err != nil {
		return errors.New("unauthorised")
	}

	err = s.recurringTransactionStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func convertToMySQLDate(isoDate string) (string, error) {
	if isoDate != "" {
		t, err := time.Parse(time.RFC3339, isoDate) // Parses "2025-03-20T07:49:00.000Z"
		if err != nil {
			return "", err
		}

		return t.Format("2006-01-02 15:04:05"), nil // Converts to "2025-03-20 07:49:00"
	}

	return "", nil
}
