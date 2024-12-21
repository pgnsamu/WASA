package database

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) SaveImageToDB(imgData []byte, table string, field string, userId int) error {

	// Prepare the INSERT statement
	// stmt, err := db.c.Prepare("INSERT INTO ? (?) VALUES (?) ")
	// stmt, err := db.c.Prepare("UPDATE ? SET ? = ? WHERE id = ?;")
	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?;", table, field)
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute the statement, passing the image data as a BLOB
	result, err := stmt.Exec(imgData, userId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("id not found")
	} else {
		return nil
		//fmt.Println("Rows were updated successfully.")
	}
}
