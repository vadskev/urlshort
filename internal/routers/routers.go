package routers

import (
	"github.com/go-chi/chi"
	"github.com/vadskev/urlshort/internal/handlers"
)

const (
	createPostfix = "/"
	getPostfix    = "/{code}"
)

func Router(h *handlers.HandlerStorage) *chi.Mux {
	r := chi.NewRouter()

	r.Post(createPostfix, h.HandlerPost)
	r.Get(getPostfix, h.HandlerGet)

	return r
}
