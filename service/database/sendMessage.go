package database

import (
	"errors"
	"time"
)

// TODO: messagetype per ora sul doc è enum invece di bool
// TODO: togliere il previewContent dal doc
// TODO: forse l'attributo replyTo su sendmessage è inutile

// errori ritornabili da SendMessage
// il messaggio non appartiene a questa conversazione

func (db *appdbimpl) SendMessage(idConversation int, idUser int, content string, photoContent []byte, replyTo *int, isForwarded int) (*[]Message, error) {
	// Controlla se l'utente appartiene alla conversazione
	var repl int
	if replyTo != nil {
		repl = *replyTo
		ex, err := db.MessageBelongsToConversation(repl, idConversation)
		if err != nil {
			return nil, err
		}
		if !ex {
			return nil, errors.New("il messaggio non appartiene a questa conversazione")
		}
	} else {
		repl = -1
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

	// Esegui la query con l'ID della conversazione (ad esempio, 20)
	rows, err := stmt.Query(idConversation)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}
	defer rows.Close()

	// Dichiarare una slice per contenere gli ID degli utenti
	var userIDs []int

	// Itera sulle righe e aggiungi gli ID degli utenti alla slice
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
		userIDs = append(userIDs, userID) // Aggiungi l'ID utente alla slice
	}

	// Controlla se ci sono errori dopo l'iterazione
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
		if id != idUser {
			_, err := stmt.Exec(id, lastInsertId, 0) // 0 = non letto
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
