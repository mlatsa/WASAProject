package database

// Changes the name of a user
func (db *appdbimpl) ChangeName(id int, new_name string) error {
	_, err := db.c.Exec("UPDATE USERS SET (user_name) = (?) WHERE id = ?", new_name, id)
	return err
}
