package database

import (
	"database/sql"
	"log"
)

func (db *appdbimpl) GetUsersDB() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT nome FROM utenti WHERE id=0").Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found
			return "", nil
		}
		// Other errors
		log.Fatal("Failed to read record: ", err)
		return "", err
	}

	return name, nil
}
