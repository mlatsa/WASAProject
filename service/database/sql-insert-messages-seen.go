package database

// Mark messages as seen but exclude messages sent by the user itself!
func (db *appdbimpl) MarkMessagesSeen(seen_by_user_id int, conversation_id int) error {
	var err error
	_, err = db.c.Exec(`INSERT OR IGNORE INTO MESSAGES_SEEN (msg_id, seen_by_user_id)
						SELECT msg_id, ?
						FROM MESSAGES
						WHERE conversation_id = ? AND from_user_id != ?`, seen_by_user_id, conversation_id, seen_by_user_id)
	return err
}
