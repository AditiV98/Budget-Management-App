package models

type Account struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	UserID            int      `json:"userID"`
	Balance           float64  `json:"balance"`
	Status            string   `json:"status"`
	ExpenseCategories []string `json:"expenseCategories"`
	SavingCategories  []string `json:"savingCategories"`
	CreatedAt         string   `json:"createdAt"`
	DeletedAt         string   `json:"deletedAt,omitempty"`
	BankEmailAddress  string   `json:"bankEmailAddress"`
}
