package database

import "errors"

// errori ritornabili da UncommentMessage
// gruppo non trovato
// autore del messaggio sbagliato
// messaggio non trovato

func (db *appdbimpl) UncommentMessage(idConversation int, idUser int, idMessageToUncomment int) error {
	exist, err := db.UserExist(idConversation, idUser)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("gruppo non trovato")
	}

	// Inizia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// controllo se il messaggio è stato mandato da chi vuole eliminarlo
	queryStr := "SELECT senderId FROM messages as m WHERE m.id = ?"
	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	defer stmt.Close()

	var sentBy int
	err = stmt.QueryRow(idMessageToUncomment).Scan(&sentBy)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	if sentBy != idUser {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return errors.New("autore del messaggio sbagliato")
	}

	// iniziamo
	updateQuery := "UPDATE messages SET answerTo = -1 WHERE id = ?"
	updateStmt, err := tx.Prepare(updateQuery)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	defer updateStmt.Close()

	result, err := updateStmt.Exec(idMessageToUncomment)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	if rows, _ := result.RowsAffected(); rows != 1 {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return errors.New("messaggio non trovato")
	}

	// Conferma la transazione
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
