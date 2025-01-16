package database

import (
	"errors"
	"time"
)

// TODO: messagetype per ora sul doc è enum invece di bool
// TODO: togliere il previewContent dal doc
// TODO: forse l'attributo replyTo su sendmessage è inutile
func (db *appdbimpl) SendMessage(idConversation int, idUser int, content string, photoContent []byte, replyTo *int, isForwarded int) (*[]Message, error) {
	resu, err := db.UserExist(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	if !resu {
		return nil, errors.New("utente non presente nella chat")
	}

	// Inizia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}

	queryStr := `
	INSERT INTO messages (content, photoContent, sentAt, conversationId, isForwarded, senderId, answerTo)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	// Prepara la query di INSERT
	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	defer stmt.Close()

	// Ottieni il timestamp corrente
	sentAt := time.Now().UnixMilli()

	repl := -1
	if replyTo != nil {
		repl = *replyTo
	}

	// Esegui la query
	resul, err := stmt.Exec(content, photoContent, sentAt, idConversation, isForwarded, idUser, repl)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	lastInsertId, err := resul.LastInsertId()
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	// get tutti gli utenti del partecipant eliminando se stesso
	queryStr = `
		SELECT u.id
		FROM participate as p, users as u
		WHERE p.userId = u.id and p.conversationId = ?
		`

	stmt, err = tx.Prepare(queryStr)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with the conversation ID (e.g., 20)
	rows, err := stmt.Query(idConversation) // Passing 20 as the conversationId
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	defer rows.Close()

	// Declare a slice to hold the user IDs
	var userIDs []int

	// Iterate over the rows and append the user IDs to the slice
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			err2 := tx.Rollback() // Rollback in caso di errore
			if err2 != nil {
				return nil, err2
			}
			return nil, err
		}
		userIDs = append(userIDs, userID) // Append the user ID to the slice
	}

	// Check for any errors after iterating
	if err := rows.Err(); err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	// Prepara la query di INSERT
	stmt, err = tx.Prepare("INSERT INTO received (userId, messageId, status) VALUES (?, ?, ?)")
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	defer stmt.Close()

	// per ogni utente appartenente alla conversation dove è stato inviato il messaggio aggiungere una riga di insert in received
	for _, id := range userIDs {
		// fmt.Println(id, idUser)
		if id != idUser {
			_, err := stmt.Exec(id, lastInsertId, 0) // TODO: cercare se nel programma ci sono ancora "delivered"
			if err != nil {
				err2 := tx.Rollback() // Rollback in caso di errore
				if err2 != nil {
					return nil, err2
				}
				return nil, err
			}
		}
	}

	// Conferma la transazione
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	res, err := db.GetMessagesFromConversation(idConversation)
	if err != nil {
		return nil, err
	}
	return res, nil

}
