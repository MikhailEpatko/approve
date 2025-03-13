package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewDB(cfg *AppConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.dbHost,
		cfg.dbPort,
		cfg.dbUser,
		cfg.dbPassword,
		cfg.dbName,
		cfg.dbSslMode,
	)
	db := sqlx.MustConnect(cfg.dbDriverName, dsn)
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
