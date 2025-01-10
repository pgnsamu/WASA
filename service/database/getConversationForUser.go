package database

import (
	"errors"
)

func (db *appdbimpl) GetConversationForUser(idUser int) (*[]Conversation, error) {

	query := "SELECT id, name, createdAt, isGroup, photo, description FROM conversations as c, participate as p WHERE c.id = p.conversationId and userId = ?"
	rows, err := db.c.Query(query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // defer keyword is used to schedule a function call to be executed after the surrounding function has completed. Deferred functions are typically used for cleanup tasks

	var convFinale []Conversation
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.Id, &conv.Name, &conv.CreatedAt, &conv.IsGroup, &conv.Photo, &conv.Description)
		if err != nil {
			return nil, err
		}
		convFinale = append(convFinale, conv)
	}
	// Check for errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If no conversations were found, return a custom error
	if len(convFinale) == 0 {
		return nil, errors.New("no conversations found")
	}
	return &convFinale, nil
}
