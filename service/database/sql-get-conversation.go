package database

import "fmt"

// Retrieve the conversations
func (db *appdbimpl) GetConversation(conversation_id int) ([]Message, error) {
	// Query data
	var sql = `SELECT MESSAGES.msg_id, ifnull(forwarded_from_msg_id, 0), conversation_id, MESSAGES.from_user_id, U1.user_name from_user_name,
				strftime("%H:%M", sent_timestamp) AS timestamp, msg, iif(substr(msg, 1, 11)='data:image/', true, false) AS is_photo,
				ifnull(R1.reaction_user, "") AS reaction, ifnull(S1.msg_id, 0) == 0 AS seen
				FROM MESSAGES
				INNER JOIN USERS AS U1 ON U1.user_id = MESSAGES.from_user_id
				LEFT JOIN
					(
						SELECT msg_id, from_user_id, group_concat(reaction || USERS.user_name) as reaction_user
						FROM REACTIONS
						INNER JOIN USERS ON USERS.user_id = REACTIONS.from_user_id
						GROUP BY msg_id
					) AS R1
					ON R1.msg_id=MESSAGES.msg_id
				LEFT JOIN USERS ON USERS.user_id = R1.from_user_id
				LEFT JOIN
					(
						SELECT MESSAGES.msg_id
						FROM MESSAGES
						LEFT JOIN CONVERSATIONS_USERS ON CONVERSATIONS_USERS.conversation_id = MESSAGES.conversation_id
						LEFT JOIN MESSAGES_SEEN ON MESSAGES_SEEN.msg_id = MESSAGES.msg_id AND MESSAGES_SEEN.seen_by_user_id = CONVERSATIONS_USERS.user_id
						WHERE MESSAGES.from_user_id != CONVERSATIONS_USERS.user_id AND MESSAGES_SEEN.seen_by_user_id is null
						GROUP BY MESSAGES.msg_id
					) AS S1 ON MESSAGES.msg_id = S1.msg_id
				WHERE conversation_id = ` + fmt.Sprint(conversation_id) + ` ORDER BY sent_timestamp DESC`

	rows, err := db.c.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Msg_id, &msg.Forwarded_from_msg_id, &msg.Conversation_id, &msg.From_user_id,
			&msg.From_user_name, &msg.Sent_timestamp, &msg.Msg, &msg.Is_photo, &msg.Reaction, &msg.Seen)
		if err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}
	err = rows.Err()
	if err != nil {
		return messages, err
	}

	return messages, nil
}
