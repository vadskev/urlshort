package main

import (
	"fmt"
	"github.com/vadskev/urlshort/internal/storage"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vadskev/urlshort/internal/handlers"
)

func TestHandlerPost(t *testing.T) {

	type expected struct {
		statusCode int
	}
	tests := []struct {
		name   string
		body   string
		status expected
	}{
		{
			name: "Test 1. Empty body",
			body: "",
			status: expected{
				statusCode: 201,
			},
		},
		{
			name: "Test 2. Body correct data",
			body: "https://practicum.yandex.ru/",
			status: expected{
				statusCode: 201,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(test.body)

			log.Printf("%v: Body: %v\n", test.name, test.body)

			req := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			hStore := &handlers.HandlerStorage{
				ShortURLAddr: "http://localhost:8080",
				Store:        *storage.NewMemStorage(),
			}
			hStore.HandlerPost(w, req)
			res := w.Result()

			defer res.Body.Close()

			fmt.Printf("Expected status code: %d status code: %d\n", test.status.statusCode, res.StatusCode)
			require.Equal(t, test.status.statusCode, res.StatusCode)
		})
	}
}

func TestHandlerGet(t *testing.T) {
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name   string
		url    string
		status expected
	}{
		{
			name: "Test 1. Empty Url",
			url:  "/",
			status: expected{
				statusCode: 400,
			},
		},
		{
			name: "Test 2. Empty body data",
			url:  "http://localhost:8080/",
			status: expected{
				statusCode: 400,
			},
		},
		{
			name: "Test 3. Missing in the database",
			url:  "http://localhost:8080/F4g44a",
			status: expected{
				statusCode: 400,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			log.Printf("%v: url: %v\n", test.name, test.url)

			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()

			hStore := &handlers.HandlerStorage{
				ShortURLAddr: "http://localhost:8080",
				Store:        *storage.NewMemStorage(),
			}
			hStore.HandlerGet(w, req)

			res := w.Result()

			defer res.Body.Close()

			fmt.Printf("Expected status code: %d status code: %d\n", test.status.statusCode, res.StatusCode)
			require.Equal(t, test.status.statusCode, res.StatusCode)
		})
	}
}
