package ping

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type DBStorage interface {
	Ping(ctx context.Context) error
}

func New(log *zap.Logger, dbStorage DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := dbStorage.Ping(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Info("ping", zap.String("connection", "error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("ping", zap.String("connection", "successful"))
	}
}
