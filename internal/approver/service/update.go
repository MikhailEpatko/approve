package service

import (
	am "approve/internal/approver/model"
	cm "approve/internal/common"
	"fmt"
)

type UpdateApproverRepository interface {
	Update(approver am.ApproverEntity) (int64, error)
	IsRouteStarted(approverId int64) (bool, error)
}

type UpdateApprover struct {
	repo UpdateApproverRepository
}

func (svc *UpdateApprover) Execute(request am.UpdateApproverRequest) (routeId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
