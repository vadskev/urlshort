package memstorage

import (
	"reflect"
	"sync"
	"testing"

	"github.com/vadskev/urlshort/internal/entity"
)

func TestMemStorage_Add(t *testing.T) {
	store := MemStorage{data: sync.Map{}}

	tests := []struct {
		name    string
		keyData entity.Links
		want    entity.Links
		wantErr error
	}{
		{
			name: "Test add key",
			keyData: entity.Links{
				Slug:   "cvbdfy",
				RawURL: "https://test.com",
			},
			want: entity.Links{
				Slug:   "cvbdfy",
				RawURL: "https://test.com",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Add(tt.keyData)
			if err != nil {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	store := MemStorage{data: sync.Map{}}

	store.data.Store("sdflk", entity.Links{
		Slug:   "sdflk",
		RawURL: "https://example.com",
	})

	store.data.Store("cvbdfy", entity.Links{
		Slug:   "cvbdfy",
		RawURL: "https://test.com",
	})

	tests := []struct {
		name    string
		key     string
		want    entity.Links
		wantErr error
	}{
		{
			name: "Test key exists",
			key:  "cvbdfy",
			want: entity.Links{
				Slug:   "cvbdfy",
				RawURL: "https://test.com",
			},
			wantErr: entity.ErrSlugExists,
		},
		{
			name:    "Test the key no exists",
			key:     "sdfsdf",
			want:    entity.Links{},
			wantErr: entity.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got, _ := store.Get(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "Test 1",
			want: New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
