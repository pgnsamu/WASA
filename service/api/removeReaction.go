package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// errori ritornabili da removeReaction
// utente non presente nella chat
// il messaggio non appartiene a questa conversazione
// reazione non trovata
// ritorna nulla

func (rt *_router) removeReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	paramId := ps.ByName("id")
	idUser, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	if idUser != claims["id"] {
		http.Error(w, "Utente non autorizzato", http.StatusUnauthorized)
		return
	}
	paramId2 := ps.ByName("conversationId")
	idConv, err := strconv.Atoi(paramId2)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	paramId3 := ps.ByName("messageId")
	messageId, err := strconv.Atoi(paramId3)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	paramId4 := ps.ByName("reactionId")
	reactionId, err := strconv.Atoi(paramId4)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}

	_, err = rt.db.RemoveReaction(idConv, idUser, messageId, reactionId)
	if err != nil && (err.Error() == "utente non presente nella chat" || err.Error() == "il messaggio non appartiene a questa conversazione" || err.Error() == "reazione non trovata") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
