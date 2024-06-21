package config

import (
	"flag"
	"os"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultBaseURL       = "http://localhost:8080"
	defaultStoragePath   = "/tmp/short-url-db.json"
	defaultLogLevel      = "info"
)

type Config struct {
	HTTPServer
	LogLevel string
	Storage  Storage
}

type HTTPServer struct {
	ServerAddress string
	BaseURL       string
}

type Storage struct {
	FileStoragePath string
}

func MustLoad() *Config {
	var cfg Config

	cfg.ServerAddress = defaultServerAddress
	cfg.BaseURL = defaultBaseURL
	cfg.Storage.FileStoragePath = defaultStoragePath
	cfg.LogLevel = defaultLogLevel

	// get flag
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address; example: -a localhost:8080")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "short url base; example: -b https://yandex.ru")
	flag.StringVar(&cfg.Storage.FileStoragePath, "f", "/tmp/short-url-db.json", "file storage path; example: -f /tmp/short-url-db.json")
	flag.Parse()

	// get env
	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		cfg.ServerAddress = envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.Storage.FileStoragePath = envFileStoragePath
	}

	return &cfg
}
