package recurringTransactions

const (
	createTransaction   = "INSERT INTO recurring_transactions (user_id, account_id, amount,type,category,description,frequency,custom_days,start_date,end_date,next_run,created_at) VALUES (?, ?, ?, ?,?,?,?,?,?,?,?,?)"
	getByIDTransactions = "SELECT t.id,t.user_id, t.account_id, t.amount,t.type,t.category,t.description,t.frequency,t.custom_days,t.start_date,t.end_date,t.last_run,t.next_run,t.created_at,t.deleted_at,a.name FROM recurring_transactions as t INNER JOIN accounts as a ON t.account_id=a.id WHERE t.id=? AND t.user_id=?"
	getAllTransactions  = "SELECT t.id,t.user_id, t.account_id, t.amount,t.type,t.category,t.description,t.frequency,t.custom_days,t.start_date,t.end_date,t.last_run,t.next_run,t.created_at,t.deleted_at,a.name FROM recurring_transactions as t INNER JOIN accounts as a ON t.account_id=a.id"
	updateTransaction   = "UPDATE recurring_transactions SET account_id=?, amount=?,type=?,category=?,description=?,frequency=?,custom_days=?,start_date=?,end_date=?,last_run=?,next_run=? WHERE id=?"
	deleteTransaction   = "UPDATE recurring_transactions SET deleted_at=? WHERE id=?"
)
