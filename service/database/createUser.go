package database

import (
	"database/sql"
	"errors"
	"log"
)

type User struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Name     *string `json:"name,omitempty"`
	Surname  *string `json:"surname,omitempty"`
	Photo    *[]byte `json:"photo,omitempty"`
}

func (db *appdbimpl) CreateUser(username string, name string, surname string) (string, error) {
	// ricerca se l'username è già nel DB
	resSearch, errore := db.SearchUser(username)

	// se la ricerca non ha trovato errori e ha ritornato -1 allora il nome utente non esiste e va creato un nuovo utente
	if resSearch == -1 && errore == nil {
		_, err := db.c.Exec("INSERT INTO users (username, name, surname) VALUES ($1, $2, $3)", username, name, surname)
		if err != nil {
			log.Fatal("Failed to read record: ", err)
			return "", err
		}
		return username, nil

	} else if resSearch == -1 && errore != nil { // se la ricerca ha trovato errori e quindi non ha trovato il record allora ritorno errore
		log.Fatal("Failed to read record: ", errore)
		return "", errore

	} else { // nel caso in cui il record già esiste ritorno utente già esistente
		return "", errors.New("utente già registrato")
	}
}

// ritorna -1 se non c'è gia, altrimenti ritorna l'id dell'utente
func (db *appdbimpl) SearchUser(username string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found
			return -1, nil
		}
		// Other errors
		log.Fatal("Failed to read record: ", err)
		return -1, err
	}

	return id, nil
}
