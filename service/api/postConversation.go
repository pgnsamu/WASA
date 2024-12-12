package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// get param
	paramId := ps.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // Limit to 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve and parse form data
	name := r.FormValue("name")                   // String, no parsing needed
	isGroupStr := r.FormValue("isGroup")          // Boolean value sent as string (e.g., "true")
	description := r.FormValue("description")     // String, no parsing needed
	partecipantsIdStr := r.Form["partecipantsId"] // Slice of strings (e.g., ["1", "2", "3"])

	// Parse `isGroup` into a boolean
	isGroup, err := strconv.ParseBool(isGroupStr)
	if err != nil {
		http.Error(w, "Invalid value for isGroup", http.StatusBadRequest)
		return
	}

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
	if !isGroup && len(partecipantsId) != 1 {
		http.Error(w, "Invalid number of participants", http.StatusBadRequest)
		return
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

	// debug
	/*
		fmt.Println(id)
		fmt.Println(name)
		fmt.Println(description)
		fmt.Println(isGroup)
		fmt.Println(imgData)
	*/
	result, err := rt.db.CreateConversation(id, name, isGroup, &imgData, &description, partecipantsId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}
