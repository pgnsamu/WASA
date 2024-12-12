package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type targetB struct {
	Id int `json:"targetConversationId"`
}

// TODO: errori non gestiti perfetti
func (rt *_router) postForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	paramId := ps.ByName("id")
	idUser, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	paramId2 := ps.ByName("conversationId")
	idConv, err := strconv.Atoi(paramId2)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	paramId3 := ps.ByName("messageId")
	idMessage, err := strconv.Atoi(paramId3)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Create a variable to hold the unmarshalled data
	var requestData targetB

	// Unmarshal the JSON body into the struct
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(idConv, requestData.Id, idUser, idMessage)
	_, err = rt.db.ForwardMessage(idConv, requestData.Id, idUser, idMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully")
}
