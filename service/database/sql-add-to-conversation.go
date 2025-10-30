package database

// Create a new conversation
func (db *appdbimpl) AddUserToConversation(conversation_id int, user_id int) error {
	var err error
	_, err = db.c.Exec("INSERT INTO CONVERSATIONS_USERS (conversation_id, user_id) VALUES (?, ?);", conversation_id, user_id)
	return err
}
