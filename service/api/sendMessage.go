package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// errori ritornabili da sendMessage
// il messaggio non appartiene a questa conversazione
// ritorna []Message lista messaggi della convo

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// Parse the form to handle file uploads
	err = r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
		return
	}
	// content, photoContent, senderId, sentAt, conversationId, answerTo

	// Retrieve values from the form
	content := r.FormValue("content")
	messageTypeStr := r.FormValue("isPhoto")

	// messageType = false not image
	// messageType = true image
	// Parse `isGroup` into a boolean
	messageType, err := strconv.ParseBool(messageTypeStr)
	if err != nil {
		http.Error(w, "Invalid value for messageType", http.StatusBadRequest)
		return
	}

	var imgData []byte
	file, _, err := r.FormFile("photo") // "file" is the form field name
	if err != nil && messageType {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	} else if err == nil {
		defer file.Close()
		imgData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file data", http.StatusInternalServerError)
			return
		}
	}

	if (content == "" && !messageType) || (messageType && file == nil) {
		http.Error(w, "content is empty", http.StatusBadRequest)
		return
	}

	resu, err := rt.db.UserExist(idConv, idUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !resu {
		http.Error(w, "user not in the group", http.StatusBadRequest)
		return
	}

	res, err := rt.db.SendMessage(idConv, idUser, content, imgData, nil, 0)
	if err != nil && err.Error() == "il messaggio non appartiene a questa conversazione" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}
}
