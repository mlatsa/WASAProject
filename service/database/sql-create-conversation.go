package database

// Create a new conversation
func (db *appdbimpl) CreateConversation(conversation_id int, user_id1 int, user_id2 int) error {
	var err error
	_, err = db.c.Exec("INSERT INTO CONVERSATIONS_USERS (conversation_id, user_id) VALUES (?, ?);", conversation_id, user_id1)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("INSERT INTO CONVERSATIONS_USERS (conversation_id, user_id) VALUES (?, ?);", conversation_id, user_id2)
	return err
}
