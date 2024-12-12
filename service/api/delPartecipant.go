package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// TODO: sistemare il Doc delle api perché è stato implementato in modo che utente x possa eliminare qualunque utente y nello stesso gruppo
func (rt *_router) delParticipant(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stringIdConv := ps.ByName("conversationId")
	stringIdUser := ps.ByName("id")
	stringIdUserToDelete := ps.ByName("toDelete")

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
	idToDelete, err := strconv.Atoi(stringIdUserToDelete)
	if err != nil {
		http.Error(w, "Errore id conversazione non intero", http.StatusBadRequest)
		return
	}

	users, err := rt.db.DeleteUserFromConv(idConv, idUser, idToDelete)
	if err != nil && err.Error() == "partecipanti non trovati" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if users == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(*users)

}
