package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type groupData struct {
	Username string `json:"groupName"`
}

func (rt *_router) postGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data groupData
	paramId := ps.ByName("id")
	paramId2 := ps.ByName("conversationId")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	idConversation, err := strconv.Atoi(paramId2)
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

	user, err := rt.db.SetupGroupName(id, idConversation, data.Username)

	if user == nil && err != nil && err.Error() == "userNotFound" {
		http.Error(w, "Errore id non registrato", http.StatusBadRequest)
		return
	} else if user == nil && err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
