package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Define a struct to match the incoming JSON
type requestBodyId struct {
	UserId int `json:"userId"`
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stringIdConv := ps.ByName("conversationId")
	stringIdUser := ps.ByName("id")

	// conversione string to int
	idConv, err := strconv.Atoi(stringIdConv)
	if err != nil {
		http.Error(w, "Errore id conversazione non intero", http.StatusBadRequest)
		return
	}
	idUser, err := strconv.Atoi(stringIdUser)
	if err != nil {
		http.Error(w, "Errore id utente non intero", http.StatusBadRequest)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Create a variable to hold the unmarshalled data
	var requestData requestBodyId

	// Unmarshal the JSON body into the struct
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	users, err := rt.db.AddToGroup(idConv, idUser, requestData.UserId)
	if err != nil && (err.Error() == "chat piena" || err.Error() == "utente da aggiungere già presente") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(*users)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}

}
