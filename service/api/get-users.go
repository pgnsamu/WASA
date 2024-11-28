package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type User struct {
	id int 			`json: "id"`
	userName string `json: "userName"`
	nome int 		`json: "nome"`
	cognome string  `json: "cognome"`
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//prendi i dati dal db
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, "Error retrieving users from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	//creazione array dove inserirli
	var users []User

	for rows.next(){
		var user User
		err := rows.Scan(&user.id, &user.userName, &user.nome, &user.cognome)
		if err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	
	// Convert the users slice to JSON
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	json.NewEncoder(w).Encode(users)                  // Write JSON response
}
