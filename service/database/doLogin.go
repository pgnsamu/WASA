package database

import (
	"database/sql"
	"errors"
)

type User struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Photo    *[]byte `json:"photo,omitempty"`
}

// nessun errore ritornabile

func (db *appdbimpl) DoLogin(username string) (*int, error) {
	// ricerca se l'username è già nel DB
	resSearch, errore := db.SearchUser(username)
	// se la ricerca non ha trovato errori e ha ritornato -1 allora il nome utente non esiste e va creato un nuovo utente
	if resSearch == -1 && errore == nil {
		res, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
		if err != nil {
			return nil, err
		}
		// ritorno l'id dell'utente appena creato
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		// conversione int64 to int
		id := int(lastInsertId)
		return &id, nil

	} else if resSearch == -1 && errore != nil { // se la ricerca ha trovato errori e quindi non ha trovato il record allora ritorno errore
		return nil, errore
	} else {
		return &resSearch, nil
	}
}

// ritorna -1 se non c'è gia, altrimenti ritorna l'id dell'utente
func (db *appdbimpl) SearchUser(username string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows found
			return -1, nil
		}
		// Other errors
		return -1, err
	}

	return id, nil
}
