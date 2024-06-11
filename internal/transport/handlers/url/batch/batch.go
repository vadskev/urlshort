package batch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadskev/urlshort/internal/config"
	resp "github.com/vadskev/urlshort/internal/lib/api/response"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/lib/verify"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

type Request struct {
	Alias  string `json:"correlation_id"`
	ResURL string `json:"short_url,omitempty"`
	URL    string `json:"original_url"`
}
type Response struct {
	resp.Response
	Alias  string `json:"correlation_id"`
	ResURL string `json:"short_url,omitempty"`
}

type URLSaver interface {
	SaveBatchURL(ctx context.Context, data []storage.URLData) error
}

func New(log *zap.Logger, cfg *config.Config, store URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []Request

		// decode json request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to decode request", zp.Err(err))
			return
		}

		// validate url
		for k, v := range req {
			if !verify.ValidateAddress(v.URL) {
				w.WriteHeader(http.StatusBadRequest)
				log.Info("invalid url")
				return
			}
			// create url
			req[k].ResURL = fmt.Sprintf("%s/%s", cfg.BaseURL, v.Alias)
		}

		var data []storage.URLData
		for _, v := range req {
			data = append(data, storage.URLData{
				Alias:  v.Alias,
				URL:    v.URL,
				ResURL: v.ResURL,
			})
		}

		ctx := context.Background()
		err = store.SaveBatchURL(ctx, data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to add url store", zp.Err(err))
			return
		}

		responseOK(w, r, req)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req []Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var data []Response
	for _, v := range req {
		data = append(data, Response{
			Alias:  v.Alias,
			ResURL: v.ResURL,
		})
	}
	render.JSON(w, r, data)
}
