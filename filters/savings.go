package filters

import "strings"

type Savings struct {
	UserID    int      `json:"email"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Category  []string `json:"category"`
	Status    string   `json:"status"`
	clause    string
	args      []interface{}
}

func (f *Savings) WhereClause() (clause string, values []interface{}) {
	if f.UserID != 0 {
		f.clause += `user_id=? AND`
		f.args = append(f.args, f.UserID)
	}

	if f.StartDate != "" && f.EndDate != "" {
		f.clause += ` start_date>=? AND start_date<=? AND`
		f.args = append(f.args, f.StartDate, f.EndDate)
	}

	if len(f.Category) != 0 {
		f.clause += ` category IN (` + placeHolders(len(f.Category)) + `) AND`

		for i := range f.Category {
			f.args = append(f.args, f.Category[i])
		}
	}

	if f.Status != "" {
		f.clause += ` status=? AND`

		f.args = append(f.args, f.Status)
	}

	if f.clause != "" {
		f.clause = " WHERE " + strings.TrimRight(f.clause, " AND")
		f.clause += " AND deleted_at IS NULL"
	}

	return f.clause, f.args
}
