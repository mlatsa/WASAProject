package database

// Get the next group id
func (db *appdbimpl) GetNextGroupId() (int, error) {
	var id int
	err := db.c.QueryRow("SELECT CASE WHEN COUNT(group_id) > 0 THEN MAX(group_id) ELSE 0 END AS ID FROM GROUPS;").Scan(&id)
	id = id + 1
	return id, err
}
