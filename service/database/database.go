/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
*/

package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	User_id    int    `json:"user_id"`
	User_name  string `json:"user_name"`
	User_photo string `json:"user_photo"`
}

type Group struct {
	Group_id        int    `json:"group_id"`
	Group_name      string `json:"group_name"`
	Group_photo     string `json:"group_photo"`
	Conversation_id int    `json:"conversation_id"`
}

type Message struct {
	Msg_id                int    `json:"msg_id"`
	Forwarded_from_msg_id int    `json:"forwarded_from_msg_id"`
	Conversation_id       int    `json:"conversation_id"`
	From_user_id          int    `json:"from_user_id"`
	From_user_name        string `json:"from_user_name"`
	Sent_timestamp        string `json:"sent_timestamp"`
	Msg                   string `json:"msg"`
	Is_photo              bool   `json:"is_photo"`
	Reaction              string `json:"reaction"`
	Seen                  int    `json:"seen"`
}

type Conversation struct {
	Conversation_id     int    `json:"conversation_id"`
	User_or_group_id    int    `json:"user_or_group_id"`
	User_or_group_name  string `json:"user_or_group_name"`
	User_or_group_photo string `json:"user_or_group_photo"`
	Is_group            int    `json:"is_group"`
	Last_timestamp      string `json:"last_timestamp"`
	Last_msg            string `json:"last_msg"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error
	GetUserId(user_name string) (int, error)
	GetUserIdAndPhoto(user_name string) (int, string, error)
	GetNextUserId() (int, error)
	GetAllUsers() ([]User, error)
	CreateUser(user_id int, user_name string) error
	SetMyUserName(id int, name string) error
	SetMyPhoto(id int, photo string) error
	GetNextConversationId() (int, error)
	GetNextMessageId() (int, error)
	GetNextCommentId() (int, error)
	GetMessage(msg_id int) (Message, error)
	GetMessageText(msg_id int) (string, error)
	ChangeName(id int, new_name string) error
	GetMyConversations(from_user_id int) ([]Conversation, error)
	GetConversation(conversation_id int) ([]Message, error)
	MarkMessagesSeen(seen_by_user_id int, conversation_id int) error
	IsUserInConversation(user_id int, conversation_id int) (bool, error)
	GetConversationId(user_id1 int, user_id2 int) (int, error)
	CreateConversation(conversation_id int, user_id1 int, user_id2 int) error
	AddUserToConversation(conversation_id int, user_id int) error
	SendMessage(conversation_id int, from_user_id int, message string, forwarded_from_msg_id int) error
	DeleteMessage(msg_id int) (sql.Result, error)
	CommentMessage(msg_id int, from_user_id int, reaction string) (int, error)
	UncommentMessage(msg_id int, from_user_id int) error
	CreateGroup(int, string, string, int) error
	GetMyGroups(int) ([]Group, error)
	GetGroupDataFromName(group_name string) (int, string, int, error)
	GetNextGroupId() (int, error)
	SetGroupName(group_id int, group_name string) error
	SetGroupPhoto(group_id int, photo string) error
	LeaveGroup(user_id int, group_id int) error
	GetMembers(group_id int) ([]User, error)
}

type appdbimpl struct {
	c *sql.DB
}

// Returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Check if table exists. If not, the database is empty, so we need to create the structure
	var tableName string
	var err error

	// Check and create table USERS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='USERS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE USERS (
					user_id	INTEGER NOT NULL,
					user_name TEXT NOT NULL UNIQUE,
					user_photo TEXT,
					PRIMARY KEY(user_id));`)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table MESSAGES
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='MESSAGES';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE MESSAGES (
					msg_id INTEGER NOT NULL,
					conversation_id INTEGER NOT NULL,
					from_user_id INTEGER NOT NULL,
					sent_timestamp TEXT NOT NULL DEFAULT '1970-01-01 00:00:00',
					msg TEXT NOT NULL,
					forwarded_from_msg_id INTEGER,
					PRIMARY KEY(msg_id));`)

		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table MESSAGES_SEEN
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='MESSAGES_SEEN';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE "MESSAGES_SEEN" (
					msg_id INTEGER NOT NULL,
					seen_by_user_id INTEGER NOT NULL,
					PRIMARY KEY(msg_id, seen_by_user_id));`)

		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table REACTIONS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='REACTIONS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE "REACTIONS" (
					"reaction_id"	INTEGER NOT NULL,
					"from_user_id"	INTEGER NOT NULL,
					"msg_id"	INTEGER NOT NULL,
					"reaction"	TEXT NOT NULL,
					PRIMARY KEY("reaction_id"));`)

		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table CONVERSATIONS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='CONVERSATIONS_USERS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE CONVERSATIONS_USERS (
					conversation_id INTEGER NOT NULL,
					user_id INTEGER NOT NULL,
					PRIMARY KEY("conversation_id","user_id"));`)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table GROUPS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='GROUPS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE GROUPS (
					group_id INTEGER NOT NULL,
					group_name INTEGER NOT NULL,
					group_photo INTEGER,
					conversation_id INTEGER NOT NULL,
					PRIMARY KEY("group_id"));`)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
