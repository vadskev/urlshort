package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vadskev/urlshort/internal/config"
	"github.com/vadskev/urlshort/internal/lib/migrator"
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

func (d *DBStorage) Setup(cfg *config.Config) error {
	const op = "storage.postgres.Setup"
	err := migrator.Migrate(cfg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (d *DBStorage) GetURL(ctx context.Context, alias string) (storage.URLData, error) {
	const op = "storage.postgres.GetURL"

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	data := storage.URLData{}

	if err := d.db.QueryRow(ctx, `SELECT url, alias, res_url FROM urls WHERE alias = $1`, alias).Scan(&data.URL, &data.Alias, &data.ResURL); err != nil {
		fmt.Println("Error occur while finding user: ", err)
		return storage.URLData{}, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (d *DBStorage) SaveURL(ctx context.Context, data storage.URLData) error {
	const op = "storage.postgres.SaveURL"
	//stmt := `INSERT INTO urls (url, alias, res_url) VALUES($1, $2, $3)`
	stmt := `INSERT INTO urls (url, alias, res_url)
			VALUES ($1, $2, $3)
			ON CONFLICT (url)
			DO NOTHING RETURNING url;`

	v, err := d.db.Exec(ctx, stmt, data.URL, data.Alias, data.ResURL)
	if v.String() == "INSERT 0 0" {
		return errors.New("url exists")
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
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

func (d *DBStorage) SaveBatchURL(ctx context.Context, data []storage.URLData) error {
	const op = "storage.postgres.SaveBatchURL"

	tx, err := d.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	d.log.Info("Begin transaction")

	for _, v := range data {
		_, err = tx.Exec(ctx, `INSERT INTO urls (url, alias, res_url) VALUES($1, $2, $3)`, v.URL, v.Alias, v.ResURL)
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	d.log.Info("Commit transaction")

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *DBStorage) CloseStorage() {
	d.db.Close()
}
