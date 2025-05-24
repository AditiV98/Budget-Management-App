package savings

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
	"time"
)

type savingsSvc struct {
	savingsStore     stores.Savings
	transactionStore stores.Transactions
}

func New(savingsStore stores.Savings, transactionStore stores.Transactions) services.Savings {
	return &savingsSvc{
		savingsStore:     savingsStore,
		transactionStore: transactionStore,
	}
}

func (s *savingsSvc) Create(ctx *gofr.Context, savings *models.Savings) (*models.Savings, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	userID, _ := ctx.Value("userID").(int)

	savings.UserID = userID

	err = s.savingsStore.Create(ctx, savings, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	newSaving, err := s.GetByID(ctx, savings.ID)
	if err != nil {
		return nil, err
	}

	return newSaving, nil
}

func (s *savingsSvc) CreateWithTx(ctx *gofr.Context, savings *models.Savings, tx *sql.Tx) error {
	err := s.savingsStore.Create(ctx, savings, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsSvc) GetByID(ctx *gofr.Context, id int) (*models.Savings, error) {
	savings, err := s.savingsStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return savings, nil
}

func (s *savingsSvc) GetByTransactionID(ctx *gofr.Context, id int) (*models.Savings, error) {
	savings, err := s.savingsStore.GetByTransactionID(ctx, id)
	if err != nil {
		return nil, err
	}

	return savings, nil
}

func (s *savingsSvc) GetAll(ctx *gofr.Context, f *filters.Savings) ([]*models.Savings, error) {
	userID, _ := ctx.Value("userID").(int)

	f.UserID = userID

	allSavings, err := s.savingsStore.GetAll(ctx, f)
	if err != nil {
		return nil, err
	}

	transactionIDs := make([]int, 0)

	for i := range allSavings {
		transactionIDs = append(transactionIDs, allSavings[i].TransactionID)
	}

	transactions, err := s.transactionStore.GetAll(ctx, &filters.Transactions{ID: transactionIDs})
	if err != nil {
		return nil, err
	}

	m := make(map[int]*models.Transaction)

	for j := range transactions {
		m[transactions[j].ID] = transactions[j]
	}

	for k := range allSavings {
		if v, ok := m[allSavings[k].TransactionID]; ok {
			allSavings[k].Account = v.Account
		}
	}

	return allSavings, nil
}

func (s *savingsSvc) Update(ctx *gofr.Context, savings *models.Savings) (*models.Savings, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	userID, _ := ctx.Value("userID").(int)

	savings.UserID = userID
	savings.MaturityDate, _ = convertToMySQLDate(savings.MaturityDate)

	err = s.savingsStore.Update(ctx, savings, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedSaving, err := s.GetByID(ctx, savings.ID)
	if err != nil {
		return nil, err
	}

	txn, err := s.transactionStore.GetByID(ctx, updatedSaving.TransactionID, userID)
	if err != nil {
		return nil, err
	}

	updatedSaving.Account = txn.Account

	return updatedSaving, nil
}

func (s *savingsSvc) UpdateWithTx(ctx *gofr.Context, savings *models.Savings, IsTransactionID bool, tx *sql.Tx) error {
	if IsTransactionID {
		err := s.savingsStore.UpdateWIthTransactionID(ctx, savings, tx)
		if err != nil {
			return err
		}
	} else {
		err := s.savingsStore.Update(ctx, savings, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *savingsSvc) Delete(ctx *gofr.Context, id int) error {
	err := s.savingsStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *savingsSvc) DeleteWithTx(ctx *gofr.Context, txnID int, tx *sql.Tx) error {
	err := s.savingsStore.DeleteWithTx(ctx, txnID, tx)
	if err != nil {
		return err
	}

	return nil
}

func convertToMySQLDate(isoDate string) (string, error) {
	t, err := time.Parse(time.RFC3339, isoDate) // Parses "2025-03-20T07:49:00.000Z"
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02 15:04:05"), nil // Converts to "2025-03-20 07:49:00"
}
