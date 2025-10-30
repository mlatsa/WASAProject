package database

// Create a new group
func (db *appdbimpl) LeaveGroup(user_id int, group_id int) error {
	var err error
	_, err = db.c.Exec(`DELETE FROM CONVERSATIONS_USERS
					WHERE ROWID IN (
						SELECT a.ROWID FROM CONVERSATIONS_USERS a
						INNER JOIN GROUPS ON GROUPS.conversation_id = CONVERSATIONS_USERS.conversation_id
						WHERE GROUPS.group_id = ? AND CONVERSATIONS_USERS.user_id = ?
					);`, group_id, user_id)
	return err
}
