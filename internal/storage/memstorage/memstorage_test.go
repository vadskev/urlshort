package memstorage

import (
	"log"
	"reflect"
	"sync"
	"testing"

	"github.com/vadskev/urlshort/internal/entity"
)

func TestMemStorage_Add(t *testing.T) {
	type fields struct {
		data sync.Map
	}
	type args struct {
		link entity.Links
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Links
		wantErr bool
	}{
		{
			name:    "Test add",
			fields:  fields{data: sync.Map{}},
			args:    args{link: entity.Links{Slug: "bfebrehbreh", RawURL: "https://practicum.yandex.ru/"}},
			want:    entity.Links{Slug: "bfebrehbreh", RawURL: "https://practicum.yandex.ru/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("%#v", tt)
			store := &MemStorage{
				data: tt.fields.data,
			}
			got, err := store.Add(tt.args.link)
			if (err != nil) != tt.wantErr {
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
	type fields struct {
		data sync.Map
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Links
		wantErr bool
	}{
		{
			name:    "Test add",
			fields:  fields{data: sync.Map{}},
			args:    args{"bfebrehbreh"},
			want:    entity.Links{Slug: "bfebrehbreh", RawURL: "https://practicum.yandex.ru/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemStorage{
				data: tt.fields.data,
			}

			store.Add(entity.Links{Slug: "bfebrehbreh", RawURL: "https://practicum.yandex.ru/"})

			got, err := store.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
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
			name: "Test new",
			want: &MemStorage{data: sync.Map{}},
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
