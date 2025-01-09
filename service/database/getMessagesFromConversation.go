package database

type Message struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	PhotoContent []byte `json:"photoContent,omitempty"` // Optional in JSON
	SentAt       int    `json:"sentAt"`
	// ConversationID int 	`json:"conversationId"`
	AnswerTo    *int `json:"answerTo,omitempty"` // Omit if nil
	IsForwarded bool `json:"isForwarded"`
	SenderID    int  `json:"senderId"`
}

// TODO: da scrivere l'endpoint
func (db *appdbimpl) GetMessagesFromConversation(conversationID int) (*[]Message, error) {
	query := `
        SELECT id, content, photoContent, sentAt, answerTo, isForwarded, senderId
        FROM messages
        WHERE conversationId = ?
        ORDER BY sentAt ASC
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

		err := rows.Scan(&msg.ID, &msg.Content, &msg.PhotoContent, &msg.SentAt, &msg.AnswerTo, &msg.IsForwarded, &msg.SenderID)
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
