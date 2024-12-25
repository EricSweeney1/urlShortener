package database

import (
	"database/sql"
	"github.com/Dashboard/urlShortener/config"
	_ "github.com/lib/pq"
)

func NewDB(config config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.DSN())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
