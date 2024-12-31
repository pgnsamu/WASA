package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		//http.Error(w, "Token mancante o invalido", http.StatusUnauthorized)
		//return
		fmt.Println("Token mancante o invalido")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validazione del token
	_, err := ValidateJWT(tokenString)
	if err != nil {
		//http.Error(w, "Token non valido", http.StatusUnauthorized)
		//return
		fmt.Println("Token non valido")
	}

	stringIdConv := ps.ByName("conversationId")
	stringIdUser := ps.ByName("id")
	stringIdMessage := ps.ByName("messageId")
	stringIdCommentToDelete := ps.ByName("commentId")

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
	/*
		if idUser != claims["id"] {
			http.Error(w, "Utente non autorizzato", http.StatusUnauthorized)
			return
		}
	*/
	idMessage, err := strconv.Atoi(stringIdMessage)
	if err != nil {
		http.Error(w, "Errore id conversazione non intero", http.StatusBadRequest)
		return
	}

	idComment, err := strconv.Atoi(stringIdCommentToDelete)
	if err != nil {
		http.Error(w, "Errore id conversazione non intero", http.StatusBadRequest)
		return
	}

	isComment, err := rt.db.IsCommentTo(idComment, idMessage, idConv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isComment {
		http.Error(w, "Il messaggio non è un commento", http.StatusBadRequest)
		return
	}
	// TODO: capire se aggiungere il controllo di è un commento
	err = rt.db.DeleteMessage(idConv, idUser, idComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
