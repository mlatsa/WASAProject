package database

import "database/sql"

// Finds the last ID
func (db *appdbimpl) GetConversationId(user_id1 int, user_id2 int) (int, error) {
	var conversation_id int
	err := db.c.QueryRow(`SELECT C1.conversation_id
					FROM CONVERSATIONS_USERS AS C1
					INNER JOIN CONVERSATIONS_USERS AS C2 ON C1.conversation_id = C2.conversation_id
					LEFT JOIN GROUPS ON GROUPS.conversation_id = C1.conversation_id
					WHERE GROUPS.conversation_id IS NULL AND C1.user_id = ? AND C2.user_id = ?`, user_id1, user_id2).Scan(&conversation_id)
	if err == sql.ErrNoRows {
		conversation_id = 0
	}
	return conversation_id, err
}
