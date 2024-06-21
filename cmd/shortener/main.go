package main

import (
	"log"
	"net/http"

	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/routers"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

func main() {
	cfg := config.Load()
	store := memstorage.New()

	err := http.ListenAndServe(cfg.Server, routers.NewRouter(cfg, store))
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
