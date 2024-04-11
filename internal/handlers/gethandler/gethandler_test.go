package gethandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/vadskev/urlshort/internal/entity"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

func TestNew(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "Test no url",
			requestMethod: http.MethodGet,
			requestPath:   "/hhjjj",
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
	}

	store := memstorage.New()

	_, err := store.Add(entity.Links{
		RawURL: "https://practicum.yandex.ru/",
		Slug:   "sdjfkh",
	})
	require.NoError(t, err)

	_, err = store.Add(entity.Links{
		RawURL: "https://yandex.ru/",
		Slug:   "asdxvsdf",
	})
	require.NoError(t, err)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			handler := New(store)
			r := chi.NewRouter()
			r.Get("/{code}", handler)

			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := http.Get(srv.URL + tt.requestPath)
			require.NoError(t, err)

			defer func() {
				err = resp.Body.Close()
				require.NoError(t, err)
			}()

			require.Equal(t, resp.StatusCode, tt.want.code)
		})
	}
}
