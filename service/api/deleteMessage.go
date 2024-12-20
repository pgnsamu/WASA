package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stringIdConv := ps.ByName("conversationId")
	stringIdUser := ps.ByName("id")
	stringIdUserToDelete := ps.ByName("messageId")

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

	err = rt.db.DeleteMessage(idConv, idUser, idToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
