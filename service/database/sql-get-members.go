package database

// Retrieve the list of users
func (db *appdbimpl) GetMembers(group_id int) ([]User, error) {
	// Query data
	rows, err := db.c.Query(`
						SELECT USERS.user_id, USERS.user_name FROM CONVERSATIONS_USERS
						INNER JOIN GROUPS ON GROUPS.conversation_id = CONVERSATIONS_USERS.conversation_id
						INNER JOIN USERS ON USERS.user_id = CONVERSATIONS_USERS.user_id
						WHERE GROUPS.group_id = ?
						;`, group_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.User_id, &user.User_name)
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
