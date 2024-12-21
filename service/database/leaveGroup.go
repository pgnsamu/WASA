package database

import (
	"errors"
)

func (db *appdbimpl) DeleteUserFromConv(idConversation int, idUser int, idUserToDelete int) (*[]User, error) {
	// prendo tutti i partecipanti di una certa conversazione in cui sia presente l'utente che sta facendo la query
	users, err := db.GetUsersOfConversation(idConversation, idUser)
	if err != nil {
		return nil, err
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
		return nil, errors.New("utente da eliminare non trovato")
	}

	// faccio select per controllare che una sola riga sia presa di mira
	stringQuery := "SELECT userId FROM participate as p WHERE p.userId = ? and p.conversationId = ?;"
	rows, err := db.c.Query(stringQuery, idUser, idConversation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idResults []string
	// Iterate over rows
	for rows.Next() {
		var tempString string
		err := rows.Scan(&tempString)

		if err != nil {
			return nil, err
		}
		idResults = append(idResults, tempString)
	}
	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(idResults) == 0 {
		return nil, errors.New("id not found")
	} else if len(idResults) > 1 {
		return nil, errors.New("too much users found")
	} else {
		stringQuery = "DELETE FROM participate as p WHERE p.userId = ? and p.conversationId = ?;"
		stmt, err := db.c.Prepare(stringQuery)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(idUserToDelete, idConversation)
		if err != nil {
			return nil, err
		}

		if idUserToDelete != idUser {
			users, err = db.GetUsersOfConversation(idConversation, idUser)
			if err != nil {
				return nil, err
			}
			// TODO: controllare se esista almeno un altro partecipante dopo l'eliminazione in caso sia una isGroup = false
			// 		 controllare se esista almeno un altro partecipante dopo l'eliminazione in caso sia una isGroup = true
			return users, nil
		} else {
			return nil, nil
		}

	}

}
