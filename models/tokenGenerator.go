package models

type Code struct {
	Scope        []string `json:"scope"`
	ExchangeCode string   `json:"exchangeCode"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
}

type RefreshToken struct {
	Mail Token `json:"mail"`
}

type Config struct {
	UserID       int           `json:"userID"`
	RefreshToken *RefreshToken `json:"refreshToken"`
	IsAutoRead   bool          `json:"isAutoRead"`
}
