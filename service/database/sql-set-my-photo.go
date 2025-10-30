package database

// Set user profile photo
func (db *appdbimpl) SetMyPhoto(id int, photo string) error {
	_, err := db.c.Exec(`UPDATE USERS SET user_photo=? WHERE user_id=?`, photo, id)
	return err
}
