package database

import (
	"errors"
	"log"
)

/*
	func (db *appdbimpl) GetUserInfo(id int) (User, error) {
		var user User
		query := "SELECT id, username, name, surname FROM users WHERE id=$1"
		err := db.c.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Name, &user.Surname)
		if err != nil {
			if err == sql.ErrNoRows {
				// No rows foundm
				return User{}, errors.New("userNotFound")
			}
			// Other errors
			log.Fatal("Failed to read record: ", err)
			return User{}, err
		}
		return user, nil
	}
*/

func (db *appdbimpl) GetUserInfo(id int) (*User, error) {
	query := "SELECT id, username, name, surname, photo FROM users WHERE id=$1"
	rows, err := db.c.Query(query, id)
	if err != nil {
		log.Fatal("Error executing query:", err)
		return nil, err
	}
	defer rows.Close() // defer keyword is used to schedule a function call to be executed after the surrounding function has completed. Deferred functions are typically used for cleanup tasks

	var user User
	// scanning multiplo anche per ricercare una singola riga in modo da passare tutti i parametri
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Photo)
		if err != nil {
			log.Fatal("Error scanning row:", err)
			return nil, err
		}
	} else {
		errorUserNotFound := errors.New("user not found")
		return nil, errorUserNotFound
	}
	return &user, nil
}
