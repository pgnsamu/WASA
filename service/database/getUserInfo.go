package database

// errori ritornabili da GetUserInfo
// utente non trovato

func (db *appdbimpl) GetUserInfo(id int) (*User, error) {
	query := "SELECT id, username, photo FROM users WHERE id= ?"
	rows, err := db.c.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// La parola chiave defer viene utilizzata per pianificare una chiamata di funzione da eseguire dopo che la funzione circostante è stata completata.

	var user User
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrUserNotFound
	}
	// Controlla se si sono verificati errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
