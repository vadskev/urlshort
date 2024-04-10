package postjsonhandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/storage/filestorage"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

func TestNew(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name      string
		inputLink string
		want      want
	}{
		{
			name:      "Test normal",
			inputLink: "https://practicum.yandex.ru/",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "Test wrong url",
			inputLink: "sdff",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "Test no url",
			inputLink: "",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	cfg := config.Load()
	store := memstorage.New()
	fstore, _ := filestorage.New(cfg.FileStoragePath)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.inputLink))
			w := httptest.NewRecorder()

			handler := New(cfg, store, fstore)
			handler(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("content-type"))

			defer func() {
				err := res.Body.Close()
				require.NoError(t, err)
			}()
		})
	}
}
