package api

import (
	"encoding/json"

	"net/http"

	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type GroupRequest struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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

	date := r.Header.Get("Date")
	if date == "" {
		http.Error(w, "Missing Date header", http.StatusBadRequest)
		return
	}
	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		context.Logger.WithError(err).Error("error parsing date header")
		return
	}
	timestamp := t.UTC().Format("2006-01-02 15:04:05")

	var group GroupRequest
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		context.Logger.WithError(err).Error("Error decoding body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	groupR, err := rt.db.CreateGroup(userId, group.Name, group.Participants, timestamp)

	if err != nil {
		context.Logger.WithError(err).Error("Error creating group")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//defaultPath := "uploads/default/default-avatar-profile-icon-social-600nw-1906669723.png"

	//filename := "groupPhoto" + context.ReqUUID.String() + ".png"
	// path := filepath.Join("uploads", "groups", strconv.Itoa(groupR.GroupId), filename)

	//if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
	//		context.Logger.WithError(err).Error("Error creating group directory")
	//		http.Error(w, "Cannot create directory", http.StatusInternalServerError)
	//		return
	//	}

	//src, err := os.Open(defaultPath)
	//	if err != nil {
	//		context.Logger.WithError(err).Error("Error opening default image")
	//		http.Error(w, "Default image missing", http.StatusInternalServerError)
	//		return
	//	}
	//	defer func(src *os.File) {
	//		err := src.Close()
	//		if err != nil {
	//			context.Logger.WithError(err).Error("Error closing default image")
	//			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//			return
	//		}
	//	}(src)

	//dst, err := os.Create(path)
	//	if err != nil {
	//		context.Logger.WithError(err).Error("Error creating image file")
	//		http.Error(w, "Cannot save image", http.StatusInternalServerError)
	//		return
	//	}
	//	defer func(dst *os.File) {
	//		err := dst.Close()
	//		if err != nil {
	//			context.Logger.WithError(err).Error("Error closing file")
	//			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//			return
	//		}
	//	}(dst)
	//
	//	if _, err := io.Copy(dst, src); err != nil {
	//		context.Logger.WithError(err).Error("Error copying default image")
	//		http.Error(w, "Write failed", http.StatusInternalServerError)
	//		return
	//	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(groupR)
	if err != nil {
		context.Logger.WithError(err).Error("Error encoding group")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
