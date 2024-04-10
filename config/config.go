package config

import (
	"flag"
	"os"
)

const (
	defaultServer      = "localhost:8080"
	defaultBaseURL     = "http://localhost:8080"
	defaultStoragePath = "./tmp/short-url-db.json"
)

type Config struct {
	Server          string
	BaseURL         string
	LogLevel        string
	FileStoragePath string
}

func Load() *Config {
	cfg := &Config{
		Server:          defaultServer,
		BaseURL:         defaultBaseURL,
		LogLevel:        "info",
		FileStoragePath: defaultStoragePath,
	}

	// get env
	if envBaseURLShortener := os.Getenv("SERVER_ADDRESS"); envBaseURLShortener != "" {
		cfg.Server = envBaseURLShortener
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	// get flag
	flag.StringVar(&cfg.Server, "a", "localhost:8080", "server address; example: -a localhost:8080")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "short url base; example: -b https://yandex.ru")
	flag.StringVar(&cfg.FileStoragePath, "f", "./tmp/short-url-db.json", "file storage path; example: -f /tmp/short-url-db.json")
	flag.Parse()

	return cfg
}
