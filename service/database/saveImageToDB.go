package database

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// errori ritornabili da SaveImageToDB
// id non trovato

func (db *appdbimpl) SaveImageToDB(imgData []byte, table string, field string, userId int) error {

	// query custom per essere fatta su varie tabelle e campi
	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?;", table, field)
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Esegui l'istruzione, passando i dati dell'immagine come BLOB
	result, err := stmt.Exec(imgData, userId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("id non trovato")
	} else {
		return nil
	}
}
