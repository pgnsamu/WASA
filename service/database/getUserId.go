package database

import (
	"errors"
)

func (db *appdbimpl) GetUserId(username string) (*int, error) {
	query := "SELECT id FROM users WHERE username = ?"
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // defer keyword is used to schedule a function call to be executed after the surrounding function has completed. Deferred functions are typically used for cleanup tasks

	var idUser int
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&idUser)
		if err != nil {
			return nil, err
		}
	} else {
		errorUserNotFound := errors.New("user not found")
		return nil, errorUserNotFound
	}
	// Check for errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &idUser, nil
}
