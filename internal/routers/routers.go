package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/handlers/gethandler"
	"github.com/vadskev/urlshort/internal/handlers/posthandler"
	"github.com/vadskev/urlshort/internal/handlers/postjsonhandler"
	"github.com/vadskev/urlshort/internal/logger"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

const (
	postPostfix    = "/"
	getPostfix     = "/{code}"
	getJsonPostfix = "/api/shorten"
)

func NewRouter(cfg *config.Config, store *memstorage.MemStorage) *chi.Mux {
	router := chi.NewRouter()
	router.Use(logger.RequestLogger)
	router.Post(postPostfix, posthandler.New(cfg, store))
	router.Get(getPostfix, gethandler.New(store))
	router.Get(getPostfix, gethandler.New(store))
	router.Post(getJsonPostfix, postjsonhandler.New(cfg, store))
	return router
}
