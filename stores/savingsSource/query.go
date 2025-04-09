package savingsSource

const (
	createSavingsSource  = "INSERT INTO savings_source (saving_id,transaction_id,amount) VALUES (?, ?, ?)"
	getByIDSavingsSource = "SELECT id,saving_id,transaction_id,amount,created_at,deleted_at FROM savings_source WHERE id=?"
	updateSavingsSource  = "UPDATE savings_source SET saving_id=?,transaction_id=?,amount=? WHERE id=?"
	deleteSavingsSource  = "UPDATE savings_source SET deleted_at=? WHERE id=?"
)
