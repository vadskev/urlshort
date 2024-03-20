package main

import (
	"net/http"

	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/handlers"
	"github.com/vadskev/urlshort/internal/routers"
	"github.com/vadskev/urlshort/internal/storage"
)

func main() {
	cfg := config.InitConfig()
	cfg.ParseFlags()

	config := config.GetConfig()
	store := storage.NewMemStorage()

	hStore := &handlers.HandlerStorage{
		ShortURLAddr: cfg.BaseURLShort,
		Store:        *store,
	}

	err := http.ListenAndServe(config.HostServer, routers.Router(hStore))
	if err != nil {
		panic(err)
	}
}
