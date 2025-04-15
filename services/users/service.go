package users

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
)

type userSvc struct {
	userStore stores.User
}

func New(userStore stores.User) services.User {
	return &userSvc{
		userStore: userStore,
	}
}

func (s *userSvc) Create(ctx *gofr.Context, user *models.User) (*models.User, error) {
	user.Status = "ACTIVE"

	err := s.userStore.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	newUser, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userSvc) GetByID(ctx *gofr.Context, id int) (*models.User, error) {
	user, err := s.userStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userSvc) GetAll(ctx *gofr.Context, f *filters.User) ([]*models.User, error) {
	users, err := s.userStore.GetAll(ctx, f)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userSvc) Update(ctx *gofr.Context, user *models.User) (*models.User, error) {
	err := s.userStore.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userSvc) Delete(ctx *gofr.Context, id int) error {
	err := s.userStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userSvc) AuthAdaptor(ctx *gofr.Context, claims *models.GoogleClaims) error {
	user, err := s.GetAll(ctx, &filters.User{Email: claims.Email})
	if err != nil {
		return err
	}

	if user == nil {
		newUser := &models.User{Email: claims.Email, FirstName: claims.GivenName, LastName: claims.FamilyName}

		_, err = s.Create(ctx, newUser)
		if err != nil {
			return err
		}

		claims.EntityID = newUser.ID

		return nil
	}

	claims.EntityID = user[0].ID

	return nil
}
