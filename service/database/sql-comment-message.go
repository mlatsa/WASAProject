package database

// Comment message
func (db *appdbimpl) CommentMessage(msg_id int, from_user_id int, reaction string) (int, error) {
	var reaction_id int
	var err error

	reaction_id, _ = db.GetNextCommentId()

	_, err = db.c.Exec(`INSERT INTO REACTIONS
		(reaction_id, from_user_id, msg_id, reaction)
		VALUES (?, ?, ?, ?)`,
		reaction_id, from_user_id, msg_id, reaction)

	return reaction_id, err
}
