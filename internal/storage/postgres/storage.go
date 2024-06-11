package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

var _ storage.Storage = (*DBStorage)(nil)

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*DBStorage, error) {
	const op = "storage.postgres.NewStorage"
	pgConnString := cfg.DataBase.DatabaseDSN

	db, err := pgxpool.New(ctx, pgConnString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &DBStorage{db: db, log: log}, nil
}

func (d *DBStorage) GetURL(alias string) (storage.URLData, error) {
	//TODO implement me
	return storage.URLData{}, nil
}

func (d *DBStorage) SaveURL(data storage.URLData) error {
	//TODO implement me
	return nil
}

func (d *DBStorage) Ping(ctx context.Context) error {
	d.log.Info("Check status connection")
	c, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := d.db.Ping(c); err != nil {
		d.log.Info("Error connect to DataBase")
		return err
	}
	d.log.Info("Connect to DataBase successful")
	return nil
}

func (d *DBStorage) CloseStorage() {
	d.db.Close()
}
