package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Group struct {
	GroupId      int    `json:"groupId"`
	Name         string `json:"name"`
	Photo        Upload `json:"photo"`
	Participants []User `json:"participants"`
}

func (rt *_router) viewGroups(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	userId, err := strconv.Atoi(params.ByName("userId"))
	if err != nil {
		context.Logger.WithError(err).Error("Error converting userId to int")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	groups, err := rt.db.GetGroups(userId)
	if err != nil {
		context.Logger.WithError(err).Error("Error getting groups")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	apiGroups := make([]Group, 0, len(groups))
	for _, dbGroup := range groups {
		var apiParticipants []User
		for _, dbUser := range dbGroup.Participants {
			apiParticipants = append(apiParticipants, User{
				UserId:   dbUser.UserID,
				Username: dbUser.Username,
			})
		}
		apiGroups = append(apiGroups, Group{
			GroupId: dbGroup.GroupId,
			Name:    dbGroup.Name,
			Photo: Upload{
				Url:    dbGroup.Photo.Url,
				Mime:   dbGroup.Photo.Mime,
				Width:  dbGroup.Photo.Width,
				Height: dbGroup.Photo.Height,
			},
			Participants: apiParticipants,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(GroupsResponse{Groups: apiGroups})
	if err != nil {
		context.Logger.WithError(err).Error("Error encoding group")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
