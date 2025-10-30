package database

// Retrieve the conversations
func (db *appdbimpl) GetMyConversations(from_user_id int) ([]Conversation, error) {
	// Query data
	rows, err := db.c.Query(`
					SELECT CONVERSATIONS_USERS.conversation_id, USERS.user_id as user_or_group_id, user_name as user_or_group_name, ifnull(user_photo, '') as user_or_group_photo, 0 is_group,
						ifnull(last_timestamp, '') as last_timestamp, ifnull(last_msg, '') as last_msg
					FROM CONVERSATIONS_USERS
					INNER JOIN USERS ON CONVERSATIONS_USERS.user_id=USERS.user_id
					LEFT JOIN
					(
						SELECT substr('JanFebMarAprMayJunJulAugSepOctNovDec', 1 + 3*strftime('%mm', date(sent_timestamp)), -3) || strftime("-%d %H:%M", sent_timestamp) AS last_timestamp, msg AS last_msg, MESSAGES.conversation_id
						FROM MESSAGES
						INNER JOIN
						(
							SELECT max(msg_id) as msg_id, conversation_id
							FROM MESSAGES GROUP BY conversation_id
						) AS M1 ON M1.msg_id=MESSAGES.msg_id AND M1.conversation_id=MESSAGES.conversation_id
					) AS LAST_MSG ON LAST_MSG.conversation_id = CONVERSATIONS_USERS.conversation_id
					WHERE CONVERSATIONS_USERS.conversation_id IN (
						SELECT CONVERSATIONS_USERS.conversation_id FROM CONVERSATIONS_USERS
						LEFT JOIN GROUPS ON CONVERSATIONS_USERS.conversation_id=GROUPS.conversation_id
						WHERE user_id = ? AND GROUPS.conversation_id IS NULL
					) AND CONVERSATIONS_USERS.user_id <> ?
					UNION ALL
					SELECT GROUPS.conversation_id, GROUPS.group_id as user_or_group_id, group_name as user_or_group_name, ifnull(group_photo, '') as user_or_group_photo, 1 is_group,
						ifnull(last_timestamp, '') as last_timestamp, ifnull(last_msg, '') as last_msg
					FROM GROUPS
					LEFT JOIN
					(
						SELECT substr('JanFebMarAprMayJunJulAugSepOctNovDec', 1 + 3*strftime('%mm', date(sent_timestamp)), -3) || strftime("-%d %H:%M", sent_timestamp) AS last_timestamp, msg AS last_msg, MESSAGES.conversation_id
						FROM MESSAGES
						INNER JOIN
						(
							SELECT max(msg_id) as msg_id, conversation_id
							FROM MESSAGES GROUP BY conversation_id
						) AS M1 ON M1.msg_id=MESSAGES.msg_id AND M1.conversation_id=MESSAGES.conversation_id
					) AS LAST_MSG ON LAST_MSG.conversation_id = GROUPS.conversation_id
					WHERE GROUPS.conversation_id IN (
					SELECT CONVERSATIONS_USERS.conversation_id FROM CONVERSATIONS_USERS WHERE user_id = ?)
					;`, from_user_id, from_user_id, from_user_id)

	/*
		SELECT conversation_id, user_name as user_or_group_name, ifnull(user_photo, '') as user_or_group_photo, 0 is_group
		FROM CONVERSATIONS_USERS
		INNER JOIN USERS ON CONVERSATIONS_USERS.user_id=USERS.user_id
		WHERE CONVERSATIONS_USERS.conversation_id IN (
			SELECT CONVERSATIONS_USERS.conversation_id FROM CONVERSATIONS_USERS
			LEFT JOIN GROUPS ON CONVERSATIONS_USERS.conversation_id=GROUPS.conversation_id
			WHERE user_id = ? AND GROUPS.conversation_id IS NULL
		) AND CONVERSATIONS_USERS.user_id <> ?
		UNION ALL
		SELECT GROUPS.conversation_id, group_name as user_or_group_name, ifnull(group_photo, '') as user_or_group_photo, 1 is_group
		FROM GROUPS
		WHERE GROUPS.conversation_id IN (
			SELECT CONVERSATIONS_USERS.conversation_id FROM CONVERSATIONS_USERS WHERE user_id = ?
		);`, from_user_id, from_user_id, from_user_id)
	*/

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.Conversation_id, &conv.User_or_group_id, &conv.User_or_group_name,
			&conv.User_or_group_photo, &conv.Is_group, &conv.Last_timestamp, &conv.Last_msg)
		if err != nil {
			return conversations, err
		}
		conversations = append(conversations, conv)
	}
	err = rows.Err()
	if err != nil {
		return conversations, err
	}

	return conversations, nil
}
