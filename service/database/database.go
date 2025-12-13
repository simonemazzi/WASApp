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
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	CreateSession(string) (string, string, time.Time, error)
	UserByToken(token string) (int, error)
	Users() ([]DBUser, error)
	UsersBySearch(token string) ([]DBUser, error)
	SetMyUserName(int, string) error
	SetMyPhoto(url string, width int, height int, mime string, userId int) error
	IDExists(token int) bool
	SetGroupPhoto(url string, width int, height int, mime string, groupId int) error
	GetConversations(userId int) ([]Conversation, error)
	CreateConversation(userId int, username string, time string) (Conversation, error)
	SearchUserByUsername(username string, time string) (int, error)
	GetPhoto(userId int, time string) (Avatar, error)
	GetConversationById(conversationId int) (Conversation, error)
	UserConversation(userId int, conversationId int) (bool, error)
	GetMessages(conversationId int, userId int) ([]Message, error)
	InsertPhoto(url string, width int, height int, mime string) (int, error)
	InsertMessage(conversationId int, userId int, text string, photoId *int) error
	GetUsername(userId int, time string) (string, error)
	IsRead(messageId int, userId int) (string, error)
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

	impl := &appdbimpl{c: db}

	if err := impl.initSchema(); err != nil {
		return nil, err
	}

	return impl, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
