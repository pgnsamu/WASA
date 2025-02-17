package database

import (
	"errors"
)

// errori ritornabili da addtogroup
// utente non trovato
// partecipanti non trovati
// chat piena
// utente da aggiungere già presente

func (db *appdbimpl) AddToGroup(idConversation int, idUser int, idUserToAdd int) (*[]User, error) {

	// controllare che l'utente esista
	// potrebbe ritornare error "user not found" -> utente non trovato
	_, err := db.GetUserInfo(idUserToAdd)
	if err != nil {
		return nil, err
	}

	// prendo tutti gli utenti della conversazione
	// potrebbe ritornare error "partecipanti non trovati"
	users, err := db.GetUsersOfConversation(idConversation, idUser)
	if err != nil {
		return nil, err
	}

	found := false
	valueOfUser := *users // unwrapping puntatore
	for i := 0; i < len(valueOfUser); i++ {
		if idUserToAdd == valueOfUser[i].Id {
			// vedo se tra questi c'è già quello da aggiungere
			found = true
			break
		}
	}

	// controllare se la conversazione è satura
	saturationQuery := `
		SELECT count(p.userId), c.id, isGroup
		FROM participate as p, conversations as c 
		WHERE p.conversationId = c.id and c.id = ?
		GROUP BY c.id
	`

	var userCount int
	var conversationID int
	var isGroup bool

	err = db.c.QueryRow(saturationQuery, idConversation).Scan(&userCount, &conversationID, &isGroup)
	if err != nil {
		return nil, err
	}

	if isGroup {
		if (userCount + 1) > 64 { // quindi l'aggiunta è di troppo
			return nil, errors.New("chat piena")
		}
	} else {
		if (userCount + 1) != 2 { // quindi l'aggiunta è di troppo
			return nil, errors.New("chat piena")
		}
	}

	if found {
		return nil, errors.New("utente da aggiungere già presente")
	}

	// faccio select per controllare che una sola riga sia presa di mira
	stringQuery := "INSERT INTO participate (userId, conversationId) VALUES (?, ?);"
	stmt, err := db.c.Prepare(stringQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(idUserToAdd, idConversation)
	if err != nil {
		return nil, err
	}

	u, err := db.GetUsersOfConversation(idConversation, idUser)
	if err != nil {
		return nil, err
	}

	return u, nil

}
