package database

import (
	"errors"
)

func (db *appdbimpl) GetConversationForUser(idUser int) (*[]Conversation, error) {

	// TODO: capire se fare la stessa cosa anche con il nome
	query2 := `
		SELECT c.id,
			CASE 
				WHEN c.isGroup = false
				THEN (
					SELECT u.username
					FROM conversations c2
					JOIN participate p1 ON c2.id = p1.conversationId
					JOIN participate p2 ON c2.id = p2.conversationId
					JOIN users u ON p2.userId = u.id
					WHERE p1.userId = ?
					AND p2.userId != ? 
					AND c2.id = c.id
					AND c2.isGroup = false
				)
				ELSE c.name
			END as name, 
		  c.createdAt, c.isGroup, 
			CASE 
				WHEN c.isGroup = false
				THEN (
					SELECT u.photo
					FROM conversations c2
					JOIN participate p1 ON c2.id = p1.conversationId
					JOIN participate p2 ON c2.id = p2.conversationId
					JOIN users u ON p2.userId = u.id
					WHERE p1.userId = ?
					AND p2.userId != ? 
					AND c2.id = c.id
					AND c2.isGroup = false
				)
				ELSE c.photo
			END as photo, 
			c.description 
		FROM conversations as c
		JOIN participate as p ON c.id = p.conversationId
		JOIN users as u ON p.userId = u.id
		WHERE p.userId = ?;`

	// query := "SELECT id, name, createdAt, isGroup, photo, description FROM conversations as c, participate as p WHERE c.id = p.conversationId and userId = ?"
	rows, err := db.c.Query(query2, idUser, idUser, idUser, idUser, idUser)
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
