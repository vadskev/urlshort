package batch

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

type Request struct {
	Alias  string `json:"alias,omitempty"`
	ResURL string `json:"result,omitempty"`
	URL    string `json:"url"`
}

type URLSaver interface {
	SaveURL(ctx context.Context, data storage.URLData) error
}

func New(log *zap.Logger, cfg *config.Config, store URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		// decode json request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to decode request", zp.Err(err))
			return
		}

		log.Info("111", zap.Any("1", req))
	}
}
