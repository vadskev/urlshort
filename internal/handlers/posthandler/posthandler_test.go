package posthandler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vadskev/urlshort/config"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg   *config.Config
		store URLStore
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {

		//req := httptest.NewRequest(http.MethodPost, "/", body)
		//w := httptest.NewRecorder()

		//handlers.HandlerPost(w, req)

		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
