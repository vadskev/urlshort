package main

import (
	"net/http"

	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/routers"
)

func main() {
	cfg := config.InitConfig()
	cfg.ParseFlags()

	config := config.GetConfig()

	err := http.ListenAndServe(config.HostServer, routers.Router())
	if err != nil {
		panic(err)
	}
}
