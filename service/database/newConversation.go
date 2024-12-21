package database

import (
	"fmt"
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
}

func (db *appdbimpl) NewConversation(userId int, name string, isGroup bool, photo *[]byte, description *string, partecipantsId []int) (*Conversation, error) {

	// TODO: aggiungere sulla tabella participate collegamento utente conversation per ogni partecipantId+userId
	// TODO: mancante nel db il check value del fatto che se isGroup allora partecipantsId == 1

	if isGroup {
		// se è un gruppo riprepariamo
		stmt, err := db.c.Prepare("INSERT INTO conversations (name, createdAt, isGroup, description, photo) VALUES (?, ?, ?, ?, ?);")
		if err != nil {
			return nil, fmt.Errorf("prepare statement: %w", err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(name, time.Now().UnixMilli(), isGroup, *description, *photo)
		if err != nil {
			return nil, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("getting last insert ID: %w", err)
		}

		tempParticipantsId := partecipantsId
		tempParticipantsId = append(tempParticipantsId, userId) // Append userId to the slice before the loop
		// fmt.Println(len(tempParticipantsId)) debug:

		for i := 0; i < len(tempParticipantsId); i++ { // TODO: è stato aggiunto un meno 1 ma non so se è giusto controllare (ora è stato tolto e funziona, non lo so )
			stmt, err := db.c.Prepare("INSERT INTO participate (userId, conversationId) VALUES (?, ?);")
			if err != nil {
				return nil, fmt.Errorf("prepare statement: %w", err)
			}
			defer stmt.Close()
			// fmt.Println(tempParticipantsId[i])
			_, err = stmt.Exec(tempParticipantsId[i], lastInsertID)
			if err != nil {
				// fmt.Println(tempParticipantsId[i])
				// fmt.Println(err.Error()) // errore qui UNIQUE constraint failed: participate.userId, participate.conversationId
				return nil, err
			}
		}

		conversation, err := db.GetConversationInfo(int(lastInsertID), userId)
		if err != nil {
			return nil, fmt.Errorf("getting last insert ID: %w", err)
		}
		return conversation, nil
	} else {
		stmt, err := db.c.Prepare("INSERT INTO conversations (name, createdAt, isGroup) VALUES (?, ?, ?);")
		if err != nil {
			return nil, fmt.Errorf("prepare statement: %w", err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(name, time.Now().UnixMilli(), isGroup)
		if err != nil {
			return nil, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("getting last insert ID: %w", err)
		}

		tempParticipantsId := partecipantsId
		tempParticipantsId = append(tempParticipantsId, userId)
		for i := 0; i < len(tempParticipantsId); i++ {
			stmt, err := db.c.Prepare("INSERT INTO participate (userId, conversationId) VALUES (?, ?);")
			if err != nil {
				return nil, fmt.Errorf("prepare statement: %w", err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(tempParticipantsId[i], lastInsertID)
			if err != nil {
				return nil, err
			}
		}

		conversation, err := db.GetConversationInfo(int(lastInsertID), userId)
		if err != nil {
			return nil, fmt.Errorf("getting last insert ID: %w", err)
		}
		return conversation, nil
	}

}
