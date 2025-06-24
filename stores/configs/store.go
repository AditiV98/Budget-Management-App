package configs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/models"
	"moneyManagement/stores"
	"time"
)

type configsStore struct{}

func New() stores.Configs {
	return &configsStore{}
}

func (s *configsStore) Create(ctx *gofr.Context, userID int) error {
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	res, err := ctx.SQL.ExecContext(ctx, createConfig, userID, createdAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		return err
	}

	return nil
}

func (s *configsStore) Update(ctx *gofr.Context, config *models.Config) error {
	token, _ := json.Marshal(config.RefreshToken)

	_, err := ctx.SQL.ExecContext(ctx, updateConfig, config.IsAutoRead, string(token), config.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *configsStore) Get(ctx *gofr.Context, userID int) (*models.Config, error) {
	var (
		currConfig   models.Config
		refreshToken sql.NullString
		createdAt    time.Time
		updatedAt    time.Time
	)

	currConfig.UserID = userID

	err := ctx.SQL.QueryRowContext(ctx, getConfig, userID).Scan(&currConfig.IsAutoRead, &refreshToken,
		&createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.Warnf("No configs found for userID : %d", userID)

			return nil, errors.New("no record found")
		}

		return nil, err
	}

	if refreshToken.Valid {
		err = json.Unmarshal([]byte(refreshToken.String), &currConfig.RefreshToken)
		if err != nil {
			return nil, err
		}
	}

	return &currConfig, nil
}
