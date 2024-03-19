package config

const (
	defaultHostServer = "localhost:8080"
)

type Config struct {
	HostServer   string
	BaseURLShort string
}

var config Config

func InitConfig() *Config {
	config = Config{
		HostServer:   defaultHostServer,
		BaseURLShort: "",
	}
	return &config
}

func GetConfig() *Config {
	return &config
}
