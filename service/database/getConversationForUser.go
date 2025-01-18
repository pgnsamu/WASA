package database

// nessun errore ritornabile

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
			c.description, (
				SELECT m.content
				FROM messages m
				WHERE m.conversationId = c.id
				ORDER BY m.id DESC
				LIMIT 1
			) as lastMessage
		FROM conversations as c
		JOIN participate as p ON c.id = p.conversationId
		JOIN users as u ON p.userId = u.id
		WHERE p.userId = ?;`

	rows, err := db.c.Query(query2, idUser, idUser, idUser, idUser, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convFinale []Conversation
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.Id, &conv.Name, &conv.CreatedAt, &conv.IsGroup, &conv.Photo, &conv.Description, &conv.LastMessage)
		if err != nil {
			return nil, err
		}
		convFinale = append(convFinale, conv)
	}
	// Controlla se ci sono stati errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &convFinale, nil
}
