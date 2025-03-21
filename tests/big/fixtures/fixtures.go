package fixtures

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	resr "approve/internal/resolution/repository"
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Fixtures struct {
	routeRepo      *rr.RouteRepository
	stepGroupRepo  *gr.StepGroupRepository
	stepRepo       *sr.StepRepository
	approverRepo   *ar.ApproverRepository
	resolutionRepo *resr.ResolutionRepository
}

func New(db *sqlx.DB) Fixtures {
	return Fixtures{
		routeRepo:      rr.NewRouteRepository(db),
		stepGroupRepo:  gr.NewStepGroupRepository(db),
		stepRepo:       sr.NewStepRepository(db),
		approverRepo:   ar.NewApproverRepository(db),
		resolutionRepo: resr.NewResolutionRepository(db),
	}
}

func (fx *Fixtures) Route(
	name string,
	routeStatus cm.Status,
) rm.RouteEntity {
	route := rm.RouteEntity{
		Name:        name,
		Description: "description",
		Status:      routeStatus,
		IsApproved:  false,
	}
	id, err := fx.routeRepo.Save(route)
	if err != nil {
		panic(err)
	}
	route.Id = id
	return route
}

func (fx *Fixtures) Group(
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
	id, err := fx.stepGroupRepo.Save(group)
	if err != nil {
		panic(err)
	}
	group.Id = id
	return group
}

func (fx *Fixtures) Step(
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
	id, err := fx.stepRepo.Save(step)
	if err != nil {
		panic(err)
	}
	step.Id = id
	return step
}

func (fx *Fixtures) Approver(
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
	id, err := fx.approverRepo.Save(approver)
	if err != nil {
		panic(err)
	}
	approver.Id = id
	return approver
}

func (fx *Fixtures) Resolution(
	approver am.ApproverEntity,
	isApproved bool,
) resm.ResolutionEntity {
	resolution := resm.ResolutionEntity{
		ApproverId: approver.Id,
		IsApproved: isApproved,
		Comment:    "comment",
	}
	id, err := fx.resolutionRepo.Save(resolution)
	if err != nil {
		panic(err)
	}
	resolution.Id = id
	return resolution
}
