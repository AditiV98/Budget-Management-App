package savings

const (
	createSavings                  = "INSERT INTO savings (user_id,transaction_id,type,category,amount,current_value,start_date,maturity_date,created_at) VALUES (?,?, ?, ?, ?, ?, ?, ?,?)"
	getByIDSavings                 = "SELECT id,user_id,transaction_id,type,category,amount,current_value,start_date,maturity_date,created_at,deleted_at FROM savings WHERE id=?"
	getAllSavings                  = "SELECT id,user_id,transaction_id,type,category,amount,current_value,start_date,maturity_date,created_at,deleted_at FROM savings"
	updateSavings                  = "UPDATE savings SET type=?,category=?,amount=?,current_value=?,start_date=?,maturity_date=? WHERE id=?"
	deleteSavings                  = "UPDATE savings SET deleted_at=? WHERE id=?"
	updateSavingsWithTransactionID = "UPDATE savings SET type=?,category=?,amount=?,current_value=?,start_date=?,maturity_date=? WHERE transaction_id=?"
	getByTransactionIDSavings      = "SELECT id,user_id,transaction_id,type,category,amount,current_value,start_date,maturity_date,created_at,deleted_at FROM savings WHERE transaction_id=?"
)
