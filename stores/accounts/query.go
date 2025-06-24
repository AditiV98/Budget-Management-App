package accounts

const (
	createAccount  = "INSERT INTO accounts (user_id, name, type,balance,status,expense_categories,saving_categories,created_at) VALUES (?, ?, ?, ?,?,?,?,?)"
	getByIDAccount = "SELECT id,user_id, name, type,balance,status,expense_categories,saving_categories,created_at,deleted_at,bank_email_address FROM accounts WHERE id=? AND user_id=?"
	getAllAccount  = "SELECT id,user_id, name, type,balance,status,expense_categories,saving_categories,created_at,deleted_at,bank_email_address FROM accounts"
	updateAccount  = "UPDATE accounts SET name=?,type=?,balance=?,status=?,expense_categories=?,saving_categories=?,bank_email_address=? WHERE id=? AND user_id=?"
	deleteAccount  = "UPDATE accounts SET status=?,deleted_at=? WHERE id=?"
)
