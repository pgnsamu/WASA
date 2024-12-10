package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessageReq(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// Parse the form to handle file uploads
	err = r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
		return
	}
	// content, photoContent, senderId, sentAt, conversationId, answerTo

	// Retrieve values from the form
	content := r.FormValue("content")
	messageTypeStr := r.FormValue("messageType")

	// messageType = false not image
	// messageType = true image
	// Parse `isGroup` into a boolean
	messageType, err := strconv.ParseBool(messageTypeStr)
	if err != nil {
		http.Error(w, "Invalid value for messageType", http.StatusBadRequest)
		return
	}
	// Convert answerTo to *int (nullable)
	answerToStr := r.FormValue("answerTo")
	var answerTo *int
	if answerToStr != "" {
		answerToVal, err := strconv.Atoi(answerToStr)
		if err != nil {
			http.Error(w, "Invalid answerTo value", http.StatusBadRequest)
			return
		}
		answerTo = &answerToVal
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

	resu, err := rt.db.UserExist(idConv, idUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !resu {
		http.Error(w, "user not in the group", http.StatusBadRequest)
	}

	_, err = rt.db.SendMessage(idConv, idUser, content, imgData, messageType, answerTo)
	if err != nil {
		http.Error(w, "Unable to read file data", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "File uploaded successfully")
}
