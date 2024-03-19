package config

import "flag"

func (config *Config) ParseFlags() {
	flag.StringVar(&config.HostServer, "a", config.HostServer, "server address; example: -a localhost:8080")
	flag.StringVar(&config.BaseURLShort, "b", config.BaseURLShort, "short url base; example: -b https://yandex.ru")
	flag.Parse()
}
