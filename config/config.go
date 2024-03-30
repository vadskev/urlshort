package config

import (
	"flag"
	"os"
)

const (
	defaultServer      = "localhost:8080"
	defaultBaseURL     = "http://localhost:8080"
	defaultStoragePath = "./"
)

type Config struct {
	Server      string
	BaseURL     string
	StoragePath string
	LogLevel    string
}

func Load() *Config {
	cfg := &Config{
		Server:   defaultServer,
		BaseURL:  defaultBaseURL,
		LogLevel: "info",
	}

	// get env
	if envBaseURLShortener := os.Getenv("SERVER_ADDRESS"); envBaseURLShortener != "" {
		cfg.Server = envBaseURLShortener
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}

	// get flag
	flag.StringVar(&cfg.Server, "a", "localhost:8080", "server address; example: -a localhost:8080")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "short url base; example: -b https://yandex.ru")
	flag.Parse()

	return cfg
}
