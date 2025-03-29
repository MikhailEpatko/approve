package service

import (
	"approve/internal/common"
	"approve/internal/database"
	routeModel "approve/internal/route/model"
	routeRepo "approve/internal/route/repository"
	"fmt"
)

func UpdateRoute(request routeModel.UpdateRouteRequest) (routeId int64, err error) {
	err = common.Validate(request)
	if err != nil {
		return 0, err
	}

	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed updating checkedRoute: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()

	checkedRoute, err := common.SafeExecuteG(err, func() (routeModel.RouteEntity, error) {
		route, innerErr := routeRepo.FindByIdTx(tx, routeId)
		switch {
		case innerErr != nil:
			return route, innerErr
		case route.Id == 0:
			return route, ErrRouteNotFound
		case route.Status == common.FINISHED:
			return route, ErrRouteIsFinished
		case route.Status == common.STARTED:
			return route, ErrRouteAlreadyStarted
		}
		return route, nil
	})

	return common.SafeExecuteG(err, func() (int64, error) {
		return routeRepo.Update(tx, request.ToEntity(checkedRoute.Status))
	})
}
