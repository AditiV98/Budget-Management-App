package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"moneyManagement/services"
	"net/http"
	"regexp"
	"strings"
)

type access struct {
	Path           string
	Method         string
	Role           string
	ProtectedRoute bool
}

type ExemptPath struct {
	Path   string
	Method string
}

func getAccessMap() []access {
	return []access{
		{"^/user$", http.MethodPost, "ADMIN,USER", true},
		{"^/user$", http.MethodGet, "ADMIN", true},
		{"^/user/[0-9]+", http.MethodGet, "ADMIN,USER", true},
		{"^/user/[0-9]+", http.MethodPut, "ADMIN,USER", true},
		{"^/user/[0-9]+", http.MethodDelete, "ADMIN", true},

		{"^/account$", http.MethodPost, "ADMIN,USER", true},
		{"^/account$", http.MethodGet, "ADMIN,USER", true},
		{"^/account/[0-9]+", http.MethodGet, "ADMIN,USER", true},
		{"^/account/[0-9]+", http.MethodPut, "ADMIN,USER", true},
		{"^/account/[0-9]+", http.MethodDelete, "ADMIN,USER", true},

		{"^/savings$", http.MethodPost, "ADMIN,USER", true},
		{"^/savings$", http.MethodGet, "ADMIN,USER", true},
		{"^/savings/[0-9]+", http.MethodGet, "ADMIN,USER", true},
		{"^/savings/[0-9]+", http.MethodPut, "ADMIN,USER", true},
		{"^/savings/[0-9]+", http.MethodDelete, "ADMIN,USER", true},

		{"^/transaction$", http.MethodPost, "ADMIN,USER", true},
		{"^/transaction$", http.MethodGet, "ADMIN,USER", true},
		{"^/transaction/[0-9]+", http.MethodGet, "ADMIN,USER", true},
		{"^/transaction/[0-9]+", http.MethodPut, "ADMIN,USER", true},
		{"^/transaction/[0-9]+", http.MethodDelete, "ADMIN,USER", true},
	}
}

func Authorization(exemptPaths []ExemptPath, validator services.Validator, userSvc services.User) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkIfExempt(exemptPaths, r) {
				inner.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				ErrorResponse(w, http.StatusUnauthorized, "missing_token", "Authorization token required")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token using your validator
			claims, err := validator.ValidateToken(tokenStr)
			if err != nil {
				ErrorResponse(w, http.StatusUnauthorized, "invalid_token", "Token validation failed: "+err.Error())
				return
			}

			userID, ok := claims["userID"].(float64)
			if !ok {
				ErrorResponse(w, http.StatusUnauthorized, "invalid_token", "Token validation failed: "+err.Error())
				return
			}

			*r = *r.WithContext(context.WithValue(r.Context(), "userID", int(userID)))

			inner.ServeHTTP(w, r)
		})
	}
}

func ErrorResponse(w http.ResponseWriter, errStatusCode int, errCode, errReason string) {
	w.WriteHeader(errStatusCode)

	data, _ := json.Marshal(map[string]string{"code": errCode, "reason": errReason})

	_, _ = fmt.Fprintln(w, string(data))
}

// containsURL checks if the request is exempted from authorization
func checkIfExempt(urls []ExemptPath, r *http.Request) bool {
	for _, url := range urls {
		match, _ := regexp.MatchString(url.Path, r.URL.Path)
		if match && r.Method == url.Method {
			return true
		}
	}

	return false
}
