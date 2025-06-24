package tokenGenerator

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
)

type tokenHandler struct {
	tokenSvc services.TokenGeneratorService
}

func New(tokenSvc services.TokenGeneratorService) handler.TokenGenerator {
	return &tokenHandler{tokenSvc: tokenSvc}
}

func (h tokenHandler) GenerateAuthURLs(ctx *gofr.Context) (interface{}, error) {
	authURLs, err := h.tokenSvc.GenerateAuthURLs(ctx)
	if err != nil {
		return nil, err
	}

	return authURLs, nil
}

func (h tokenHandler) GenerateTokens(ctx *gofr.Context) (interface{}, error) {
	var token models.Code

	err := ctx.Bind(&token)
	if err != nil {
		return nil, errors.New("bind error")
	}

	tokens, err := h.tokenSvc.GenerateTokens(ctx, token)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
