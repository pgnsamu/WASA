package database

import (
	"errors"
)

// ritornare l'utente aggiornato
func (db *appdbimpl) SetMyUserName(id int, username string) (*User, error) {

	// controllo che l'username non sia già esistente
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username già esistente")
	}

	// preparazione della query
	stmt, err := db.c.Prepare("UPDATE users SET username = ? WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// esecuzione
	res, err := stmt.Exec(username, id)
	if err != nil {
		return nil, err
	}

	// controllo righe interessate che in questo caso dovranno essere == 1
	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 1 {
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New("utente non trovato")
		}
	}

	var user *User
	user, err = db.GetUserInfo(id)
	if err != nil {
		return nil, err

	}

	return user, nil
}
