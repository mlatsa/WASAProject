package database

import "database/sql"

// Get the user id
func (db *appdbimpl) GetUserId(user_name string) (int, error) {
	var user_id int
	err := db.c.QueryRow("SELECT user_id FROM USERS WHERE user_name=?;", user_name).Scan(&user_id)
	if err == sql.ErrNoRows {
		user_id = 0
		err = nil
	}
	return user_id, err
}
