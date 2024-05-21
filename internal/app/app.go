package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/storage/filestorage"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
	"github.com/vadskev/urlshort/internal/transport/handlers/url/redirect"
	"github.com/vadskev/urlshort/internal/transport/handlers/url/save"
	"github.com/vadskev/urlshort/internal/transport/middleware/compress"
	"github.com/vadskev/urlshort/internal/transport/middleware/logger"
	"go.uber.org/zap"
)

func RunServer(log *zap.Logger, cfg *config.Config) error {
	log.Info("Running server",
		zap.String("address", cfg.ServerAddress),
	)

	// init storage
	store := memstorage.NewMemStorage(log)

	// init file storage
	filestore := filestorage.NewFileStorage(cfg.Storage.FileStoragePath, log)
	err := filestore.Get(store)
	if err != nil {
		log.Info("Error get file store", zp.Err(err))
	}

	// init router
	router := chi.NewRouter()

	// use middleware logger
	router.Use(logger.New(log))

	// use middleware compress
	router.Use(compress.New(log))

	// add url router
	router.Route("/", func(r chi.Router) {
		r.Post("/", save.New(log, cfg, store, filestore))
	})

	// add url router json
	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", save.NewJSON(log, cfg, store, filestore))
	})

	// add response router
	router.Route("/{code}", func(r chi.Router) {
		r.Get("/", redirect.New(log, store))
	})

	err = http.ListenAndServe(cfg.ServerAddress, router)
	if err != nil {
		log.Info("Failed to start server")
	}

	return nil
}
