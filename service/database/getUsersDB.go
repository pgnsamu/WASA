package database

import (
	"fmt"
)

func (db *appdbimpl) GetUsersDB() (*[]User, error) {
	// Query to get all users
	query := "SELECT id, username, name, surname, photo FROM users"
	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Slice to hold all users
	var users []User

	// Iterate over rows and scan the data into the User struct
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Photo); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &users, nil
}
