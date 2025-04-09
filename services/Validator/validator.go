package Validator

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"moneyManagement/services"
)

type validatorSvc struct {
	AccessSecret string
}

func New(accessSecret string) services.Validator {
	return &validatorSvc{
		AccessSecret: accessSecret,
	}
}

func (s *validatorSvc) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure token algorithm is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.AccessSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	return claims, nil
}
