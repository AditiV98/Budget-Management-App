package transactions

import (
	"database/sql"
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
	"time"
)

type transactionSvc struct {
	transactionStore stores.Transactions
	accountSvc       services.Account
	savingsSvc       services.Savings
	userSvc          services.User
}

func New(transactionStore stores.Transactions, accountSvc services.Account, savingsSvc services.Savings, userSvc services.User) services.Transactions {
	return &transactionSvc{
		transactionStore: transactionStore,
		accountSvc:       accountSvc,
		savingsSvc:       savingsSvc,
		userSvc:          userSvc,
	}
}

func (s *transactionSvc) Create(ctx *gofr.Context, transaction *models.Transaction) (*models.Transaction, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	// Ensure rollback only if an error occurs
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	transaction.UserID = userID
	transaction.TransactionDate, _ = convertToMySQLDate(transaction.TransactionDate)

	// 1️⃣ Lock the Account Row First using FOR UPDATE
	account, err := s.accountSvc.GetByIDForUpdate(ctx, transaction.Account.ID, userID, tx)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Update Account Balance Immediately
	if transaction.Type == "INCOME" {
		a := account.Balance + transaction.Amount
		ctx.Logger.Warnf("transaction.Type :%v,transaction.Amount:%v,previous account balance :%v, new account balance: %v", "INCOME", transaction.Amount, account.Balance, a)
		account.Balance = a
	} else if transaction.Type == "SAVINGS" {
		b := account.Balance - transaction.Amount
		ctx.Logger.Warnf("transaction.Type :%v,transaction.Amount:%v,previous account balance :%v, new account balance: %v", "SAVINGS", transaction.Amount, account.Balance, b)
		account.Balance = b
	} else if transaction.Type == "EXPENSE" {
		c := account.Balance - transaction.Amount
		ctx.Logger.Warnf("transaction.Type :%v,transaction.Amount:%v,previous account balance :%v, new account balance: %v", "EXPENSE", transaction.Amount, account.Balance, c)
		account.Balance = c
	}

	// ✅ Update account balance before inserting transaction
	_, err = s.accountSvc.UpdateWithTx(ctx, account, tx)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Insert the Transaction Record
	err = s.transactionStore.Create(ctx, transaction, tx)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Insert into Savings only if it's a "SAVINGS" transaction
	if transaction.Type == "SAVINGS" {
		savings := &models.Savings{
			UserID: transaction.UserID, Amount: transaction.Amount, Type: transaction.Category,
			StartDate: transaction.TransactionDate, TransactionID: transaction.ID,
		}

		_, err = s.savingsSvc.CreateWithTx(ctx, savings, tx)
		if err != nil {
			return nil, err
		}
	}

	// 5️⃣ Commit Transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// 6️⃣ Fetch and return the newly created transaction after commit
	newTransaction, err := s.GetByID(ctx, transaction.ID)
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (s *transactionSvc) GetByID(ctx *gofr.Context, id int) (*models.Transaction, error) {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	transaction, err := s.transactionStore.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionSvc) GetAll(ctx *gofr.Context, f *filters.Transactions) ([]*models.Transaction, error) {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	f.UserID = userID

	allTransactions, err := s.transactionStore.GetAll(ctx, f)
	if err != nil {
		return nil, err
	}

	return allTransactions, nil
}

func (s *transactionSvc) Update(ctx *gofr.Context, transaction *models.Transaction) (*models.Transaction, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	transaction.UserID = userID

	transaction.TransactionDate, _ = convertToMySQLDate(transaction.TransactionDate)

	// Fetch the original transaction to compare values
	originalTransaction, err := s.GetByID(ctx, transaction.ID)
	if err != nil {
		return nil, err
	}

	// Get the associated account
	account, err := s.accountSvc.GetByID(ctx, transaction.Account.ID)
	if err != nil {
		return nil, err
	}

	// Reverse the effect of the original transaction
	if originalTransaction.Type == "INCOME" {
		account.Balance -= originalTransaction.Amount
	} else if originalTransaction.Type == "SAVINGS" {
		account.Balance += originalTransaction.Amount
	} else if originalTransaction.Type == "EXPENSE" {
		account.Balance += originalTransaction.Amount
	}

	// Apply the effect of the updated transaction
	if transaction.Type == "INCOME" {
		account.Balance += transaction.Amount
	} else if transaction.Type == "SAVINGS" {
		savings := &models.Savings{
			UserID:        transaction.UserID,
			Amount:        transaction.Amount,
			Type:          transaction.Category,
			StartDate:     transaction.TransactionDate,
			TransactionID: transaction.ID,
		}
		// Check if savings entry already exists for this transaction
		_, err = s.savingsSvc.GetByTransactionID(ctx, transaction.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Create new savings record if not found
				_, err = s.savingsSvc.CreateWithTx(ctx, savings, tx)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// Update existing if it exists
			_, err = s.savingsSvc.UpdateWithTx(ctx, savings, true, tx)
			if err != nil {
				return nil, err
			}
		}

		account.Balance -= transaction.Amount
	} else if transaction.Type == "EXPENSE" {
		saving, er := s.savingsSvc.GetByTransactionID(ctx, transaction.ID)
		if er == nil {
			err = s.savingsSvc.Delete(ctx, saving.ID)
			if err != nil {
				return nil, err
			}
		}
		account.Balance -= transaction.Amount
	}

	// Update transaction record
	err = s.transactionStore.Update(ctx, transaction, tx)
	if err != nil {
		return nil, err
	}

	// Update account balance
	_, err = s.accountSvc.UpdateWithTx(ctx, account, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Retrieve and return the updated transaction
	updatedTransaction, err := s.GetByID(ctx, transaction.ID)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (s *transactionSvc) Delete(ctx *gofr.Context, id int) error {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	originalTransaction, err := s.GetByID(ctx, id)
	if err != nil {
		return errors.New("unauthorised")
	}

	account, err := s.accountSvc.GetByID(ctx, originalTransaction.Account.ID)
	if err != nil {
		return err
	}

	if originalTransaction.Type == "EXPENSE" || originalTransaction.Type == "SAVINGS" {
		account.Balance += originalTransaction.Amount
	} else if originalTransaction.Type == "INCOME" {
		account.Balance -= originalTransaction.Amount
	}

	err = s.transactionStore.Delete(ctx, id, tx)
	if err != nil {
		return err
	}

	_, err = s.accountSvc.UpdateWithTx(ctx, account, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
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
