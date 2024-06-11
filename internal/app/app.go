package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/storage"
	"github.com/vadskev/urlshort/internal/storage/filestorage"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
	"github.com/vadskev/urlshort/internal/storage/postgres"
	"github.com/vadskev/urlshort/internal/transport/handlers/database/ping"
	"github.com/vadskev/urlshort/internal/transport/handlers/url/redirect"
	"github.com/vadskev/urlshort/internal/transport/handlers/url/save"
	"github.com/vadskev/urlshort/internal/transport/middleware/compress"
	"github.com/vadskev/urlshort/internal/transport/middleware/logger"
	"go.uber.org/zap"
)

func RunServer(log *zap.Logger, cfg *config.Config) error {
	const op = "internal.app.RunServer"

	ctx := context.Background()

	log.Info("Running server",
		zap.String("address", cfg.ServerAddress),
	)

	// init storage
	var stor storage.Storage
	cfg.DataBase.DatabaseDSN = ""
	switch {
	case cfg.DataBase.DatabaseDSN != "":
		// init postgresql storage
		dbstore, err := postgres.New(ctx, cfg, log)
		if err != nil {
			log.Info("Failed to init storage", zp.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
		err = dbstore.Setup(cfg)
		if err != nil {
			log.Info("Failed to migrate", zp.Err(err))
		}
		stor = dbstore
		defer dbstore.CloseStorage()
	default:
		// init mem storage
		memstore := memstorage.NewMemStorage(log)

		// init file storage
		filestore := filestorage.NewFileStorage(cfg.Storage.FileStoragePath, log)
		err := filestore.Get(ctx, memstore)
		if err != nil {
			log.Info("Error get file store", zp.Err(err))
		}
		stor = filestore
	}

	/**/

	// init router
	router := chi.NewRouter()

	// use middleware logger
	router.Use(logger.New(log))

	// use middleware compress
	router.Use(compress.New(log))

	// add url router
	router.Route("/", func(r chi.Router) {
		r.Post("/", save.New(log, cfg, stor))
	})

	// add url router json
	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", save.NewJSON(log, cfg, stor))
	})

	// add response router
	router.Route("/{code}", func(r chi.Router) {
		r.Get("/", redirect.New(log, stor))
	})

	// add ping router
	router.Route("/ping", func(r chi.Router) {
		r.Get("/", ping.New(log, stor))
	})

	err := http.ListenAndServe(cfg.ServerAddress, router)
	if err != nil {
		log.Info("Failed to start server")
	}

	return nil
}
