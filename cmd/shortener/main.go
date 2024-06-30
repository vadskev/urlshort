package main

import (
	"os"

	"github.com/vadskev/urlshort/internal/app"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := setupLogger()

	if err := app.RunServer(log, cfg); err != nil {
		log.Info("error to start server", zp.Err(err))
	} else {
		log.Info("server was shutdown")
	}
}

func setupLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	return zap.Must(cfg.Build())
}
