package service

import (
	am "approve/internal/approver/model"
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	"fmt"
)

func UpdateApprover(request am.UpdateApproverRequest) (routeId int64, err error) {
	isRouteStarted, err := approverRepo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return approverRepo.Update(request.ToEntity()) })
}
