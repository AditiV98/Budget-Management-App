package transactions

const (
	createTransaction   = "INSERT INTO transactions (user_id, account_id, amount,type,category,description,transaction_date,created_at) VALUES (?, ?, ?, ?,?,?,?,?)"
	getByIDTransactions = "SELECT t.id,t.user_id, t.account_id, t.amount,t.type,t.category,t.description,t.transaction_date,t.created_at,t.deleted_at,a.name FROM transactions as t INNER JOIN accounts as a ON t.account_id=a.id WHERE t.id=? AND t.user_id=?"
	getAllTransactions  = "SELECT t.id,t.user_id, t.account_id, t.amount,t.type,t.category,t.description,t.transaction_date," +
		"t.created_at,t.deleted_at,a.name FROM transactions as t INNER JOIN accounts as a ON t.account_id=a.id"
	updateTransaction = "UPDATE transactions SET account_id=?, amount=?,type=?,category=?,description=?,transaction_date=? WHERE id=?"
	deleteTransaction = "UPDATE transactions SET deleted_at=? WHERE id=?"
)
