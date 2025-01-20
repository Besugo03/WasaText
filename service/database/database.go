/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	CheckUsersTable() error
	GetUserID(username string) (string, error)
	GetUserName(userid string) (string, error)
	RegisterUser(username string, userid string, photo string) error // userid is basically its auth token and identifier.
	UpdateUser(username string, userid string, photo string) error   // update the user's photo. Arguments can be nil to keep the old value.
	GetAllUsers() ([]string, error)
	CheckUserExistence(userid string) (bool, error)

	CheckConversationsTable() error
	CreateConversation(isGroup bool, name string, picture string) (string, error)
	UpdateConversation(conversationID string, isGroup bool, name string, picture string) error
	DeleteConversation(conversationID string) error
	GetConversationInfo(conversationID string, userID string) (string, bool, string, string, error)
	GetUserConversations(userID string) ([]string, error)

	AddParticipant(conversationID string, userID string) error
	RemoveParticipant(conversationID string, userID string) error
	GetParticipants(conversationID string) ([]string, error)

	SendMessage(conversationID string, senderID string, content string, isForwarded bool) (string, error)
	GetMessages(conversationID string) (messageIDs []string, conversationIDs []string, senderIDs []string, contents []string, timestmaps []string, err error)
	GetLastMessage(conversationID string) (string, string, string, string, error)
	DeleteMessage(messageID string) error
	GetMessageInfo(messageID string, userID string) (string, string, string, string, string, bool, error)

	CheckCommentsTable() error
	AddComment(messageID string, senderID string, content string) error
	GetComments(messageID string) (commentIDs []string, messageIDs []string, senderIDs []string, contents []string, timestamps []string, err error)
	DeleteComment(commentID string) error
	GetCommentInfo(commentID string) (string, string, string, string, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
