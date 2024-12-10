package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getParticipants(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	w.Header().Set("content-type", "application/json")
	users, err := rt.db.GetUsersOfConversation(idConv, idUser)
	if err != nil && err.Error() == "partecipanti non trovati" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(*users)

}
