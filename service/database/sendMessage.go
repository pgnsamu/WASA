package database

import (
	"errors"
	"time"
)

// TODO: messagetype per ora sul doc è enum invece di bool
// TODO: togliere il previewContent dal doc

func (db *appdbimpl) SendMessage(idConversation int, idUser int, content string, photoContent []byte, messageType bool, replyTo *int) (*[]Message, error) {
	resu, err := db.UserExist(idConversation, idUser)
	if err != nil {
		return nil, err
	}
	if !resu {
		return nil, errors.New("utente non presente nella chat")
	}
	// fare insert
	query := `
	INSERT INTO messages (content, photoContent, senderId, sentAt, conversationId, answerTo)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	// Ottieni il timestamp corrente
	sentAt := time.Now().UnixMilli()

	// Esegui la query
	_, err = db.c.Exec(query, content, photoContent, idUser, sentAt, idConversation, replyTo)
	if err != nil {
		return nil, err
	}

	res, err := db.GetMessagesByConversation(idConversation)
	if err != nil {
		return nil, err
	}
	return res, nil

}
