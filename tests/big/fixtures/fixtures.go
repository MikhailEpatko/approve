package fixtures

import (
	am "approve/internal/approver/model"
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	resolutionRepo "approve/internal/resolution/repository"
	rm "approve/internal/route/model"
	routeRepo "approve/internal/route/repository"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"fmt"
)

func Route(
	name string,
	routeStatus cm.Status,
) rm.RouteEntity {
	route := rm.RouteEntity{
		Name:        name,
		Description: "description",
		Status:      routeStatus,
		IsApproved:  false,
	}
	id, err := routeRepo.Save(route)
	if err != nil {
		panic(err)
	}
	route.Id = id
	return route
}

func Group(
	route rm.RouteEntity,
	number int,
	groupStatus cm.Status,
	stepOrder cm.OrderType,
	isApproved bool,
) gm.StepGroupEntity {
	group := gm.StepGroupEntity{
		RouteId:    route.Id,
		Name:       fmt.Sprintf("%s-group-%d", route.Name, number),
		Number:     number,
		Status:     groupStatus,
		StepOrder:  stepOrder,
		IsApproved: isApproved,
	}
	id, err := stepGroupRepo.Save(group)
	if err != nil {
		panic(err)
	}
	group.Id = id
	return group
}

func Step(
	group gm.StepGroupEntity,
	number int,
	orderStatus cm.Status,
	approverOrder cm.OrderType,
	isApproved bool,
) sm.StepEntity {
	step := sm.StepEntity{
		StepGroupId:   group.Id,
		Name:          fmt.Sprintf("%s-step-%d", group.Name, number),
		Number:        number,
		Status:        orderStatus,
		ApproverOrder: approverOrder,
		IsApproved:    isApproved,
	}
	id, err := stepRepo.Save(step)
	if err != nil {
		panic(err)
	}
	step.Id = id
	return step
}

func Approver(
	step sm.StepEntity,
	number int,
	status cm.Status,
) am.ApproverEntity {
	approver := am.ApproverEntity{
		StepId:   step.Id,
		Guid:     fmt.Sprintf("guid-%d", number),
		Name:     fmt.Sprintf("%s-approver-%d", step.Name, number),
		Position: fmt.Sprintf("position-%d", number),
		Email:    fmt.Sprintf("email-%d@mail.ru", number),
		Number:   number,
		Status:   status,
	}
	id, err := approverRepo.Save(approver)
	if err != nil {
		panic(err)
	}
	approver.Id = id
	return approver
}

func Resolution(
	approver am.ApproverEntity,
	isApproved bool,
) resm.ResolutionEntity {
	resolution := resm.ResolutionEntity{
		ApproverId: approver.Id,
		IsApproved: isApproved,
		Comment:    "comment",
	}
	id, err := resolutionRepo.Save(resolution)
	if err != nil {
		panic(err)
	}
	resolution.Id = id
	return resolution
}
