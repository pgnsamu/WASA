package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// errori che possono essere ritornati da deleteMessage
// utente non nel gruppo
// autore del messaggio sbagliato
// nessun messaggio trovato con quell'id
// ritorna nulla

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Token mancante o invalido", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validazione del token
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

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
	if idUser != claims["id"] {
		http.Error(w, "Utente non autorizzato", http.StatusUnauthorized)
		return
	}
	idToDelete, err := strconv.Atoi(stringIdUserToDelete)
	if err != nil {
		http.Error(w, "Errore id conversazione non intero", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(idConv, idUser, idToDelete)
	if err != nil && (err.Error() == "utente non nel gruppo" ||
		err.Error() == "autore del messaggio sbagliato" ||
		err.Error() == "nessun messaggio trovato con quell'id") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
