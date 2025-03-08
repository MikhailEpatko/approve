package service

import (
	resm "approve/internal/resolution/model"
	"errors"
	"github.com/jmoiron/sqlx"
)

type FinishStepAndStartNext struct {
}

func (svc *FinishStepAndStartNext) Execute(*sqlx.Tx, resm.ApprovingInfoEntity) error {
	return errors.New("not implemented")
}
