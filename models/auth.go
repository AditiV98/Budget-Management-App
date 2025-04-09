package models

import "github.com/google/uuid"

type CodeRequest struct {
	Code string `json:"code"`
}

type LoginRequest struct {
	TenantID     uuid.UUID              `json:"tenantID"`
	Provider     string                 `json:"provider"`
	ProviderData map[string]interface{} `json:"providerData"`
	AppData      map[string]interface{} `json:"appData"`
	Platform     string                 `json:"platform"`
	Roles        []string               `json:"roles"`
}

type RefreshRequest struct {
	TenantID     uuid.UUID              `json:"tenantID"`
	RefreshToken string                 `json:"refreshToken"`
	AppData      map[string]interface{} `json:"appData"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type GoogleClaims struct {
	Sub        string
	Email      string
	Name       string
	Picture    string
	GivenName  string
	FamilyName string
	EntityID   int
}
