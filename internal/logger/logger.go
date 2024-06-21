package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

func New(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	zl, err := cfg.Build()

	if err != nil {
		return err
	}

	defer func() {
		err = zl.Sync()
	}()

	Log = zl
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			Log.Info("got incoming HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("length", ww.BytesWritten()),
				zap.Duration("time", time.Since(timeStart)),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
