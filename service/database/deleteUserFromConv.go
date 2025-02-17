package database

import (
	"errors"
)

// errori ritornabili da DeleteUserFromConv
// partecipanti non trovati
// utente da eliminare non trovato
// ID non trovato
// troppi utenti trovati

func (db *appdbimpl) DeleteUserFromConv(idConversation int, idUser int, idUserToDelete int) error {
	// prendo tutti i partecipanti di una certa conversazione in cui sia presente l'utente che sta facendo la query
	users, err := db.GetUsersOfConversation(idConversation, idUser)
	if err != nil {
		return err
	}

	found := false
	valueOfUser := *users // unwrapping pointer
	for i := 0; i < len(valueOfUser); i++ {
		if idUserToDelete == valueOfUser[i].Id { // ricerco se esiste tra questi, quello da eliminare
			found = true
			break
		}
	}

	if !found {
		return errors.New("utente da eliminare non trovato")
	}

	// faccio select per controllare che una sola riga sia presa di mira
	stringQuery := "SELECT userId FROM participate as p WHERE p.userId = ? and p.conversationId = ?;"
	rows, err := db.c.Query(stringQuery, idUser, idConversation)
	if err != nil {
		return err
	}
	defer rows.Close()

	var idResults []string
	for rows.Next() {
		var tempString string
		err := rows.Scan(&tempString)
		if err != nil {
			return err
		}

		idResults = append(idResults, tempString)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if len(idResults) == 0 {
		return errors.New("ID non trovato")
	} else if len(idResults) > 1 {
		return errors.New("troppi utenti trovati")
	} else {
		stringQuery = "DELETE FROM participate as p WHERE p.userId = ? and p.conversationId = ?;"
		stmt, err := db.c.Prepare(stringQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(idUserToDelete, idConversation)
		if err != nil {
			return err
		}

		return nil

	}

}
