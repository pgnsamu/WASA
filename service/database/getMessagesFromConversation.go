package database

type Message struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	PhotoContent []byte `json:"photoContent,omitempty"` // Optional in JSON
	SentAt       int    `json:"sentAt"`
	// ConversationID int 	`json:"conversationId"`
	AnswerTo       *int              `json:"answerTo,omitempty"` // Omit if nil
	IsForwarded    bool              `json:"isForwarded"`
	SenderID       int               `json:"senderId"`
	SenderUsername string            `json:"senderUsername"`
	Status         int               `json:"status"`
	Reactions      *[]SimpleReaction `json:"reactions,omitempty"`
}
type SimpleReaction struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	SentBy  string `json:"sentBy"`
	SentAt  int    `json:"sentAt"`
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
		// Get reactions for the current message
		reactionQuery := `
			SELECT r.id, r.reaction, u.username, r.sentAt
			FROM reactions as r, users u
			WHERE r.messageId = ? AND r.userId = u.id
		`
		reactionRows, err := db.c.Query(reactionQuery, msg.ID)
		if err != nil {
			return nil, err
		}
		defer reactionRows.Close()

		var reactions []SimpleReaction
		for reactionRows.Next() {
			var reaction SimpleReaction
			err := reactionRows.Scan(&reaction.ID, &reaction.Content, &reaction.SentBy, &reaction.SentAt)
			if err != nil {
				return nil, err
			}
			reactions = append(reactions, reaction)
		}

		if err := reactionRows.Err(); err != nil {
			return nil, err
		}

		msg.Reactions = &reactions
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &messages, nil
}
