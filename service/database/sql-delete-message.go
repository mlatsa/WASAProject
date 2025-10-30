package database

import "database/sql"

// Set user name
func (db *appdbimpl) DeleteMessage(msg_id int) (sql.Result, error) {
	var count sql.Result
	var err error
	count, err = db.c.Exec(`DELETE FROM MESSAGES WHERE msg_id = ?`, msg_id)
	return count, err
}
