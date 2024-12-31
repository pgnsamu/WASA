package database

import "errors"

func (db *appdbimpl) DeleteMessage(idConversation int, idUser int, idMessageToDelete int) error {
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
	err = stmt.QueryRow(idMessageToDelete).Scan(&sentBy)
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

	queryStr = "DELETE FROM received WHERE messageId = ?"
	stmt, err = tx.Prepare(queryStr)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	defer stmt.Close()

	// Execute the deletion
	res, err := stmt.Exec(idMessageToDelete)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}

	// Check how many rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no message found with that id")
	}

	// cancellazione di tutti i commenti collegati al messaggio da eliminare
	queryStr = "DELETE FROM messages WHERE answerTo = ?"
	stmt, err = tx.Prepare(queryStr)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	defer stmt.Close()

	// Execute the deletion
	_, err = stmt.Exec(idMessageToDelete)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}

	queryStr = "DELETE FROM messages WHERE id = ?"
	stmt, err = tx.Prepare(queryStr)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	defer stmt.Close()

	// Execute the deletion
	res, err = stmt.Exec(idMessageToDelete)
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}

	// Check how many rows were affected
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		err2 := tx.Rollback() // Rollback in caso di errore
		if err2 != nil {
			return err2
		}
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no message found with that id")
	}

	// Conferma la transazione
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
