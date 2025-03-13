package config

import (
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Begin() (*sqlx.Tx, error)
}

type transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) Transaction {
	return &transaction{db}
}

func (tr *transaction) Begin() (*sqlx.Tx, error) {
	return tr.db.Beginx()
}
