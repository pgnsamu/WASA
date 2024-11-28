package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	name, err := rt.db.GetUsersDB()
	if err != nil {
		return 
	}
	json.NewEncoder(w).Encode(name)

}
