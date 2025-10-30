package database

// Retrieve the groups that belongs to the user
func (db *appdbimpl) GetMyGroups(user_id int) ([]Group, error) {
	// Query data
	rows, err := db.c.Query(`SELECT group_id, group_name, group_photo, GROUPS.conversation_id
					FROM GROUPS
					INNER JOIN CONVERSATIONS_USERS ON CONVERSATIONS_USERS.conversation_id = GROUPS.conversation_id
					WHERE CONVERSATIONS_USERS.user_id = ?;`, user_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var groups []Group
	for rows.Next() {
		var grp Group
		err := rows.Scan(&grp.Group_id, &grp.Group_name, &grp.Group_photo, &grp.Conversation_id)
		if err != nil {
			return groups, err
		}
		groups = append(groups, grp)
	}
	err = rows.Err()
	if err != nil {
		return groups, err
	}

	return groups, nil
}
