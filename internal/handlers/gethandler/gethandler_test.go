package gethandler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vadskev/urlshort/internal/entity"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

func TestNew(t *testing.T) {

	storeTest := memstorage.New()
	storeTest.Add(entity.Links{RawURL: "https://practicum.yandex.ru/", Slug: "bfebrehbreh"})
	storeTest.Add(entity.Links{RawURL: "https://yandex.ru/", Slug: "asdxvsdf"})

	type args struct {
		store URLStore
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		{
			name: "Test 1",
			args: args{store: storeTest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
