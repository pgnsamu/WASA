package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type requestBody struct {
	Username string `json:"username"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	// Name     string `json:"name"`
	// Surname  string `json:"surname"`
}

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
		var username string
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&username)
		fmt.Printf("Decoded name: %s\n", username)
	*/

	// Leggo il body della request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusInternalServerError)
		return
	}
	var reqBody requestBody
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := rt.db.CreateUser(reqBody.Username, "", "")
	if err.Error() == "Utente già registrato" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

	/*
		// Respond back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Body received successfully"}`))

		w.Header().Set("content-type", "text/plain")
		name, err := rt.db.GetUsersDB()
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(name)
	*/

}
