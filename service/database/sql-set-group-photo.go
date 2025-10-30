package database

// Set group photo
func (db *appdbimpl) SetGroupPhoto(group_id int, photo string) error {
	_, err := db.c.Exec(`UPDATE GROUPS SET group_photo=? WHERE group_id=?`, photo, group_id)
	return err
}
