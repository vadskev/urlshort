package gethandler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadskev/urlshort/internal/entity"
)

var (
	ErrMethodRequest = errors.New("GET Is not GET Request")
	ErrURLNotFound   = errors.New("GET URL not found")
	ErrURLEmpty      = errors.New("GET URL empty")
)

type URLStore interface {
	Get(key string) (entity.Links, error)
}

func New(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, ErrMethodRequest.Error(), http.StatusBadRequest)
			return
		}

		shortCode := chi.URLParam(r, "code")

		if shortCode == "" {
			http.Error(w, ErrURLEmpty.Error(), http.StatusBadRequest)
			return
		}

		url, err := store.Get(shortCode)
		if err != nil {
			http.Error(w, ErrURLNotFound.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Location", url.RawURL)
	}
}
