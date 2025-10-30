package database

// Get the next conversation id
func (db *appdbimpl) GetNextConversationId() (int, error) {
	var id int
	err := db.c.QueryRow("SELECT CASE WHEN COUNT(conversation_id) > 0 THEN MAX(conversation_id) ELSE 0 END AS ID FROM CONVERSATIONS_USERS;").Scan(&id)
	id = id + 1
	return id, err
}
