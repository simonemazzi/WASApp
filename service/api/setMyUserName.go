package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type setUsernameRequest struct {
	NewUsername string `json:"newusername"`
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	autHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(autHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	userIdStr := params.ByName("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(400)
		_, err := w.Write([]byte("user id must be int"))
		if err != nil {
			return
		}
		return
	}
	var req setUsernameRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("invalid JSON body"))
		if err != nil {
			return
		}
		return
	}

	usernamePropost := req.NewUsername
	if usernamePropost == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		context.Logger.WithError(err).Warn("username must not be empty")
		return
	}
	if !rt.db.IDExists(userId) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		context.Logger.WithError(err).Warn("username does not exist")
		return
	}

	err = rt.db.SetMyUserName(userId, usernamePropost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		context.Logger.WithError(err).Warn("failed to set user name")
		return
	}
	w.WriteHeader(http.StatusOK)
	var u User
	u.UserId = strconv.Itoa(userId)
	u.Username = usernamePropost

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		context.Logger.WithError(err).Warn("failed to encode user response")
		return

	}

}
