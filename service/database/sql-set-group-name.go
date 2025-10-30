package database

// Set group name
func (db *appdbimpl) SetGroupName(group_id int, group_name string) error {
	_, err := db.c.Exec(`UPDATE GROUPS SET group_name=? WHERE group_id=?`, group_name, group_id)
	return err
}
