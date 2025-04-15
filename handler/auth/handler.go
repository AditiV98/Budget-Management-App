package auth

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
)

type auth struct {
	authSvc services.Auth
	userSvc services.User
}

func New(authSvc services.Auth, userSvc services.User) handler.Auth {
	return &auth{authSvc: authSvc, userSvc: userSvc}
}

func (h *auth) CreateToken(ctx *gofr.Context) (interface{}, error) {
	var req models.CodeRequest

	err := ctx.Bind(&req)
	if err != nil {
		return nil, errors.New("bind error")
	}

	result, err := h.authSvc.GenerateGoogleToken(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *auth) Login(ctx *gofr.Context) (interface{}, error) {
	var req models.LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return nil, errors.New("bind error")
	}

	// Get id_token from providerData
	idToken, ok := req.ProviderData["token"].(string)
	if !ok || idToken == "" {
		return nil, errors.New("missing id_token")
	}

	// Validate and parse the id_token
	claims, err := h.authSvc.VerifyGoogleIDToken(ctx.Context, idToken)
	if err != nil {
		return nil, err
	}

	// Generate a custom refresh token (valid for 1 day)
	refreshToken, err := h.authSvc.GenerateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{RefreshToken: refreshToken}, nil
}

func (h *auth) Refresh(ctx *gofr.Context) (interface{}, error) {
	var req models.RefreshRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, errors.New("bind error")
	}

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return nil, errors.New("missing refresh token")
	}

	// Validate and parse refresh token
	jwtClaims, err := h.authSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	email, _ := jwtClaims["email"].(string)
	sub, _ := jwtClaims["sub"].(string)
	name, _ := jwtClaims["name"].(string)
	pic, _ := jwtClaims["picture"].(string)
	givenName, _ := jwtClaims["given_name"].(string)
	familyName, _ := jwtClaims["family_name"].(string)

	claims := &models.GoogleClaims{
		Sub:        sub,
		Email:      email,
		Name:       name,
		Picture:    pic,
		GivenName:  givenName,
		FamilyName: familyName,
	}

	err = h.userSvc.AuthAdaptor(ctx, claims)
	if err != nil {
		return nil, err
	}

	// Generate short-lived access token (valid for 5 minutes)
	accessToken, err := h.authSvc.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{AccessToken: accessToken}, nil
}
