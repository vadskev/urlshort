package migrator

import (
	"database/sql"

	"github.com/pressly/goose"
	"github.com/vadskev/urlshort/internal/config"
)

func Migrate(cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.DataBase.DatabaseDSN)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	return goose.Up(db, "./migrations/")
}
