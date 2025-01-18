package database

import (
	"errors"
)

// errori ritornabili da GetConversationInfo
// utente non registrato
// utente non trovato

func (db *appdbimpl) GetConversationInfo(idConversation int, idUser int) (*Conversation, error) {

	// controllo utente se è all'interno
	resu, err := db.UserExist(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	if !resu {
		return nil, errors.New("utente non registrato")
	}

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
		WHERE c.id = ?;`

	rows, err := db.c.Query(query2, idUser, idUser, idUser, idUser, idConversation)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // defer keyword is used to schedule a function call to be executed after the surrounding function has completed. Deferred functions are typically used for cleanup tasks

	var conversation Conversation
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&conversation.Id, &conversation.Name, &conversation.CreatedAt, &conversation.IsGroup, &conversation.Photo, &conversation.Description)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("utente non trovato")
	}
	// TODO: capire se è qua il problema stessa cosa su getUserInfo riga 44
	// Check for errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &conversation, nil
}
