package config

import (
	. "approve/internal/approver/repository"
	. "approve/internal/resolution/repository"
	. "approve/internal/route/repository"
	. "approve/internal/step/repository"
	. "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var (
	dbDriverName = os.Getenv("DB_DRIVER_NAME")
	dbHost       = os.Getenv("DB_HOST")
	dbPort       = os.Getenv("DB_PORT")
	dbUser       = os.Getenv("DB_USER")
	dbPassword   = os.Getenv("DB_PASSWORD")
	dbName       = os.Getenv("DB_NAME")
	dbSslMode    = os.Getenv("DB_SSL_MODE")
)

func BuildServer() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbSslMode,
	)
	db := sqlx.MustConnect(dbDriverName, dsn)
	err := db.Ping()
	if err != nil {
		return err
	}
	routeRepo := NewRouteRepository(db)
	stepGroupeRepo := NewStepGroupRepository(db)
	stepRepo := NewStepRepository(db)
	approverRepo := NewApproverRepository(db)
	resolutionRepo := NewResolutionRepository(db)

	fmt.Printf("%#v\n%#v\n%#v\n%#v\n%#v\n", routeRepo, stepGroupeRepo, stepRepo, approverRepo, resolutionRepo)
	return nil
}
