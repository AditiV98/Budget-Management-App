package models

type Frequency string

const (
	DAILY   Frequency = "DAILY"
	WEEKLY  Frequency = "WEEKLY"
	MONTHLY Frequency = "MONTHLY"
	CUSTOM  Frequency = "CUSTOM"
)

type RecurringTransaction struct {
	ID          int            `json:"id"`
	UserID      int            `json:"userID"`
	Account     AccountDetails `json:"account"`
	Amount      float64        `json:"amount"`
	Type        Type           `json:"type"`
	Category    string         `json:"category"`
	Description string         `json:"description"`
	Frequency   Frequency      `json:"frequency"`
	CustomDays  int            `json:"customDays"`
	StartDate   string         `json:"startDate"`
	EndDate     string         `json:"endDate"`
	LastRun     string         `json:"lastRun"`
	NextRun     string         `json:"nextRun"`
	CreatedAt   string         `json:"createdAt"`
	DeletedAt   string         `json:"deletedAt,omitempty"`
}
