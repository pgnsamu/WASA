package database

import (
	"errors"
	"time"
)

// errori ritornabili da MessageBelongsToConversation
// nessuno

func (db *appdbimpl) MessageBelongsToConversation(idMessage int, idConversation int) (bool, error) {
	queryStr := "SELECT EXISTS(SELECT 1 FROM messages WHERE id = ? AND conversationId = ?);"
	var exists int

	err := db.c.QueryRow(queryStr, idMessage, idConversation).Scan(&exists)
	if err != nil {
		return false, err
	}

	// Controlla il risultato
	if exists != 1 {
		return false, nil
	}

	return true, nil
}

type Reaction struct {
	MessageID int    `json:"messageId"`
	UserID    int    `json:"userId"`
	Content   string `json:"content"`
	SentAt    int64  `json:"sentAt,omitempty"`
}

// errori ritornabili da SendReaction
// utente non presente nella chat
// il messaggio non appartiene a questa conversazione

func (db *appdbimpl) SendReaction(idConversation int, idUser int, content string, replyTo int) (*[]Reaction, error) {
	resu, err := db.UserExist(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	if !resu {
		return nil, errors.New("utente non presente nella chat")
	}

	// Controlla se il messaggio appartiene alla conversazione
	messageBelongs, err := db.MessageBelongsToConversation(replyTo, idConversation)
	if err != nil {
		return nil, err
	}
	if !messageBelongs {
		return nil, errors.New("il messaggio non appartiene a questa conversazione")
	}

	queryStr := "INSERT INTO reactions (messageId, userId, reaction, sentAt) VALUES (?, ?, ?, ?)"
	stmt, err := db.c.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ottieni il timestamp corrente
	sentAt := time.Now().UnixMilli()

	// Inserisci la reazione nel database
	resul, err := stmt.Exec(replyTo, idUser, content, sentAt)
	if err != nil {
		return nil, err
	}

	// Recupera l'ID della reazione appena inserita
	_, err = resul.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Crea l'oggetto Reaction da restituire
	reaction := Reaction{
		MessageID: replyTo,
		UserID:    idUser,
		Content:   content,
		SentAt:    sentAt,
	}

	// Restituisci la reazione inserita
	return &[]Reaction{reaction}, nil

}
