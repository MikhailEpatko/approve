package service

import (
	cm "approve/internal/common"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"fmt"
)

func UpdateStepGroup(request gm.UpdateStepGroupRequest) (groupId int64, err error) {
	isRouteStarted, err := stepGroupRepo.IsRouteProcessing(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return stepGroupRepo.Update(request.ToEntity()) })
}
