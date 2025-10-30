package database

import (
	"fmt"
)

// Retrieve a message
func (db *appdbimpl) GetMessage(msg_id int) (Message, error) {
	var msg Message
	var err error

	// Query data
	var sql = `SELECT msg_id, ifnull(forwarded_from_msg_id, 0), conversation_id, from_user_id,
				sent_timestamp, seen, msg FROM MESSAGES
				WHERE msg_id = ` + fmt.Sprint(msg_id)

	row := db.c.QueryRow(sql)

	err = row.Scan(&msg.Msg_id, &msg.Forwarded_from_msg_id, &msg.Conversation_id,
		&msg.From_user_id, &msg.Sent_timestamp, &msg.Seen, &msg.Msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
