package database

// Check if a user belong to a conversation
func (db *appdbimpl) IsUserInConversation(user_id int, conversation_id int) (bool, error) {
	var belongTo bool
	err := db.c.QueryRow(`SELECT count() FROM CONVERSATIONS_USERS WHERE user_id=? AND conversation_id=?`, user_id, conversation_id).Scan(&belongTo)
	return belongTo, err
}
