package database

import (
	"errors"
	"log"
	"time"
)

// TODO: messagetype per ora sul doc è enum invece di bool
// TODO: togliere il previewContent dal doc
// TODO: forse l'attributo replyTo su sendmessage è inutile
func (db *appdbimpl) SendMessage(idConversation int, idUser int, content string, photoContent []byte, messageType bool, replyTo *int) (*[]Message, error) {
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
		log.Fatal(err)
	}

	queryStr := `
	INSERT INTO messages (content, photoContent, sentAt, conversationId)
	VALUES (?, ?, ?, ?)
	`
	// Prepara la query di INSERT
	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		tx.Rollback() // Rollback in caso di errore
		log.Fatal(err)
	}
	defer stmt.Close()

	// Ottieni il timestamp corrente
	sentAt := time.Now().UnixMilli()

	// Esegui la query
	resul, err := stmt.Exec(content, photoContent, sentAt, idConversation)
	if err != nil {
		tx.Rollback() // Rollback in caso di errore
		return nil, err
	}

	lastInsertId, err := resul.LastInsertId()
	if err != nil {
		tx.Rollback() // Rollback in caso di errore
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
		tx.Rollback() // Rollback in caso di errore
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with the conversation ID (e.g., 20)
	rows, err := stmt.Query(idConversation) // Passing 20 as the conversationId
	if err != nil {
		tx.Rollback() // Rollback in caso di errore
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
			tx.Rollback() // Rollback in caso di errore
			return nil, err
		}
		userIDs = append(userIDs, userID) // Append the user ID to the slice
	}

	// Check for any errors after iterating
	if err := rows.Err(); err != nil {
		tx.Rollback() // Rollback in caso di errore
		return nil, err
	}

	// Prepara la query di INSERT
	stmt, err = tx.Prepare("INSERT INTO sent (userId, messageId, status, answerTo) VALUES (?, ?, ?, ?)") //TODO: forse meglio usare il sent? BOH
	if err != nil {
		tx.Rollback() // Rollback in caso di errore
		return nil, err
	}
	defer stmt.Close()

	repl := -1
	if replyTo != nil {
		repl = *replyTo
	}

	// per ogni utente appartenente alla conversation dove è stato inviato il messaggio aggiungere una riga di insert in sent
	for _, id := range userIDs {
		_, err := stmt.Exec(id, lastInsertId, "delivered", repl)
		if err != nil {
			tx.Rollback() // Rollback in caso di errore
			return nil, err
		}
	}

	// Conferma la transazione
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	res, err := db.GetMessagesByConversation(idConversation)
	if err != nil {
		return nil, err
	}
	return res, nil

}
