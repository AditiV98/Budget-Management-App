package filters

import "strings"

type User struct {
	Email  string `json:"email"`
	clause string
	args   []interface{}
}

func (f *User) WhereClause() (clause string, values []interface{}) {
	if f.Email != "" {
		f.clause += `email=? AND`
		f.args = append(f.args, f.Email)
	}

	if f.clause != "" {
		f.clause = " WHERE " + strings.TrimRight(f.clause, " AND")
		f.clause += " AND deleted_at IS NULL"
	}

	return f.clause, f.args
}
