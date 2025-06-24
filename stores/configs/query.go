package configs

const (
	createConfig = "INSERT INTO configs (user_id, created_at) VALUES (?, ?)"
	updateConfig = "UPDATE configs SET is_auto_read=?,refresh_token=? WHERE user_id=?"
	getConfig    = "SELECT is_auto_read,refresh_token,created_at,updated_at FROM configs WHERE user_id=?"
)
