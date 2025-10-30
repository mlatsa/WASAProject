package database

// Create a new group
func (db *appdbimpl) CreateGroup(group_id int, group_name string, group_photo string, conversation_id int) error {
	var err error
	_, err = db.c.Exec("INSERT INTO GROUPS (group_id, group_name, group_photo, conversation_id) VALUES (?, ?, ?, ?);", group_id, group_name, group_photo, conversation_id)
	return err
}
