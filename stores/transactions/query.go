package transactions

const (
	createTransaction   = "INSERT INTO transactions (user_id, account_id, amount,type,category,description,transaction_date,created_at,withdraw_from,meta_data) VALUES (?, ?, ?, ?,?,?,?,?,?,?)"
	getByIDTransactions = "SELECT t.id,t.user_id, t.account_id, t.amount,t.type,t.category,t.description,t.transaction_date,t.created_at,t.deleted_at,t.withdraw_from,t.meta_data,a.id,a.name FROM transactions as t INNER JOIN accounts as a ON t.account_id=a.id WHERE t.id=? AND t.user_id=?"
	getAllTransactions  = "SELECT t.id,t.user_id,t.account_id,t.amount, t.type,t.category,t.description,t.transaction_date,t.created_at,t.deleted_at,t.withdraw_from,t.meta_data,a.id,a.name,s.id,s.status,s.transaction_id,s.amount,s.current_value FROM transactions AS t INNER JOIN accounts AS a ON t.account_id = a.id LEFT JOIN savings AS s ON t.id = s.transaction_id AND s.deleted_at IS NULL"
	updateTransaction   = "UPDATE transactions SET account_id=?, amount=?,type=?,category=?,description=?,transaction_date=?,meta_data=? WHERE id=?"
	deleteTransaction   = "UPDATE transactions SET deleted_at=? WHERE id=?"
)
