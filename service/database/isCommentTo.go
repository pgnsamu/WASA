package database

func (db *appdbimpl) IsCommentTo(idComment int, idMessage int, idConversation int) (bool, error) {

	queryStr := "SELECT senderId FROM messages as m WHERE m.id = ? and m.conversationId = ? and m.answerTo = ?"
	stmt, err := db.c.Prepare(queryStr)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Eseguire la query con un parametro
	rows, err := stmt.Query(idComment, idConversation, idMessage)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil

}
