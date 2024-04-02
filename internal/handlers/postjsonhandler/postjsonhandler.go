package postjsonhandler

import (
	"encoding/json"
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
	ErrInvalidURL     = errors.New("POST invalid body url")
	ErrAddStore       = errors.New("POST storage add error")
	ErrFailedResponse = errors.New("POST Failed to write response")
)

type Response struct {
	URL string `json:"result"`
}

type URLStore interface {
	Add(link entity.Links) (entity.Links, error)
}

func New(cfg *config.Config, store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var link entity.Links
		var res Response

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

		err = link.UnmarshalJSON(body)
		if err != nil {
			http.Error(w, ErrReadRequest.Error(), http.StatusBadRequest)
			return
		}

		if !util.ValidateAddress(link.RawURL) {
			http.Error(w, ErrInvalidURL.Error(), http.StatusBadRequest)
			return
		}

		shortCode := app.GenerateRandomString()
		link.Slug = shortCode

		_, err = store.Add(link)
		if err != nil {
			http.Error(w, ErrAddStore.Error(), http.StatusBadRequest)
			return
		}

		res.URL = cfg.BaseURL + "/" + shortCode

		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, ErrFailedResponse.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, ErrFailedResponse.Error(), http.StatusInternalServerError)
			return
		}
	}
}
