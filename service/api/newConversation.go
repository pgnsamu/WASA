package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) newConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// get param
	paramId := ps.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	if id != claims["id"] {
		http.Error(w, "Utente non autorizzato", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // Limit to 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve and parse form data
	name := r.FormValue("name")                       // String, no parsing needed
	isGroupStr := r.FormValue("isGroup")              // Boolean value sent as string (e.g., "true")
	description := r.FormValue("description")         // String, no parsing needed
	partecipantsStr := r.Form["partecipantsUsername"] // Slice of strings (e.g., ["1", "2", "3"])

	// Parse `isGroup` into a boolean
	isGroup, err := strconv.ParseBool(isGroupStr)
	if err != nil {
		http.Error(w, "Invalid value for isGroup", http.StatusBadRequest)
		return
	}
	/*
		// Parse `partecipantsId` into a slice of integers
		partecipantsId := []int{}
		for _, idStr := range partecipantsIdStr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid participant ID", http.StatusBadRequest)
				return
			}
			partecipantsId = append(partecipantsId, id)
		}
	*/
	if !isGroup && len(partecipantsStr) != 1 {
		http.Error(w, "Invalid number of participants", http.StatusBadRequest)
		return
	}

	var partecipantsId = make([]int, 0, len(partecipantsStr))
	for _, username := range partecipantsStr {
		tempId, err := rt.db.GetUserId(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if *tempId == id {
			http.Error(w, "Invalid participant ID", http.StatusBadRequest)
			return
		}
		partecipantsId = append(partecipantsId, *tempId)
	}

	var imgData []byte

	file, _, err := r.FormFile("photo") // "file" is the form field name
	if err != nil && isGroup {
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

	result, err := rt.db.NewConversation(id, name, isGroup, &imgData, &description, partecipantsId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}

}
