package database

// Get the next comment id
func (db *appdbimpl) GetNextCommentId() (int, error) {
	var id int
	err := db.c.QueryRow("SELECT CASE WHEN COUNT(reaction_id) > 0 THEN MAX(reaction_id) ELSE 0 END AS ID FROM REACTIONS;").Scan(&id)
	id = id + 1
	return id, err
}
