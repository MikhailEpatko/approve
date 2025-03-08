package service

import (
	resm "approve/internal/resolution/model"
	"errors"
	"github.com/jmoiron/sqlx"
)

type ProcessParallelAllOf struct {
}

func (*ProcessParallelAllOf) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	approved bool,
) error {
	return errors.New("not implemented")
}
