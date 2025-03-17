package fixtures

import (
	cm "approve/internal/common"
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
	routeRepo     rr.RouteRepository
	stepGroupRepo gr.StepGroupRepository
	stepRepo      sr.StepRepository
}

func New(db *sqlx.DB) Fixtures {
	return Fixtures{
		routeRepo:     rr.NewRouteRepository(db),
		stepGroupRepo: gr.NewStepGroupRepository(db),
		stepRepo:      sr.NewStepRepository(db),
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
) gm.StepGroupEntity {
	group := gm.StepGroupEntity{
		RouteId:   route.Id,
		Name:      fmt.Sprintf("%s-group-%d", route.Name, number),
		Number:    number,
		Status:    groupStatus,
		StepOrder: stepOrder,
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
