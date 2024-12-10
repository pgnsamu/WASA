package database

import (
	"errors"
	"fmt"
	"log"
)

func (db *appdbimpl) GetConversationInfo(idConversation int, idUser int) (*Conversation, error) {

	fmt.Println(idConversation, idUser)
	// controllo utente se è all'interno
	stmt, err := db.c.Prepare("SELECT userId FROM participate as p WHERE p.conversationId = ? and p.userId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Count the number of rows
	count := 0
	for rows.Next() {
		count++
	}

	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if count != 1 {
		return nil, errors.New("utente non registrato nel gruppo")
	}

	query := "SELECT id, name, createdAt, isGroup, description, photo FROM conversations WHERE id=?"
	rows, err = db.c.Query(query, idConversation)
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
