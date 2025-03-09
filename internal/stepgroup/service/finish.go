package service

import (
	resm "approve/internal/resolution/model"
	"errors"
	"github.com/jmoiron/sqlx"
)

type FinishGroupAndStartNext struct {
}

func (svc *FinishGroupAndStartNext) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isStepApproved bool,
) error {
	return errors.New("not implemented")
}
