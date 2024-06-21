package routers

import (
	"github.com/go-chi/chi"
	"github.com/vadskev/urlshort/internal/handlers"
)

const (
	createPostfix = "/"
	getPostfix    = "/{code}"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post(createPostfix, handlers.HandlerPost)
	r.Get(getPostfix, handlers.HandlerGet)

	return r
}
