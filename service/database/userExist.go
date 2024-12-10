package database

func (db *appdbimpl) UserExist(idConv int, idUser int) (bool, error) {

	// controllo utente se è all'interno
	stmt, err := db.c.Prepare("SELECT userId FROM participate as p WHERE p.conversationId = ? and p.userId = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(idConv, idUser)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Count the number of rows
	count := 0
	for rows.Next() {
		count++
	}

	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return false, err
	}

	if count != 1 {
		return false, nil
	} else {
		return true, nil
	}
}
