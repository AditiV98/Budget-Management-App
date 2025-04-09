package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	"google.golang.org/api/idtoken"
	"moneyManagement/models"
	"moneyManagement/services"
	"net/http"
	"net/url"
	"time"
)

type authSvc struct {
	RefreshSecret string
	AccessSecret  string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
}

func New(refreshSecret, accessSecret, ClientID, ClientSecret, RedirectURL string) services.Auth {
	return &authSvc{
		RefreshSecret: refreshSecret,
		AccessSecret:  accessSecret,
		ClientID:      ClientID,
		ClientSecret:  ClientSecret,
		RedirectURL:   RedirectURL,
	}
}

func (s *authSvc) GenerateGoogleToken(ctx *gofr.Context, code string) (map[string]interface{}, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)
	data.Set("redirect_uri", s.RedirectURL) // same as in Google console
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *authSvc) GenerateRefreshToken(claims *models.GoogleClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"userID":      claims.Sub,
		"email":       claims.Email,
		"name":        claims.Name,
		"given_name":  claims.GivenName,
		"family_name": claims.FamilyName,
		"picture":     claims.Picture,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(s.RefreshSecret))
}

func (s *authSvc) GenerateAccessToken(claims *models.GoogleClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"userID":      claims.EntityID,
		"email":       claims.Email,
		"name":        claims.Name,
		"given_name":  claims.GivenName,
		"family_name": claims.FamilyName,
		"picture":     claims.Picture,
		"exp":         time.Now().Add(5 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(s.AccessSecret))
}

func (s *authSvc) ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.RefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func (s *authSvc) VerifyGoogleIDToken(ctx context.Context, idToken string) (*models.GoogleClaims, error) {
	// Validate the id_token using Google's token verifier
	payload, err := idtoken.Validate(ctx, idToken, s.ClientID)
	if err != nil {
		return nil, err
	}

	// Extract useful claims
	email, _ := payload.Claims["email"].(string)
	sub, _ := payload.Claims["sub"].(string)
	name, _ := payload.Claims["name"].(string)
	pic, _ := payload.Claims["picture"].(string)
	givenName, _ := payload.Claims["given_name"].(string)
	familyName, _ := payload.Claims["family_name"].(string)

	if email == "" || sub == "" {
		return nil, errors.New("invalid id_token: missing sub or email")
	}

	return &models.GoogleClaims{
		Sub:        sub,
		Email:      email,
		Name:       name,
		Picture:    pic,
		GivenName:  givenName,
		FamilyName: familyName,
	}, nil
}
