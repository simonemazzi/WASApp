package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// SessionRequest : Request's schema
type SessionRequest struct {
	Name string `json:"name"`
}

type SessionResponse struct {
	UserId string    `json:"userId"`
	Token  string    `json:"token"`
	Time   time.Time `json:"time"`
}

func (rt *_router) postSession(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "text/plain")

	// Reading body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		context.Logger.WithError(err).Error("Error reading body")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid body"))
		return
	}

	// Parsing JSON
	var req SessionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		context.Logger.WithError(err).Error("Error parsing JSON")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid JSON"))
		return
	}

	var resp SessionResponse
	// Execute query and return response
	resp.UserId, resp.Token, resp.Time, err = rt.db.CreateSession(req.Name)
	fmt.Println("UserID:", resp.UserId)
	fmt.Println("Token:", resp.Token)
	fmt.Println("Time:", resp.Time)
	if err != nil {
		fmt.Println("Failed to create session:", err)
		http.Error(w, "Failed to create session1", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResp)
}
