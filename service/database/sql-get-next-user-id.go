package database

// Get the next user id
func (db *appdbimpl) GetNextUserId() (int, error) {
	var id int
	err := db.c.QueryRow("SELECT CASE WHEN COUNT(user_id) > 0 THEN MAX(user_id) ELSE 0 END AS ID FROM USERS;").Scan(&id)
	id = id + 1
	return id, err
}
