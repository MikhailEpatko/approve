package service

import (
	resm "approve/internal/resolution/model"
	"errors"
	"github.com/jmoiron/sqlx"
)

type ProcessParallelAnyOf struct {
}

func (s *ProcessParallelAnyOf) Execute(
	*sqlx.Tx,
	resm.ApprovingInfoEntity,
	bool,
) error {
	return errors.New("not implemented")
}
