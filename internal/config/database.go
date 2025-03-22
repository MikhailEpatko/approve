package config

import (
	"approve/internal/common"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDatabase(cfg *AppConfig) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.dbHost,
		cfg.dbPort,
		cfg.dbUser,
		cfg.dbPassword,
		cfg.dbName,
		cfg.dbSslMode,
	)
	DB = sqlx.MustConnect(cfg.dbDriverName, dsn)
	common.Logger.Info("connected to database")
}
