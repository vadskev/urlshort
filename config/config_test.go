package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Test",
			want: &Config{
				Server:          "localhost:8080",
				BaseURL:         "http://localhost:8080",
				LogLevel:        "info",
				FileStoragePath: "./tmp/short-url-db.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
