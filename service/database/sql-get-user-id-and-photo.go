package database

import "database/sql"

// Finds the last ID
func (db *appdbimpl) GetUserIdAndPhoto(user_name string) (int, string, error) {
	var user User
	err := db.c.QueryRow("SELECT user_id, ifnull(user_photo, '') FROM USERS WHERE user_name=?;", user_name).Scan(&user.User_id, &user.User_photo)
	if err == sql.ErrNoRows {
		err = nil
	}
	return user.User_id, user.User_photo, err
}
