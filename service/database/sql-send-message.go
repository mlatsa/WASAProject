package database

import "time"

// Send message
func (db *appdbimpl) SendMessage(conversation_id int, from_user_id int, message string, forwarded_from_msg_id int) error {
	var msg_id int
	var err error

	msg_id, _ = db.GetNextMessageId()

	_, err = db.c.Exec(`INSERT INTO MESSAGES
		(msg_id, conversation_id, from_user_id, sent_timestamp, msg, forwarded_from_msg_id)
		VALUES (?, ?, ?, ?, ?, IIF(? == 0, NULL, ?))`,
		msg_id, conversation_id, from_user_id, time.Now(), message, forwarded_from_msg_id, forwarded_from_msg_id)

	return err
}
