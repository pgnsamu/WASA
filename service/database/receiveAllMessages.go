package database

// nessun errore ritornabile

func (db *appdbimpl) ReceiveAllMessages(idUser int) (int, error) {
	query := `
		UPDATE received
		SET status = 1
		WHERE userId = ? AND status = 0;
	`

	stmt, err := db.c.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(idUser)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
