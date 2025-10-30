package database

// Retrieve the list of users
func (db *appdbimpl) GetAllUsers() ([]User, error) {
	// Query data
	rows, err := db.c.Query(`SELECT user_id, user_name, ifnull(user_photo, '') FROM USERS ORDER BY user_name;`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.User_id, &user.User_name, &user.User_photo)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}

	return users, nil
}
