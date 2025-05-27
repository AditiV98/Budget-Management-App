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

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	userID, _ := ctx.Value("userID").(int)

	transaction.UserID = userID
	transaction.TransactionDate, _ = convertToMySQLDate(transaction.TransactionDate)

	account, err := s.accountSvc.GetByIDForUpdate(ctx, transaction.Account.ID, userID, tx)
	if err != nil {
		return nil, err
	}

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
	} else if transaction.Type == "WITHDRAW" {
		d := account.Balance + transaction.Amount
		ctx.Logger.Warnf("transaction.Type :%v,transaction.Amount:%v,previous account balance :%v, new account balance: %v", "WITHDRAW", transaction.Amount, account.Balance, d)
		account.Balance = d
	} else if transaction.Type == "SELF TRANSFER" {
		d := account.Balance - transaction.Amount
		ctx.Logger.Warnf("transaction.Type :%v,transaction.Amount:%v,previous account balance :%v, new account balance: %v", "WITHDRAW", transaction.Amount, account.Balance, d)
		account.Balance = d

		accountToTransfer, er := s.accountSvc.GetByID(ctx, transaction.MetaData.TransferTo)
		if er != nil {
			return nil, er
		}

		accountToTransfer.Balance += transaction.Amount

		_, err = s.accountSvc.UpdateWithTx(ctx, accountToTransfer, tx)
		if err != nil {
			return nil, err
		}

		transaction.MetaData.TransferFrom = account.ID

		updatedTxn := *transaction

		updatedTxn.Account.ID = accountToTransfer.ID

		err = s.transactionStore.Create(ctx, &updatedTxn, tx)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.accountSvc.UpdateWithTx(ctx, account, tx)
	if err != nil {
		return nil, err
	}

	err = s.transactionStore.Create(ctx, transaction, tx)
	if err != nil {
		return nil, err
	}

	if transaction.Type == "SAVINGS" {
		savings := &models.Savings{
			UserID: transaction.UserID, Amount: transaction.Amount, Category: transaction.Category,
			StartDate: transaction.TransactionDate, TransactionID: transaction.ID, Status: "ACTIVE",
		}

		err = s.savingsSvc.CreateWithTx(ctx, savings, tx)
		if err != nil {
			return nil, err
		}
	}

	if transaction.Type == "WITHDRAW" {
		saving, er := s.savingsSvc.GetByTransactionID(ctx, int(transaction.WithdrawFrom))
		if er != nil {
			return nil, er
		}

		if (saving.CurrentValue != 0 && saving.CurrentValue < transaction.Amount) || (saving.CurrentValue == 0 && saving.Amount < transaction.Amount) {
			return nil, errors.New("withdrawal amount cannot exceed the saved amount")
		}

		saving.WithdrawnAmount += transaction.Amount

		if (saving.CurrentValue != 0 && saving.CurrentValue == transaction.Amount) || (saving.CurrentValue == 0 && saving.Amount == transaction.Amount) {
			saving.Status = "INACTIVE"
		}

		err = s.savingsSvc.UpdateWithTx(ctx, saving, true, tx)
		if err != nil {
			return nil, err
		}

	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	newTransaction, err := s.GetByID(ctx, transaction.ID)
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (s *transactionSvc) GetByID(ctx *gofr.Context, id int) (*models.Transaction, error) {
	userID, _ := ctx.Value("userID").(int)

	transaction, err := s.transactionStore.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionSvc) GetAll(ctx *gofr.Context, f *filters.Transactions) ([]*models.Transaction, error) {
	userID, _ := ctx.Value("userID").(int)

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
	} else if originalTransaction.Type == "SAVINGS" || originalTransaction.Type == "EXPENSE" {
		account.Balance += originalTransaction.Amount
	} else if originalTransaction.Type == "SELF TRANSFER" {
		fromAcc, err := s.accountSvc.GetByIDForUpdate(ctx, originalTransaction.MetaData.TransferFrom, userID, tx)
		if err != nil {
			return nil, err
		}
		toAcc, err := s.accountSvc.GetByIDForUpdate(ctx, originalTransaction.MetaData.TransferTo, userID, tx)
		if err != nil {
			return nil, err
		}

		fromAcc.Balance += originalTransaction.Amount
		toAcc.Balance -= originalTransaction.Amount

		_, err = s.accountSvc.UpdateWithTx(ctx, fromAcc, tx)
		if err != nil {
			return nil, err
		}
		_, err = s.accountSvc.UpdateWithTx(ctx, toAcc, tx)
		if err != nil {
			return nil, err
		}
	}

	// Apply the effect of the updated transaction
	if transaction.Type == "INCOME" {
		_, er := s.savingsSvc.GetByTransactionID(ctx, transaction.ID)
		if er == nil {
			err = s.savingsSvc.DeleteWithTx(ctx, transaction.ID, tx)
			if err != nil {
				return nil, err
			}
		}

		account.Balance += transaction.Amount
	} else if transaction.Type == "SAVINGS" {
		savings := &models.Savings{
			Status:        "ACTIVE",
			UserID:        transaction.UserID,
			Amount:        transaction.Amount,
			Category:      transaction.Category,
			StartDate:     transaction.TransactionDate,
			TransactionID: transaction.ID,
		}
		// Check if savings entry already exists for this transaction
		_, err = s.savingsSvc.GetByTransactionID(ctx, transaction.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Create new savings record if not found
				err = s.savingsSvc.CreateWithTx(ctx, savings, tx)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// Update existing if it exists
			err = s.savingsSvc.UpdateWithTx(ctx, savings, true, tx)
			if err != nil {
				return nil, err
			}
		}

		account.Balance -= transaction.Amount
	} else if transaction.Type == "EXPENSE" {
		_, er := s.savingsSvc.GetByTransactionID(ctx, transaction.ID)
		if er == nil {
			err = s.savingsSvc.DeleteWithTx(ctx, transaction.ID, tx)
			if err != nil {
				return nil, err
			}
		}

		account.Balance -= transaction.Amount
	} else if transaction.Type == "SELF TRANSFER" {
		_, er := s.savingsSvc.GetByTransactionID(ctx, transaction.ID)
		if er == nil {
			err = s.savingsSvc.DeleteWithTx(ctx, transaction.ID, tx)
			if err != nil {
				return nil, err
			}
		}

		d := account.Balance - transaction.Amount
		account.Balance = d

		accountToTransfer, er := s.accountSvc.GetByID(ctx, transaction.MetaData.TransferTo)
		if er != nil {
			return nil, er
		}

		accountToTransfer.Balance += transaction.Amount

		_, err = s.accountSvc.UpdateWithTx(ctx, accountToTransfer, tx)
		if err != nil {
			return nil, err
		}

		transaction.MetaData.TransferFrom = account.ID

		updatedTxn := *transaction

		updatedTxn.Account.ID = accountToTransfer.ID

		err = s.transactionStore.Create(ctx, &updatedTxn, tx)
		if err != nil {
			return nil, err
		}
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
	} else if originalTransaction.Type == "SELF TRANSFER" && originalTransaction.Account.ID == originalTransaction.MetaData.TransferFrom {
		account.Balance += originalTransaction.Amount
	} else if originalTransaction.Type == "SELF TRANSFER" && originalTransaction.Account.ID == originalTransaction.MetaData.TransferTo {
		account.Balance -= originalTransaction.Amount
	}

	err = s.transactionStore.Delete(ctx, id, tx)
	if err != nil {
		return err
	}

	err = s.savingsSvc.DeleteWithTx(ctx, id, tx)
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
