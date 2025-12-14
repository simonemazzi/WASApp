package api

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Send struct {
	Text  string `json:"text"`
	Photo Upload `json:"photo"`
}

const JPEG = "image/jpeg"
const PNG = "image/png"

func (rt *_router) postMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	userId, err := strconv.Atoi(params.ByName("userId"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !rt.db.IDExists(userId) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	conversationId, err := strconv.Atoi(params.ByName("conversationId"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	isThere, err := rt.db.UserConversation(userId, conversationId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !isThere {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	contentType := r.Header.Get("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/")

	const maxSize = 10 << 20 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	var text string
	var photoId *int

	// -------- CASE 1: JSON SENZA FOTO --------
	if !isMultipart {
		var payload struct {
			Body string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		text = payload.Body
	} else {
		// -------- CASE 2: MULTIPART (TESTO + FACOLTATIVA FOTO) --------
		if err := r.ParseMultipartForm(maxSize); err != nil {
			http.Error(w, "Invalid multipart form", http.StatusBadRequest)
			return
		}

		text = r.FormValue("body")

		file, header, err := r.FormFile("photo")

		if err == nil {

			mime := header.Header.Get("Content-Type")

			if mime != JPEG && mime != PNG {
				http.Error(w, "Invalid image type", http.StatusBadRequest)
				return
			}

			ext := ".png"
			if mime == JPEG {
				ext = ".jpg"
			}

			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					context.Logger.WithError(err).Error("closing file")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}(file)

			path := fmt.Sprintf("uploads/conversations/%d_%d_%s%s",
				conversationId, userId, context.ReqUUID, ext,
			)

			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				http.Error(w, "Cannot create directory", http.StatusInternalServerError)
				return
			}

			out, err := os.Create(path)
			if err != nil {
				http.Error(w, "Cannot save image", http.StatusInternalServerError)
				return
			}
			defer func(out *os.File) {
				err := out.Close()
				if err != nil {
					http.Error(w, "Cannot save image", http.StatusInternalServerError)
					context.Logger.Errorf("Cannot save image: %s", err)
				}
			}(out)

			if _, err := io.Copy(out, file); err != nil {
				http.Error(w, "Write failed", http.StatusInternalServerError)
				return
			}
			imgFile, err := os.Open(path)
			if err != nil {
				context.Logger.WithError(err).WithField("path", path).Error("Cannot open file")
				http.Error(w, "Cannot save image", http.StatusInternalServerError)
				return
			}

			img, _, err := image.Decode(imgFile)

			if err != nil {
				context.Logger.WithError(err).Error("Error decoding photo file")
				http.Error(w, "Cannot decode image", http.StatusInternalServerError)
				return
			}
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			id, err := rt.db.InsertPhoto(path, width, height, mime)
			if err != nil {
				http.Error(w, "Cannot save photo to DB", http.StatusInternalServerError)
				return
			}

			photoId = &id
		}
	}

	// ALMENO TESTO O FOTO
	if text == "" && photoId == nil {
		http.Error(w, "Message must have text or photo", http.StatusBadRequest)
		return
	}

	// SALVA MESSAGGIO
	if err := rt.db.InsertMessage(conversationId, userId, text, photoId); err != nil {
		http.Error(w, "Cannot save message", http.StatusInternalServerError)
		return
	}

	// RISPONDE CON TUTTI I MESSAGGI
	messages, err := rt.db.GetMessages(conversationId, userId)
	if err != nil {
		context.Logger.WithError(err).Error("Cannot get messages")
		http.Error(w, "Error getting image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		context.Logger.WithError(err).Error("Cannot save message")
		return
	}

}
