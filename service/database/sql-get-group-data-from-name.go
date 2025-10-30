package database

import "database/sql"

// Get the group id
func (db *appdbimpl) GetGroupDataFromName(group_name string) (int, string, int, error) {
	var group_id int
	var group_photo string
	var conversation_id int
	err := db.c.QueryRow("SELECT group_id, group_photo, conversation_id FROM GROUPS WHERE group_name=?;", group_name).Scan(&group_id, &group_photo, &conversation_id)
	if err == sql.ErrNoRows {
		group_id = 0
		err = nil
	}
	return group_id, group_photo, conversation_id, err
}
