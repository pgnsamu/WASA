package database

type Message struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	PhotoContent []byte `json:"photoContent,omitempty"` // Optional in JSON
	SentAt       int    `json:"sentAt"`
	// ConversationID int 	`json:"conversationId"`
	AnswerTo       *int   `json:"answerTo,omitempty"` // Omit if nil
	IsForwarded    bool   `json:"isForwarded"`
	SenderID       int    `json:"senderId"`
	SenderUsername string `json:"senderUsername"`
	Status         int    `json:"status"`
}

// TODO: da scrivere l'endpoint
func (db *appdbimpl) GetMessagesFromConversation(conversationID int) (*[]Message, error) {
	query := `
		SELECT m.id, m.content, m.photoContent, m.sentAt, m.answerTo, m.isForwarded, m.senderId, u.username, MIN(r.status) as status
		FROM messages as m, received as r, users u
		WHERE m.conversationId = ? AND m.id = r.messageId  AND u.id = m.senderId
		GROUP BY m.id
		ORDER BY m.sentAt ASC
	`

	rows, err := db.c.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message
		// var answerTo sql.NullInt64 // Handle nullable integer for the answerTo column

		err := rows.Scan(&msg.ID, &msg.Content, &msg.PhotoContent, &msg.SentAt, &msg.AnswerTo, &msg.IsForwarded, &msg.SenderID, &msg.SenderUsername, &msg.Status)
		if err != nil {
			return nil, err
		}
		/*
			// Convert sql.NullInt64 to *int
			if answerTo.Valid {
				val := int(answerTo.Int64)
				msg.AnswerTo = &val
			} else {
				msg.AnswerTo = nil
			}
		*/
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &messages, nil
}
