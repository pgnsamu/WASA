package database

// errori ritornabili da GetProfilePhoto
// utente non trovato

func (db *appdbimpl) GetProfilePhoto(id int) ([]byte, error) {
	query := "SELECT photo FROM users WHERE id = ?"
	rows, err := db.c.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photo []byte
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&photo)
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
	return photo, nil
}
