package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type usernameData struct {
	Username string `json:"username"`
}

func (rt *_router) putUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data usernameData
	paramId := ps.ByName("id")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}

	// Decode the JSON body
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil || data.Username == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := rt.db.UpdateUsername(id, data.Username)

	if user == nil && err != nil && err.Error() == "userNotFound" {
		http.Error(w, "Errore id non registrato", http.StatusBadRequest)
		return
	} else if user == nil && err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Println(user.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
