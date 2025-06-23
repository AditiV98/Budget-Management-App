package savings

const (
	createSavings                  = "INSERT INTO savings (user_id,transaction_id,category,amount,current_value,start_date,maturity_date,created_at,status,description) VALUES (?,?, ?, ?, ?, ?, ?,?,?,?)"
	getByIDSavings                 = "SELECT id,user_id,transaction_id,category,amount,current_value,start_date,maturity_date,created_at,deleted_at,status,withdrawn_amount,description FROM savings WHERE id=?"
	getAllSavings                  = "SELECT id,user_id,transaction_id,category,amount,current_value,start_date,maturity_date,created_at,deleted_at,status,withdrawn_amount,description FROM savings"
	updateSavings                  = "UPDATE savings SET current_value=?,maturity_date=?,status=?,description=? WHERE id=?"
	deleteSavings                  = "UPDATE savings SET deleted_at=? WHERE id=?"
	updateSavingsWithTransactionID = "UPDATE savings SET category=?,amount=?,current_value=?,start_date=?,maturity_date=?,withdrawn_amount=?,status=? WHERE transaction_id=?"
	getByTransactionIDSavings      = "SELECT id,user_id,transaction_id,category,amount,current_value,start_date,maturity_date,created_at,deleted_at,withdrawn_amount FROM savings WHERE transaction_id=? AND deleted_at IS NULL"
	deleteSavingsByTransactionID   = "UPDATE savings SET deleted_at=?,status=? WHERE transaction_id=?"
)
