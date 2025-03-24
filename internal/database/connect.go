package database

import (
	"approve/internal/common"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

var DB *sqlx.DB

type DbConfig struct {
	dbDriverName string
	dbHost       string
	dbPort       string
	dbUser       string
	dbPassword   string
	dbName       string
	dbSslMode    string
}

func Connect() {
	cfg := getConfig()
	ConnectWithCfg(cfg)
}

func ConnectWithCfg(cfg *DbConfig) {
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

func getConfig() *DbConfig {
	err := godotenv.Load(".env")
	if err != nil {
		common.Logger.Info("not found .env file")
	}
	return &DbConfig{
		dbDriverName: os.Getenv("DB_DRIVER_NAME"),
		dbHost:       os.Getenv("DB_HOST"),
		dbPort:       os.Getenv("DB_PORT"),
		dbUser:       os.Getenv("DB_USER"),
		dbPassword:   os.Getenv("DB_PASSWORD"),
		dbName:       os.Getenv("DB_NAME"),
		dbSslMode:    os.Getenv("DB_SSL_MODE"),
	}
}
