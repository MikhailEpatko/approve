package service

import (
	resm "approve/internal/resolution/model"
	"errors"
	"github.com/jmoiron/sqlx"
)

type ProcessAnyOfStep struct {
}

func (s *ProcessAnyOfStep) Execute(
	*sqlx.Tx,
	resm.ApprovingInfoEntity,
	bool,
) error {
	return errors.New("not implemented")
}
