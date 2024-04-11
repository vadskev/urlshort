package main

import (
	"log"
	"net/http"

	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/logger"
	"github.com/vadskev/urlshort/internal/routers"
	"github.com/vadskev/urlshort/internal/storage/filestorage"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	store := memstorage.New()

	fstore := filestorage.New(cfg.FileStoragePath)
	if err := fstore.Load(store); err != nil {
		log.Println("Error load fstore")
	}

	if err := logger.New(cfg.LogLevel); err != nil {
		log.Fatal("Error logger")
	}

	logger.Log.Info("Running server", zap.String("address", cfg.Server))

	err := http.ListenAndServe(cfg.Server, routers.NewRouter(cfg, store, fstore))
	if err != nil {
		logger.Log.Info("Failed to start server")
	}
}
