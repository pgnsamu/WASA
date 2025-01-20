package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Conversation struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	IsGroup     bool    `json:"isGroup"`
	CreatedAt   int     `json:"createdAt"`
	Description *string `json:"description,omitempty"`
	Photo       *[]byte `json:"photo,omitempty"`
	LastMessage *string `json:"lastMessage,omitempty"`
}

// errori ritornabili da newConversation
// utente non registrato
// utente non trovato
// chat già esistente
// TODO: forse si puo mettere tutto in uno e uno in tutti

func (db *appdbimpl) NewConversation(userId int, name string, isGroup bool, photo *[]byte, description *string, partecipantsId []int) (*Conversation, error) {

	// TODO: aggiungere sulla tabella participate collegamento utente conversation per ogni partecipantId+userId
	// TODO: mancante nel db il check value del fatto che se isGroup allora partecipantsId == 1

	if isGroup {
		// se è un gruppo riprepariamo
		stmt, err := db.c.Prepare("INSERT INTO conversations (name, createdAt, isGroup, description, photo) VALUES (?, ?, ?, ?, ?);")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		result, err := stmt.Exec(name, time.Now().UnixMilli(), isGroup, *description, *photo)
		if err != nil {
			return nil, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		tempParticipantsId := partecipantsId
		tempParticipantsId = append(tempParticipantsId, userId) // Append userId to the slice before the loop

		// aggiunta dei partecipanti a un gruppo nella tabella
		for i := 0; i < len(tempParticipantsId); i++ { // TODO: è stato aggiunto un meno 1 ma non so se è giusto controllare (ora è stato tolto e funziona, non lo so )

			stmt, err := db.c.Prepare("INSERT INTO participate (userId, conversationId) VALUES (?, ?);")
			if err != nil {
				return nil, err
			}
			defer stmt.Close()

			_, err = stmt.Exec(tempParticipantsId[i], lastInsertID)
			if err != nil {
				return nil, err
			}
		}

		conversation, err := db.GetConversationInfo(int(lastInsertID), userId)
		if err != nil {
			return nil, err
		}
		return conversation, nil
	} else {
		// controllo se l'utente ha già una chat privata con l'utente
		var existingConversationId int
		q := `SELECT c.id 
			FROM conversations c 
			JOIN participate p1 ON c.id = p1.conversationId 
			JOIN participate p2 ON c.id = p2.conversationId 
			WHERE p1.userId = ? AND p2.userId = ? AND c.isGroup = false
		`
		err := db.c.QueryRow(q, userId, partecipantsId[0]).Scan(&existingConversationId)
		if err == nil {
			return nil, errors.New("chat già esistente")
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		// se è una chat privata riprepariamo
		stmt, err := db.c.Prepare("INSERT INTO conversations (name, createdAt, isGroup) VALUES (?, ?, ?);")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		result, err := stmt.Exec(name, time.Now().UnixMilli(), isGroup)
		if err != nil {
			return nil, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		tempParticipantsId := partecipantsId
		tempParticipantsId = append(tempParticipantsId, userId)
		// TODO: vedere se è da controllare il numero <3
		for i := 0; i < len(tempParticipantsId); i++ {
			stmt, err := db.c.Prepare("INSERT INTO participate (userId, conversationId) VALUES (?, ?);")
			if err != nil {
				return nil, err
			}
			defer stmt.Close()

			_, err = stmt.Exec(tempParticipantsId[i], lastInsertID)
			if err != nil {
				return nil, err
			}
		}

		conversation, err := db.GetConversationInfo(int(lastInsertID), userId)
		if err != nil {
			return nil, err
		}
		return conversation, nil
	}

}
