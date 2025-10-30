package database

import (
	"fmt"
)

// Retrieve a message
func (db *appdbimpl) GetMessageText(msg_id int) (string, error) {
	var msg string
	err := db.c.QueryRow(`SELECT msg FROM MESSAGES WHERE msg_id = ` + fmt.Sprint(msg_id)).Scan(&msg)
	return msg, err
}
