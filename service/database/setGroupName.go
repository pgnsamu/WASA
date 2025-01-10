package database

import (
	"errors"
)

// ritornare la conversation aggiornata
func (db *appdbimpl) SetGroupName(idUser int, idConversation int, name string) (*Conversation, error) {

	// controllo utente se è all'interno
	stmt, err := db.c.Prepare("SELECT userId FROM participate as p WHERE p.conversationId = ? and p.userId = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Count the number of rows
	count := 0
	for rows.Next() {
		count++
	}

	// Check for errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if count != 1 {
		return nil, errors.New("utente non registrato nel gruppo")
	}

	// preparazione della query
	stmt, err = db.c.Prepare("UPDATE conversations SET name = ? WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// esecuzione
	res, err := stmt.Exec(name, idConversation)
	if err != nil {
		return nil, err
	}

	// controllo righe interessate che in questo caso dovranno essere == 1
	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 1 {
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New("too much rows")
		}
	}
	var conv *Conversation
	conv, err = db.GetConversationInfo(idConversation, idUser)
	if err != nil {
		return nil, err
	}

	return conv, nil
}
