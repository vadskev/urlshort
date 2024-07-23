package save

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/storage"
	"github.com/vadskev/urlshort/internal/storage/filestorage"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestServeHTTP(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		query         string
		method        string
		reqBody       string
		statusWant    int
		wantEmptyBody bool
	}{
		{
			name:          "Normal POST (should work)",
			query:         "/",
			method:        http.MethodPost,
			statusWant:    http.StatusCreated,
			reqBody:       "https://practicum.yandex.ru/55",
			wantEmptyBody: false,
		},
	}

	conf := &config.Config{
		Storage:  config.Storage{FileStoragePath: "/tmp/short-url-db.json"},
		LogLevel: "info",
		HTTPServer: config.HTTPServer{
			ServerAddress: "localhost:8080",
			BaseURL:       "http://localhost:8080",
		},
	}

	// init logger
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	log := zap.Must(cfg.Build())

	var stor storage.Storage

	// init storage
	store := memstorage.NewMemStorage(log)

	filestore, _ := filestorage.NewFileStorage(conf.Storage.FileStoragePath, log)
	_ = filestore.Get(ctx, store)

	stor = filestore

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Post("/", New(log, conf, stor))
	})

	ts := httptest.NewServer(router)

	//tests run
	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, ts.URL+tt.query, strings.NewReader(tt.reqBody))
		require.NoError(t, err, tt.name)

		resp, err := ts.Client().Do(req)
		require.NoError(t, err, tt.name)
		assert.Equal(t, tt.statusWant, resp.StatusCode, tt.name)

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			require.NotEmpty(t, resp.Body, tt.name)
			response, err := io.ReadAll(resp.Body)
			require.NoError(t, err, tt.name)

			require.NotEmpty(t, resp.Body, tt.name)
			require.NotEmpty(t, response, tt.name)
		}
	}
}

func Test_JSON_ServeHTTP(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		query         string
		method        string
		reqBody       string
		statusWant    int
		wantEmptyBody bool
	}{
		{
			name:          "Normal POST (should work)",
			query:         "/api/shorten",
			method:        http.MethodPost,
			statusWant:    http.StatusCreated,
			reqBody:       "{\"url\":\"https://practicum.yandex.su\"}",
			wantEmptyBody: false,
		},
	}

	conf := &config.Config{
		Storage:  config.Storage{FileStoragePath: "/tmp/short-url-db.json"},
		LogLevel: "info",
		HTTPServer: config.HTTPServer{
			ServerAddress: "localhost:8080",
			BaseURL:       "http://localhost:8080",
		},
	}

	// init logger
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	log := zap.Must(cfg.Build())

	var stor storage.Storage
	// init storage
	store := memstorage.NewMemStorage(log)

	filestore, _ := filestorage.NewFileStorage(conf.Storage.FileStoragePath, log)
	_ = filestore.Get(ctx, store)

	stor = filestore

	router := chi.NewRouter()
	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", NewJSON(log, conf, stor))
	})

	ts := httptest.NewServer(router)

	//tests run
	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, ts.URL+tt.query, strings.NewReader(tt.reqBody))
		require.NoError(t, err, tt.name)

		resp, err := ts.Client().Do(req)
		require.NoError(t, err, tt.name)
		assert.Equal(t, tt.statusWant, resp.StatusCode, tt.name)

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			require.NotEmpty(t, resp.Body, tt.name)
			response, err := io.ReadAll(resp.Body)
			require.NoError(t, err, tt.name)

			require.NotEmpty(t, resp.Body, tt.name)
			require.NotEmpty(t, response, tt.name)

			result := struct {
				Result string `json:"result"`
			}{}
			require.NoError(t, json.Unmarshal(response, &result), "Error while unmarshalling json response")

		}
	}
}
