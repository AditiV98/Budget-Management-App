package models

type Savings struct {
	ID              int            `json:"id"`
	UserID          int            `json:"userID"`
	TransactionID   int            `json:"transactionID"`
	Amount          float64        `json:"amount"`
	WithdrawnAmount float64        `json:"withdrawnAmount"`
	Status          string         `json:"status"`
	Category        string         `json:"category"`
	CurrentValue    float64        `json:"currentValue"`
	StartDate       string         `json:"startDate"`
	MaturityDate    string         `json:"maturityDate,omitempty"`
	CreatedAt       string         `json:"createdAt"`
	DeletedAt       string         `json:"deletedAt,omitempty"`
	Account         AccountDetails `json:"account"`
}
