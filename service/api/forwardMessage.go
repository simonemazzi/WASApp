package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// TODO: SIstemare roba degli id dei messaggi inoltrati
type ForwardMessageRequest struct {
	ForwardToConversation []int `json:"forwardToConversation"`
	ForwardToGroup        []int `json:"forwardToGroup"`
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// --- userId ---
	userId, err := strconv.Atoi(params.ByName("userId"))
	if err != nil {
		context.Logger.WithError(err).Error("Error converting userId to int")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !rt.db.IDExists(userId) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// --- conversationId (origine) ---
	conversationId, err := strconv.Atoi(params.ByName("conversationId"))
	if err != nil {
		context.Logger.WithError(err).Error("Error converting conversationId")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	isThere, err := rt.db.UserConversation(userId, conversationId)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking user conversation")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !isThere {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// --- messageId ---
	messageId, err := strconv.Atoi(params.ByName("messageId"))
	if err != nil {
		context.Logger.WithError(err).Error("Error converting messageId")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	originalMessageId := messageId

	// --- body ---
	var request ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		context.Logger.WithError(err).Error("Error decoding request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// --- forward verso conversazioni ---
	for _, convId := range request.ForwardToConversation {
		isThere, err := rt.db.UserConversation(userId, convId)
		if err != nil {
			context.Logger.WithError(err).Error("Error checking user conversation")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !isThere {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = rt.db.ForwardToConversation(userId, convId, originalMessageId)
		if err != nil {
			context.Logger.WithError(err).Error("Error forwarding to conversation")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// --- forward verso gruppi ---
	dateHeader := r.Header.Get("Date")
	if dateHeader == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	t, err := time.Parse(time.RFC1123, dateHeader)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	date := t.UTC().Format("2006-01-02 15:04:05")

	for _, groupId := range request.ForwardToGroup {
		isThere, err := rt.db.UserGroup(userId, groupId, date)
		if err != nil {
			context.Logger.WithError(err).Error("Error checking user group")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !isThere {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = rt.db.ForwardToGroup(userId, groupId, originalMessageId)
		if err != nil {
			context.Logger.WithError(err).Error("Error forwarding to group")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
