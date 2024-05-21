package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(alias string) (storage.URLData, error)
}

func New(log *zap.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get alias from url
		alias := chi.URLParam(r, "code")
		if alias == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("alias is empty")
			return
		}

		//get from store
		res, err := urlGetter.GetURL(alias)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("alias is empty", zp.Err(err))
			return
		}
		log.Info("store", zap.String("", res.URL))

		w.Header().Set("Location", res.URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
