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
	// GetName() (string, error)
	// SetName(name string) error
	SetMyUserName(id int, username string) (*User, error)
	GetUsersDB() (*[]User, error)
	GetUserInfo(id int) (*User, error)
	GetConversationInfo(idConversation int, idUser int) (*Conversation, error)
	NewConversation(userId int, name string, isGroup bool, photo *[]byte, description *string, partecipantsId []int) (*Conversation, error)
	GetUsersOfConversation(idConversation int, idUser int) (*[]User, error)
	DeleteUserFromConv(idConversation int, idUser int, idUserToDelete int) (*[]User, error)
	AddToGroup(idConversation int, idUser int, idUserToAdd int) (*[]User, error)
	GetConversationForUser(idUser int) (*[]Conversation, error)
	SetGroupName(idUser int, idConversation int, name string) (*Conversation, error)
	GetMessagesFromConversation(conversationID int) (*[]Message, error)
	SendMessage(idConversation int, idUser int, content string, photoContent []byte, messageType bool, replyTo *int, isForwarded int) (*[]Message, error)
	ForwardMessage(idConversationSource int, idConversationDest int, idUser int, idMessage int) (*Conversation, error)
	DeleteMessage(idConversation int, idUser int, idMessageToDelete int) error
	UncommentMessage(idConversation int, idUser int, idMessageToUncomment int) error

	SaveImageToDB(imgData []byte, table string, field string, userId int) error

	SearchUser(username string) (int, error)
	DoLogin(username string, name string, surname string) (*int, error)

	UserExist(idConv int, idUser int) (bool, error)
	SeeAllMessages(idConv int, idUser int) (int, error)
	ReceiveAllMessages(idUser int) (int, error)
	GetProfilePhoto(id int) ([]byte, error)
	GetUserId(username string) (*int, error) // non viene usata per ora
	IsCommentTo(idComment int, idMessage int, idConversation int) (bool, error)
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

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='participate';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			CREATE TABLE participate (
				userId INTEGER NOT NULL,                -- Foreign key to the user table
				conversationId INTEGER NOT NULL,        -- Foreign key to the conversation table
				PRIMARY KEY (userId, conversationId),   -- Composite primary key
				FOREIGN KEY (userId) REFERENCES users(id),
				FOREIGN KEY (conversationId) REFERENCES conversations(id)
			)`
		_, err = db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='received';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			CREATE TABLE received (
			userId INTEGER NOT NULL,
			messageId INTEGER NOT NULL,
			status INTEGER NOT NULL,
			CHECK (status IN (0, 1, 2)),
			PRIMARY KEY (userId, messageId),
			FOREIGN KEY (userId) REFERENCES users(id),
			FOREIGN KEY (messageId) REFERENCES messages(id)
		)`
		_, err = db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversations';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			CREATE TABLE conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for the conversation
			name TEXT NOT NULL,                         -- Conversation name (e.g., group chat name)
			createdAt INTEGER NOT NULL,                 -- Timestamp of when the conversation was created (Unix time)
			isGroup BOOLEAN NOT NULL,                   -- Indicates if it's a group chat
			photo BLOB,                                 -- Photo for the conversation (e.g., group avatar)
			description TEXT                            -- Optional description field
			CHECK (
				( isGroup = TRUE and description is not null and photo is not null )
				or 
				( isGroup = FALSE and description is null and photo is null)
			
			)
		)`
		_, err = db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			CREATE TABLE users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for the user
				username TEXT NOT NULL,                     -- Username, must be unique (consider adding UNIQUE constraint)
				name TEXT,                                  -- First name of the user
				surname TEXT,                               -- Last name of the user
				photo BLOB                                  -- Photo stored as a binary large object
			)
		`
		_, err = db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			CREATE TABLE "messages" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT,
				photoContent BLOB,
				sentAt INTEGER NOT NULL,
				conversationId INTEGER NOT NULL,
				answerTo INTEGER,
				isForwarded INTEGER,
				senderId INTEGER NOT NULL,
				CHECK (answerTo <> id),
				FOREIGN KEY (answerTo) REFERENCES messages(id),
				FOREIGN KEY (conversationId) REFERENCES conversations(id),
				FOREIGN KEY (senderId) REFERENCES users(id)
				)
			`
		_, err = db.Exec(query)
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
