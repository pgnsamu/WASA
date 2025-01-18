package database

import "errors"

// errori ritornabili da RemoveReaction
// utente non presente nella chat
// il messaggio non appartiene a questa conversazione
// reazione non trovata

func (db *appdbimpl) RemoveReaction(idConversation int, idUser int, idMessage int, idReaction int) (*[]Reaction, error) {
	resu, err := db.UserExist(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	if !resu {
		return nil, errors.New("utente non presente nella chat")
	}
	messageBelongs, err := db.MessageBelongsToConversation(idMessage, idConversation)
	if err != nil {
		return nil, err
	}
	if !messageBelongs {
		return nil, errors.New("il messaggio non appartiene a questa conversazione")
	}

	// Controlla se la reazione esistes
	queryStr := "SELECT EXISTS(SELECT 1 FROM reactions WHERE id = ? AND userId = ?);"
	var exists int

	err = db.c.QueryRow(queryStr, idReaction, idUser).Scan(&exists)
	if err != nil {
		return nil, err
	}

	// Controlla il risultato
	if exists != 1 {
		return nil, errors.New("reazione non trovata")
	}

	queryStr = "DELETE FROM reactions WHERE id = ?"
	stmt, err := db.c.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(idReaction)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
