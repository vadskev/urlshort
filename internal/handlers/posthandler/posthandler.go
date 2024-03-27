package posthandler

import (
	"errors"
	"io"
	"net/http"

	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/app"
	"github.com/vadskev/urlshort/internal/entity"
	"github.com/vadskev/urlshort/internal/util"
)

var (
	ErrReadRequest    = errors.New("POST fails to read request body")
	ErrInvalidUrl     = errors.New("POST invalid body url")
	ErrAddStore       = errors.New("POST storage add error")
	ErrFailedResponse = errors.New("POST Failed to write response")
	ErrMethodRequest  = errors.New("POST Is not POST Request")
)

type URLStore interface {
	Add(link entity.Links) (entity.Links, error)
}

func New(cfg *config.Config, store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, ErrMethodRequest.Error(), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, ErrReadRequest.Error(), http.StatusBadRequest)
			return
		}

		if !util.ValidateAddress(string(body)) {
			http.Error(w, ErrInvalidUrl.Error(), http.StatusBadRequest)
			return
		}

		shortCode := app.GenerateRandomString()

		_, err = store.Add(entity.Links{Slug: shortCode, RawURL: string(body)})
		if err != nil {
			http.Error(w, ErrAddStore.Error(), http.StatusBadRequest)
			return
		}

		url := cfg.BaseURL + "/" + shortCode

		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(url))
		if err != nil {
			http.Error(w, ErrFailedResponse.Error(), http.StatusInternalServerError)
			return
		}
	}
}
