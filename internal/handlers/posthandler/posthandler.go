package posthandler

import (
	"errors"
	"io"
	"net/http"

	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/app"
	"github.com/vadskev/urlshort/internal/entity"
	"github.com/vadskev/urlshort/internal/storage/filestorage"
	"github.com/vadskev/urlshort/internal/util"
)

var (
	ErrReadRequest    = errors.New("POST fails to read request body")
	ErrInvalidURL     = errors.New("POST invalid body url")
	ErrAddStore       = errors.New("POST storage add error")
	ErrFailedResponse = errors.New("POST Failed to write response")
	ErrMethodRequest  = errors.New("POST Is not POST Request")
)

type URLStore interface {
	Add(link entity.Links) (entity.Links, error)
}

func New(cfg *config.Config, store URLStore, fstore *filestorage.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, ErrReadRequest.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			err = r.Body.Close()
			if err != nil {
				return
			}
		}()

		if !util.ValidateAddress(string(body)) {
			http.Error(w, ErrInvalidURL.Error(), http.StatusBadRequest)
			return
		}

		shortCode := app.GenerateRandomString()

		link, err := store.Add(entity.Links{Slug: shortCode, RawURL: string(body)})
		if err != nil {
			http.Error(w, ErrAddStore.Error(), http.StatusBadRequest)
			return
		}

		url := cfg.BaseURL + "/" + shortCode

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(url))
		if err != nil {
			http.Error(w, ErrFailedResponse.Error(), http.StatusInternalServerError)
			return
		}

		err = fstore.SaveToFileStorage(&link)
		if err != nil {
			http.Error(w, ErrFailedResponse.Error(), http.StatusInternalServerError)
			return
		}
	}
}
