package service

import (
	cm "approve/internal/common"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
)

type UpdateStepGroup struct {
	repo gr.StepGroupRepository
}

func (svc *UpdateStepGroup) Execute(request gm.UpdateStepGroupRequest) (groupId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
