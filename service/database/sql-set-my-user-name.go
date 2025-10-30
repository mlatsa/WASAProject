package database

// Set user name
func (db *appdbimpl) SetMyUserName(id int, name string) error {
	_, err := db.c.Exec(`UPDATE USERS SET user_name=? WHERE user_id=?`, name, id)
	return err
}
