package filters

import "strings"

type RecurringTransactions struct {
	Type      []string `json:"type"`
	UserID    int      `json:"userID"`
	AccountID int      `json:"accountID"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Category  []string `json:"category"`
	clause    string
	args      []interface{}
}

func (t *RecurringTransactions) WhereClause() (clause string, values []interface{}) {
	if len(t.Type) != 0 {
		t.clause += `t.type IN (` + placeHolders(len(t.Type)) + `) AND`

		for i := range t.Type {
			t.args = append(t.args, t.Type[i])
		}
	}

	if len(t.Category) != 0 {
		t.clause += ` t.category IN (` + placeHolders(len(t.Category)) + `) AND`

		for i := range t.Category {
			t.args = append(t.args, t.Category[i])
		}
	}

	if t.UserID != 0 {
		t.clause += ` t.user_id=? AND`
		t.args = append(t.args, t.UserID)
	}

	if t.AccountID != 0 {
		t.clause += ` t.account_id=? AND`
		t.args = append(t.args, t.AccountID)
	}

	if t.StartDate != "" && t.EndDate != "" {
		t.clause += ` t.transaction_date>=? AND t.transaction_date<=? AND`
		t.args = append(t.args, t.StartDate, t.EndDate)
	}

	if t.clause != "" {
		t.clause = " WHERE " + strings.TrimRight(t.clause, " AND")
		t.clause += " AND t. deleted_at IS NULL"
	}

	return t.clause, t.args
}
