package posthandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vadskev/urlshort/config"
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
				code:        201,
				contentType: "text/plain",
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

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			handler := New(cfg, store)
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.inputLink)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tt.want.code)
			require.Equal(t, rr.Result().Header.Get("content-type"), tt.want.contentType)

			require.NoError(t, err)

			err = rr.Result().Body.Close()
			if err != nil {
				require.NoError(t, err)
			}

			err = req.Body.Close()
			if err != nil {
				require.NoError(t, err)
			}
		})
	}
}
