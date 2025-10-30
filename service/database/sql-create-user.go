package database

// Set user name
func (db *appdbimpl) CreateUser(user_id int, user_name string) error {
	_, err := db.c.Exec(`INSERT INTO USERS VALUES (?, ?, NULL)`, user_id, user_name)
	return err
}
