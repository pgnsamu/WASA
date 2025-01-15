package database

import (
	"errors"
)

/*
	func (db *appdbimpl) GetUserInfo(id int) (User, error) {
		var user User
		query := "SELECT id, username, name, surname FROM users WHERE id=$1"
		err := db.c.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Name, &user.Surname)
		if err != nil {
			if err == sql.ErrNoRows {
				// No rows foundm
				return User{}, errors.New("userNotFound")
			}
			// Other errors
			log.Fatal("Failed to read record: ", err)
			return User{}, err
		}
		return user, nil
	}
*/

// ritornare l'utente aggiornato
func (db *appdbimpl) SetMyUserName(id int, username string) (*User, error) {

	// preparazione della query
	stmt, err := db.c.Prepare("UPDATE users SET username = ? WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// TODO: aggiungere controllo username già esistente
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
