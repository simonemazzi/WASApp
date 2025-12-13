package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	autHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(autHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	username := params.ByName("username")
	if username == "" {
		dbUsers, err := rt.db.Users()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var users []User
		for _, u := range dbUsers {
			users = append(users, User{
				UserId:   u.UserId,
				Username: u.Username,
			})
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		dbUsers, err := rt.db.UsersBySearch(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var users []User
		for _, u := range dbUsers {
			users = append(users, User{
				UserId:   u.UserId,
				Username: u.Username,
			})
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
