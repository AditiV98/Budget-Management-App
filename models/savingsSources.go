package models

type SavingsSources struct {
	ID            int     `json:"id"`
	SavingID      int     `json:"savingID"`
	TransactionID int     `json:"transactionID"`
	Amount        float64 `json:"amount"`
	CreatedAt     string  `json:"createdAt"`
	DeletedAt     string  `json:"deletedAt"`
}
