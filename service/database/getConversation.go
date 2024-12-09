package database

import (
	"errors"
	"log"
)

func (db *appdbimpl) GetConversation(id int) (*Conversation, error) {
	query := "SELECT id, name, createdAt, isGroup, description, photo FROM conversations WHERE id=$1"
	rows, err := db.c.Query(query, id)
	if err != nil {
		log.Fatal("Error executing query:", err)
		return nil, err
	}
	defer rows.Close() // defer keyword is used to schedule a function call to be executed after the surrounding function has completed. Deferred functions are typically used for cleanup tasks

	var conversation Conversation
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&conversation.Id, &conversation.Name, &conversation.CreatedAt, &conversation.IsGroup, &conversation.Description, &conversation.Photo)
		if err != nil {
			log.Fatal("Error scanning row:", err)
			return nil, err
		}
	} else {
		errorUserNotFound := errors.New("user not found")
		return nil, errorUserNotFound
	}
	return &conversation, nil
}
