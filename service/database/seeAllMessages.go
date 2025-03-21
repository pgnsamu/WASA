package database

// nessun errore ritornabile

func (db *appdbimpl) SeeAllMessages(idUser int, idConv int) (int, error) {
	query := `
		UPDATE received
		SET status = 2
		WHERE userId = ? AND status = 1 AND messageId IN (
			SELECT m.id
			FROM messages as m
			WHERE m.conversationId = ?
		)
	`

	stmt, err := db.c.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(idUser, idConv)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
