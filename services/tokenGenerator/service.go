package tokenGenerator

import (
	"errors"
	"fmt"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"moneyManagement/models"
	"moneyManagement/services"
	"moneyManagement/stores"
	"os"
	"time"
)

var scopes = []string{"https://mail.google.com/",
	"https://www.googleapis.com/auth/gmail.modify",
	"https://www.googleapis.com/auth/gmail.readonly",
	"https://www.googleapis.com/auth/gmail.compose",
	"https://www.googleapis.com/auth/gmail.send"}

type tokenService struct {
	configsStore stores.Configs
}

func New(configsStore stores.Configs) services.TokenGeneratorService {
	return &tokenService{configsStore: configsStore}
}

func (s *tokenService) GenerateAuthURLs(ctx *gofr.Context) (string, error) {
	credentials := os.Getenv("GCP_CREDENTIALS")
	b := []byte(credentials)

	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}

	config.RedirectURL = "http://localhost:3000/settings"
	//os.Getenv("REDIRECT_URL")

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return authURL, nil
}

func (s *tokenService) GenerateTokens(ctx *gofr.Context, m models.Code) (*models.RefreshToken, error) {
	userID, _ := ctx.Value("userID").(int)

	credentials := os.Getenv("GCP_CREDENTIALS")
	b := []byte(credentials)

	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}

	config.RedirectURL = "http://localhost:3000/settings"

	token, err := config.Exchange(ctx, m.ExchangeCode)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve token from web: %v", err))
	}

	newToken := models.Token{RefreshToken: token.RefreshToken, AccessToken: token.AccessToken,
		Expiry: token.Expiry.Format(time.RFC3339Nano), TokenType: token.TokenType}

	tokens := models.RefreshToken{Mail: newToken}

	configs, err := s.configsStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	configs.RefreshToken = &tokens

	err = s.configsStore.Update(ctx, configs)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}
