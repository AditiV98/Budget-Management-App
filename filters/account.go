package filters

import "strings"

type Account struct {
	UserID int `json:"email"`
	clause string
	args   []interface{}
}

func (f *Account) WhereClause() (clause string, values []interface{}) {
	if f.UserID != 0 {
		f.clause += `user_id=? AND`
		f.args = append(f.args, f.UserID)
	}

	if f.clause != "" {
		f.clause = " WHERE " + strings.TrimRight(f.clause, " AND")
		f.clause += " AND deleted_at IS NULL"
	}

	return f.clause, f.args
}
