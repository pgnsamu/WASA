package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postGroupchatPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// Get the file from the form data
	file, _, err := r.FormFile("file") // "file" is the form field name
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
		if err.Error() == "id not found" {
			http.Error(w, "id not found", http.StatusBadRequest)
		} else {
			// fmt.Println(err.Error())
			http.Error(w, "Unable to save the image in the Database", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "File uploaded successfully")
}
