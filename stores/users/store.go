package users

import (
	"database/sql"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type userStore struct{}

func New() stores.User {
	return &userStore{}
}

func (s *userStore) Create(ctx *gofr.Context, user *models.User) error {
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	res, err := ctx.SQL.ExecContext(ctx, createUser, user.FirstName, user.LastName, user.Email, user.Status, createdAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return err
	}

	user.ID = int(id)

	return nil
}

func (s *userStore) GetByID(ctx *gofr.Context, id int) (*models.User, error) {
	var (
		user      models.User
		createdAt time.Time
		deletedAt sql.NullString
	)

	err := ctx.SQL.QueryRowContext(ctx, getByID, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Status, &createdAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, datasource.ErrorDB{Err: err, Message: "error fetching user by id"}
	}

	user.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

	if deletedAt.Valid {
		user.DeletedAt = deletedAt.String
	}

	return &user, nil
}

func (s *userStore) GetAll(ctx *gofr.Context, f *filters.User) ([]*models.User, error) {
	var allUsers []*models.User

	clause, val := f.WhereClause()

	q := getAll + clause

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
			user      models.User
			createdAt time.Time
			deletedAt sql.NullString
		)

		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Status, &createdAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		// Format as UTC in RFC3339 (standard for API responses)
		user.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")

		if deletedAt.Valid {
			user.DeletedAt = deletedAt.String
		}

		allUsers = append(allUsers, &user)
	}

	return allUsers, nil
}

func (s *userStore) Update(ctx *gofr.Context, user *models.User) error {
	_, err := ctx.SQL.ExecContext(ctx, updateUser, user.FirstName, user.LastName, user.Email, user.Status, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userStore) Delete(ctx *gofr.Context, id int) error {
	deletedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := ctx.SQL.ExecContext(ctx, deleteUser, "INACTIVE", deletedAt, id)
	if err != nil {
		return err
	}

	return nil
}
