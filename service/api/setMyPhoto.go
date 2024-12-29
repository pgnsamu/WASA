package api

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	paramId := ps.ByName("id")
	id, err := strconv.Atoi(paramId)
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

	err = rt.db.SaveImageToDB(imgData, "users", "photo", id)
	if err != nil {
		if err.Error() == "id not found" {
			http.Error(w, "id not found", http.StatusBadRequest)
		} else {
			// fmt.Println(err.Error())
			http.Error(w, "Unable to save the image in the Database", http.StatusInternalServerError)
		}
		return
	}

	/*
		// Create a new file to save the uploaded image
		outFile, err := os.Create("uploaded_image.jpg") // Save the file as uploaded_image.jpg
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy the uploaded file to the new file
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error while copying the file", http.StatusInternalServerError)
			return
		}
	*/
	// Respond with a success message
	// TODO: capire se va bene
	io.WriteString(w, "File uploaded successfully")
	log.Println("File uploaded successfully")
	// fmt.Fprintf(w, "File uploaded successfully")
}
