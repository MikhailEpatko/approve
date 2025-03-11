package service

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	"fmt"
)

type UpdateApprover struct {
	repo ar.ApproverRepository
}

func (svc *UpdateApprover) Execute(request am.UpdateApproverRequest) (routeId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
