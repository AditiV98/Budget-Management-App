package models

import (
	"errors"
	"time"
)

type Type string

const (
	INCOME  Type = "INCOME"
	EXPENSE Type = "EXPENSE"
	SAVINGS Type = "SAVINGS"
)

var ExpenseCategories = map[string]struct{}{
	"Housing": {}, "Utilities": {}, "Groceries": {}, "Transportation": {}, "Education": {},
	"Healthcare": {}, "Loan & Debt Payments": {}, "Dining Out": {}, "Entertainment": {},
	"Shopping": {}, "Travel": {}, "Gifts & Donations": {}, "Fitness & Wellness": {},
	"Childcare": {}, "Home Maintenance": {}, "Pet Care": {}, "Self-Development": {},
}

var SavingsCategories = map[string]struct{}{
	"FD": {}, "Mutual Funds": {}, "Stocks": {}, "Gold ETFs": {}, "Other": {},
}

type Transaction struct {
	ID              int            `json:"id"`
	UserID          int            `json:"userID"`
	Account         AccountDetails `json:"account"`
	Amount          float64        `json:"amount"`
	Type            Type           `json:"type"`
	Category        string         `json:"category"`
	Description     string         `json:"description"`
	TransactionDate string         `json:"transactionDate"`
	CreatedAt       string         `json:"createdAt"`
	DeletedAt       string         `json:"deletedAt,omitempty"`
	WithdrawFrom    int64          `json:"withdrawFrom"`
	MetaData        MetaData       `json:"metaData"`
	Saving          SavingDetails  `json:"saving"`
}

type MetaData struct {
	TransferTo   int `json:"transferTo"`
	TransferFrom int `json:"transferFrom"`
}

type AccountDetails struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SavingDetails struct {
	ID            int     `json:"id"`
	Status        string  `json:"status"`
	TransactionID int     `json:"transactionID"`
	Amount        float64 `json:"amount"`
	CurrentValue  float64 `json:"currentValue"`
}

// Validate checks if the transaction fields are valid
func (t *Transaction) Validate() error {
	// Validate category for EXPENSES type
	if t.Type == EXPENSE {
		if _, exists := ExpenseCategories[t.Category]; !exists {
			return errors.New("invalid category for EXPENSES type")
		}
	}

	// Validate dates
	if _, err := time.Parse("2006-01-02", t.TransactionDate); err != nil {
		return errors.New("invalid transaction date format, use YYYY-MM-DD")
	}
	if _, err := time.Parse("2006-01-02", t.CreatedAt); err != nil {
		return errors.New("invalid createdAt date format, use YYYY-MM-DD")
	}

	return nil
}
