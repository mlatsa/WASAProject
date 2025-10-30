package database

// Get the next message id
func (db *appdbimpl) GetNextMessageId() (int, error) {
	var id int
	err := db.c.QueryRow("SELECT CASE WHEN COUNT(msg_id) > 0 THEN MAX(msg_id) ELSE 0 END AS ID FROM MESSAGES;").Scan(&id)
	id = id + 1
	return id, err
}
