package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type requestBody struct {
	Username string `json:"username"`
}

// non ritorna errori
// non ritorna valori

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Leggo il body della request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusInternalServerError)
		return
	}
	var reqBody requestBody
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := rt.db.DoLogin(reqBody.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var token string
	// Genero il token
	tok, erro := GenerateJWT(reqBody.Username, *id)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}
	token = tok

	claims, err := ValidateJWT(token)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	idF, ok := claims["id"].(int)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	user, err := rt.db.GetUserInfo(idF)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"token":    token,
		"user_id":  user.Id,
		"username": user.Username,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}

	_, err = rt.db.ReceiveAllMessages(idF)
	if err != nil {
		return
	}

}
