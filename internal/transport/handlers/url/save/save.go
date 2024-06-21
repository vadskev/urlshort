package save

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadskev/urlshort/internal/config"
	resp "github.com/vadskev/urlshort/internal/lib/api/response"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/lib/random"
	"github.com/vadskev/urlshort/internal/lib/verify"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

type Request struct {
	Alias  string `json:"alias,omitempty"`
	ResURL string `json:"result,omitempty"`
	URL    string `json:"url"`
}

type Response struct {
	resp.Response
	Result string `json:"result"`
}

/**/

type URLSaver interface {
	SaveURL(data storage.URLData) error
}

func New(log *zap.Logger, cfg *config.Config, store URLSaver, fstore URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to read body", zp.Err(err))
			return
		}
		defer func() {
			err = r.Body.Close()
			if err != nil {
				return
			}
		}()

		// validate url
		if !verify.ValidateAddress(string(body)) {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("invalid url")
			return
		}

		var req Request

		// create random alias
		req.Alias = random.GenerateRandomString()

		req.URL = string(body)

		// create url
		req.ResURL = fmt.Sprintf("%s/%s", cfg.BaseURL, req.Alias)

		// add to store
		err = store.SaveURL(storage.URLData{URL: req.URL, ResURL: req.ResURL, Alias: req.Alias})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to add url store", zp.Err(err))
			return
		}

		// add to file store
		err = fstore.SaveURL(storage.URLData{URL: req.URL, ResURL: req.ResURL, Alias: req.Alias})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to add url fstore", zp.Err(err))
			return
		}

		// response OK
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(req.ResURL))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Info("failed to write response", zp.Err(err))
			return
		}
	}
}

func NewJSON(log *zap.Logger, cfg *config.Config, store URLSaver, fstore URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		// decode json request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to decode request", zp.Err(err))
			return
		}

		// validate url
		if !verify.ValidateAddress(req.URL) {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("invalid url")
			return
		}

		// create random alias
		req.Alias = random.GenerateRandomString()

		// create url
		req.ResURL = fmt.Sprintf("%s/%s", cfg.BaseURL, req.Alias)

		// add to store
		err = store.SaveURL(storage.URLData{URL: req.URL, ResURL: req.ResURL, Alias: req.Alias})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to add url", zp.Err(err))
			return
		}

		// add to file store
		err = fstore.SaveURL(storage.URLData{URL: req.URL, ResURL: req.ResURL, Alias: req.Alias})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("failed to add url fstore", zp.Err(err))
			return
		}

		// response OK
		responseOK(w, r, req.ResURL)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, result string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, Response{
		Result: result,
	})
}
