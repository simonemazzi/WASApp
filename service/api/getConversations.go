package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Avatar struct {
	Url    string `json:"url"`
	Mime   string `json:"mime"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Conversation struct {
	ConversationID string `json:"conversationId"`
	Name           string `json:"name"`
	Avatar         Avatar `json:"avatar"`
}

func (rt *_router) getConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {

	userId, err := strconv.Atoi(params.ByName("userId"))
	if err != nil {
		context.Logger.WithError(err).Error("Error converting userId to int")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !rt.db.IDExists(userId) {
		context.Logger.WithError(err).Error("User does not exist")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	convs, err := rt.db.GetConversations(userId)
	if err != nil {
		context.Logger.WithError(err).Error("Error getting conversations")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := make([]Conversation, 0, len(convs))
	for _, dbConv := range convs {
		response = append(response, Conversation{
			// Converti l'ID da int (DB) a string (API)
			ConversationID: strconv.Itoa(dbConv.ConversationID),
			Name:           dbConv.Name,
			// Mappa l'Avatar (anche se identico, sono tipi diversi per Go)
			Avatar: Avatar{
				Url:    dbConv.Avatar.Url,
				Mime:   dbConv.Avatar.Mime,
				Width:  dbConv.Avatar.Width,
				Height: dbConv.Avatar.Height,
			},
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ConversationsResponse{Conversations: response})
	if err != nil {
		context.Logger.WithError(err).Error("error encoding conversations")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
