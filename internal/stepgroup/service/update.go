package service

import (
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
)

type UpdateStepGroup struct {
	repo gr.StepGroupRepository
}

func (svc *UpdateStepGroup) Execute(request gm.UpdateStepGroupRequest) (groupId int64, err error) {
	if res, err := svc.repo.IsRouteStarted(request.Id); err != nil {
		return 0, err
	} else if res {
		return 0, fmt.Errorf("route was started and cannot be updated")
	}
	return svc.repo.Update(request.ToEntity())
}
