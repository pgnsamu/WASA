package database

// errrori ritornabili da GetUserId
// utente non trovato

func (db *appdbimpl) GetUserId(username string) (*int, error) {
	query := "SELECT id FROM users WHERE username = ?"
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idUser int
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&idUser)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrUserNotFound
	}
	// Controlla se ci sono errori che possono essersi verificati durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &idUser, nil
}
