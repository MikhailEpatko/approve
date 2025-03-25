package service

import (
	cm "approve/internal/common"
	"approve/internal/database"
	rm "approve/internal/route/model"
	routeRepo "approve/internal/route/repository"
	"fmt"
)

func UpdateRoute(request rm.UpdateRouteRequest) (routeId int64, err error) {
	err = cm.Validate(request)
	if err != nil {
		return 0, err
	}

	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed updating route: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()

	err = cm.SafeExecute(err, func() error {
		route, innerErr := routeRepo.FindByIdTx(tx, routeId)
		switch {
		case innerErr != nil:
			return innerErr
		case route.Id == 0:
			return ErrRouteNotFound
		case route.Status == cm.FINISHED:
			return ErrRouteIsFinished
		case route.Status == cm.STARTED:
			return ErrRouteAlreadyStarted
		}
		return nil
	})

	return cm.SafeExecuteG(err, func() (int64, error) { return routeRepo.Update(tx, request.ToEntity()) })
}
