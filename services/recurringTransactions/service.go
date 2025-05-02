package recurringTransactions

import (
	"errors"
	"fmt"
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
	nextRun, err := calculateNextRun(recurringTransaction.StartDate, "", recurringTransaction.Frequency, recurringTransaction.CustomDays)
	if err != nil {
		return nil, err
	}

	recurringTransaction.NextRun = nextRun

	err = s.recurringTransactionStore.Create(ctx, recurringTransaction)
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

	oldTxn, err := s.GetByID(ctx, recurringTransaction.ID)
	if err != nil {
		return nil, err
	}

	recurringTransaction.UserID = userID
	recurringTransaction.StartDate, _ = convertToMySQLDate(recurringTransaction.StartDate)
	recurringTransaction.EndDate, _ = convertToMySQLDate(recurringTransaction.EndDate)
	recurringTransaction.NextRun, _ = convertToMySQLDate(recurringTransaction.NextRun)
	recurringTransaction.LastRun, _ = convertToMySQLDate(recurringTransaction.LastRun)

	if recurringTransaction.Frequency != oldTxn.Frequency ||
		(recurringTransaction.Frequency == models.CUSTOM && recurringTransaction.Frequency == oldTxn.Frequency && recurringTransaction.CustomDays != oldTxn.CustomDays) {
		nextRun, err := calculateNextRun(recurringTransaction.StartDate, "", recurringTransaction.Frequency, recurringTransaction.CustomDays)
		if err != nil {
			return nil, err
		}

		recurringTransaction.NextRun = nextRun
	}

	err = s.recurringTransactionStore.Update(ctx, recurringTransaction)
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

func calculateNextRun(startDateStr, lastRunStr string, freq models.Frequency, customDays int) (string, error) {
	var base time.Time
	var err error

	layout := "2006-01-02 15:04:05"
	now := time.Now()

	// 1. Determine base time
	if lastRunStr != "" {
		base, err = time.Parse(layout, lastRunStr)
		if err != nil {
			return "", err
		}
	} else if startDateStr != "" {
		base, err = time.Parse(layout, startDateStr)
		if err != nil {
			return "", err
		}

		if base.After(now) {
			return base.Format(layout), nil
		}
	} else {
		return "", errors.New("no base date provided")
	}

	// 2. Calculate next run
	for {
		switch freq {
		case models.DAILY:
			base = base.AddDate(0, 0, 1)
		case models.WEEKLY:
			base = base.AddDate(0, 0, 7)
		case models.MONTHLY:
			base = base.AddDate(0, 1, 0)
		case models.CUSTOM:
			if customDays <= 0 {
				return "", errors.New("invalid customDays: must be > 0")
			}
			base = base.AddDate(0, 0, customDays)
		default:
			return "", errors.New("unsupported frequency")
		}

		if base.After(now) {
			break
		}
	}

	return base.Format(layout), nil
}

func (s *recurringTransactionSvc) SkipNextRun(ctx *gofr.Context, id int) error {
	userID, _ := ctx.Value("userID").(int)

	transaction, err := s.recurringTransactionStore.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	layout := time.RFC3339
	next, err := time.Parse(layout, transaction.NextRun)
	if err != nil {
		return fmt.Errorf("invalid NextRun format: %v", err)
	}

	var newNextRun time.Time

	switch transaction.Frequency {
	case "DAILY":
		newNextRun = next.AddDate(0, 0, 1)
	case "WEEKLY":
		newNextRun = next.AddDate(0, 0, 7)
	case "MONTHLY":
		newNextRun = next.AddDate(0, 1, 0)
	case "CUSTOM":
		newNextRun = next.AddDate(0, 0, transaction.CustomDays)
	default:
		return errors.New("unsupported frequency")
	}

	transaction.NextRun = newNextRun.Format(layout)

	transaction.NextRun, _ = convertToMySQLDate(transaction.NextRun)
	transaction.StartDate, _ = convertToMySQLDate(transaction.StartDate)
	transaction.EndDate, _ = convertToMySQLDate(transaction.EndDate)

	err = s.recurringTransactionStore.Update(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}
