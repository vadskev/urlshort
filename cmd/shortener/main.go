package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

const (
	baseURL       = "localhost:8080"
	createPostfix = "/"
	getPostfix    = "/{code}"
)

var urlDBase map[string]string

// handlerPost
func handlerPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if len(body) > 0 && isURL(string(body)) {
		shortURL := generateShortURL()

		// записываем в базу, перезаписав
		urlDBase[shortURL] = string(body)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + shortURL))

	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte("No URL in request"))
	}
}

// handlerGet
func handlerGet(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.String()[1:]

	if value, ok := urlDBase[shortURL]; ok {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", value)
		w.WriteHeader(307)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", "URL not found")
		w.WriteHeader(400)
	}
}

// generateShortURL Алгоритм сокращения URL
func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := make([]byte, 7)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}

// isURL проверяем валидность URL
func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func main() {
	urlDBase = make(map[string]string)

	r := chi.NewRouter()
	log.Println("Server is starter")

	r.Post(createPostfix, handlerPost)
	r.Get(getPostfix, handlerGet)

	err := http.ListenAndServe(baseURL, r)
	if err != nil {
		panic(err)
	}
}
