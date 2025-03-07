package service

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	"fmt"
)

type UpdateApprover struct {
	repo ar.ApproverRepository
}

func (svc *UpdateApprover) Execute(request am.UpdateApproverRequest) (routeId int64, err error) {
	if res, err := svc.repo.IsRouteStarted(request.Id); err != nil {
		return 0, err
	} else if res {
		return 0, fmt.Errorf("route was started and cannot be updated")
	}
	return svc.repo.Update(request.ToEntity())
}
