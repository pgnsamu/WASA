package database

import "errors"

func (db *appdbimpl) GetUsersOfConversation(idConversation int, idUser int) (*[]User, error) {
	// controllo se l'utente è presente in participate in collegamento con idConv
	stringQuery := `SELECT u.id, u.username, u.photo
					FROM users AS u
					JOIN participate AS p
					ON p.userId = u.id
					WHERE p.conversationId = ? 
						AND p.conversationId IN (
							SELECT 
								p2.conversationId 
							FROM 
								participate AS p2
							WHERE 
								p2.userId = ?
						);`
	// Execute the query
	rows, err := db.c.Query(stringQuery, idConversation, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Prepare a slice to hold the results
	var users []User

	// Iterate over rows
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Photo)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, errors.New("partecipanti non trovati")
	}

	return &users, nil
}
