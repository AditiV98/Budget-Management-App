package accounts

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/sql"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
)

type accountSvc struct {
	accountStore stores.Account
	userSvc      services.User
}

func New(accountStore stores.Account, userSvc services.User) services.Account {
	return &accountSvc{
		accountStore: accountStore,
		userSvc:      userSvc,
	}
}

func (s *accountSvc) Create(ctx *gofr.Context, account *models.Account) (*models.Account, error) {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	account.Status = "ACTIVE"
	account.UserID = userID

	id, err := s.accountStore.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	newAccount, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (s *accountSvc) GetByID(ctx *gofr.Context, id int) (*models.Account, error) {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	account, err := s.accountStore.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountSvc) GetAll(ctx *gofr.Context, f *filters.Account) ([]*models.Account, error) {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	f.UserID = userID

	accounts, err := s.accountStore.GetAll(ctx, f)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountSvc) Update(ctx *gofr.Context, account *models.Account) (*models.Account, error) {
	tx, err := ctx.SQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	userID, _ := ctx.Value("userID").(int)

	account.UserID = userID

	err = s.accountStore.Update(ctx, account, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedAccount, err := s.GetByID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (s *accountSvc) UpdateWithTx(ctx *gofr.Context, account *models.Account, tx *sql.Tx) (*models.Account, error) {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return nil, err
	//}

	account.UserID = userID

	err := s.accountStore.Update(ctx, account, tx)
	if err != nil {
		return nil, err
	}

	updatedAccount, err := s.GetByID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (s *accountSvc) Delete(ctx *gofr.Context, id int) error {
	userID, _ := ctx.Value("userID").(int)

	//user, err := s.userSvc.GetAll(ctx, &filters.User{Email: userEmail})
	//if err != nil {
	//	return err
	//}

	_, err := s.accountStore.GetByID(ctx, id, userID)
	if err != nil {
		return errors.New("unauthorised")
	}

	err = s.accountStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *accountSvc) GetByIDForUpdate(ctx *gofr.Context, id, userID int, tx *sql.Tx) (*models.Account, error) {
	account, err := s.accountStore.GetByIDForUpdate(ctx, id, userID, tx)
	if err != nil {
		return nil, err
	}

	return account, nil
}
