package handlers

import (
	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/storage"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vadskev/urlshort/internal/app"
	"github.com/vadskev/urlshort/internal/util"
)

type HandlerStorage struct {
	ShortURLAddr string
	Store        storage.MemStorage
}

// HandlerPost
func (h *HandlerStorage) HandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		http.Error(w, "This is not a POST request, use POST request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		http.Error(w, "Reading POST body request failed", http.StatusBadRequest)
		return
	}

	if len(body) < 1 && util.IsURL(string(body)) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		http.Error(w, "Empty POST body request", http.StatusBadRequest)
		return
	}

	shortCode := app.GenerateRandomString()

	err = h.Store.AddURL(shortCode, string(body))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		http.Error(w, "URL not correct failed to shorten", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	url := "http://" + config.GetConfig().HostServer + "/" + shortCode
	_, err = w.Write([]byte(url))
	if err != nil {
		return
	}
}

// HandlerGet
func (h *HandlerStorage) HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")

	expandedURL, err := h.Store.GetURL(shortCode)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if expandedURL != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", expandedURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", "URL not found")
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}
