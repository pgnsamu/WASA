package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}
	if id != claims["id"] {
		http.Error(w, "Utente non autorizzato", http.StatusUnauthorized)
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
	ph, err := rt.db.GetProfilePhoto(id)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"photo": ph,
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
	/*
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(imgData)
		if err != nil {
			http.Error(w, "Unable to write image data to response", http.StatusInternalServerError)
			return
		}

			_, err = io.WriteString(w, "File uploaded successfully")
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
				return
			}
	*/
	// fmt.Fprintf(w, "File uploaded successfully")
}
