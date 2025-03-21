package service

import (
	cm "approve/internal/common"
	gm "approve/internal/stepgroup/model"
	"fmt"
)

type UpdateStepGroupRepository interface {
	IsRouteProcessing(stepGroupId int64) (bool, error)
	Update(group gm.StepGroupEntity) (int64, error)
}

type UpdateStepGroup struct {
	repo UpdateStepGroupRepository
}

func (svc *UpdateStepGroup) Execute(request gm.UpdateStepGroupRequest) (groupId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteProcessing(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
