package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	paramId := ps.ByName("id")

	// conversione string to int
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Errore id non intero", http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUserInfo(id)

	if user == nil && err.Error() == "user not found" {
		http.Error(w, "Errore id non registrato", http.StatusBadRequest)
		return
	} else if user == nil && err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
