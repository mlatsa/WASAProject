package database

// Uncomment a message
func (db *appdbimpl) UncommentMessage(msg_id int, from_user_id int) error {
	var err error
	_, err = db.c.Exec(`DELETE FROM REACTIONS WHERE msg_id = ? AND from_user_id = ?`, msg_id, from_user_id)
	return err
}
