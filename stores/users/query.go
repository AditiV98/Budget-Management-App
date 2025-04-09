package users

const (
	createUser = "INSERT INTO users (first_name, last_name, email,status, created_at) VALUES (?, ?, ?, ?,?)"
	getByID    = "SELECT id,first_name,last_name,email,status,created_at,deleted_at FROM users WHERE id=?"
	getAll     = "SELECT id,first_name,last_name,email,status,created_at,deleted_at FROM users"
	updateUser = "UPDATE users SET first_name=?,last_name=?,email=?,status=? WHERE id=?"
	deleteUser = "UPDATE users SET status=?,deleted_at=? WHERE id=?"
)
