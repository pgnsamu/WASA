package database

import (
	"errors"
)

// errori ritornabili da ForwardMessage
// utente non presente nella chat sorgente
// utente non presente nella chat di destinazione
// messaggio non trovato

func (db *appdbimpl) ForwardMessage(idConversationSource int, idConversationDest int, idUser int, idMessage int) (*Conversation, error) {

	// TODO: capire se da errore anche se non esiste la conversazione nel db
	// controllo se l'utente che vuole mandare sta nella chat sorgente
	resu, err := db.UserExist(idConversationSource, idUser)
	if err != nil {
		return nil, err
	}
	if !resu {
		return nil, errors.New("utente non presente nella chat sorgente")
	}

	// controllo se l'utente esiste nella conversazione di dest
	res, err := db.UserExist(idConversationDest, idUser)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, errors.New("utente non presente nella chat di destinazione")
	}

	queryStr := "SELECT EXISTS(SELECT 1 FROM messages WHERE id = ? AND conversationId = ?);"
	var esiste int

	err = db.c.QueryRow(queryStr, idMessage, idConversationSource).Scan(&esiste)
	if err != nil {
		return nil, err
	}

	// Controlla il risultato
	if esiste != 1 {
		return nil, errors.New("messaggio non trovato")
	}

	queryStr = `
		SELECT content, photoContent
		FROM messages
		WHERE id = ?;
	`
	stmt, err := db.c.Prepare(queryStr)
	if err != nil {
		return nil, err
	}

	var content string
	var photoContent []byte

	err = stmt.QueryRow(idMessage).Scan(&content, &photoContent)
	if err != nil {
		return nil, err
	}

	_, err = db.SendMessage(idConversationDest, idUser, content, photoContent, nil, 1)
	if err != nil {
		return nil, err
	}

	result, err := db.GetConversationInfo(idConversationDest, idUser)
	if err != nil {
		return nil, err
	}
	return result, nil
}
