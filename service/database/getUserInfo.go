package database

import (
	"errors"
)

func (db *appdbimpl) GetUserInfo(id int) (*User, error) {
	query := "SELECT id, username, photo FROM users WHERE id= ?"
	rows, err := db.c.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // defer keyword is used to schedule a function call to be executed after the surrounding function has completed. Deferred functions are typically used for cleanup tasks

	var user User
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("user not found")
	}
	// Check for errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
