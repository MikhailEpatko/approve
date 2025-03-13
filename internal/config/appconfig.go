package config

import (
	_ "github.com/lib/pq"
	"os"
)

type AppConfig struct {
	dbDriverName string
	dbHost       string
	dbPort       string
	dbUser       string
	dbPassword   string
	dbName       string
	dbSslMode    string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		dbDriverName: os.Getenv("DB_DRIVER_NAME"),
		dbHost:       os.Getenv("DB_HOST"),
		dbPort:       os.Getenv("DB_PORT"),
		dbUser:       os.Getenv("DB_USER"),
		dbPassword:   os.Getenv("DB_PASSWORD"),
		dbName:       os.Getenv("DB_NAME"),
		dbSslMode:    os.Getenv("DB_SSL_MODE"),
	}
}
