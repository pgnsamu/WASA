package api

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// errori ritornabili da setGroupPhoto
// id non trovato
// ritorna nulla

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// Get the file from the form data
	file, _, err := r.FormFile("photo") // "file" is the form field name
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file data", http.StatusInternalServerError)
		return
	}
	resu, err := rt.db.UserExist(idConv, idUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !resu {
		http.Error(w, "user not in the group", http.StatusBadRequest)
	}

	err = rt.db.SaveImageToDB(imgData, "conversations", "photo", idConv)
	if err != nil {
		if err.Error() == "id non trovato" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
