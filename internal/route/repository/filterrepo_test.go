package repository

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	conf "approve/internal/config"
	rm "approve/internal/route/model"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	name1        = "route1"
	description1 = "test route description 1"
	status1      = cm.TEMPLATE
	name2        = "route2"
	description2 = "test route description 2"
	status2      = cm.STARTED
	guid         = "guid"
	routes       RouteRepository
	stepGroups   gr.StepGroupRepository
	steps        sr.StepRepository
	approvers    ar.ApproverRepository
	log          = cm.Logger
)

func TestFindByFilterRouteRepository(t *testing.T) {
	a := assert.New(t)
	cfg := conf.NewAppConfig()
	db, err := conf.NewDB(cfg)
	a.Nil(err)
	cleaner := cm.NewCleaner(db)
	routes = NewRouteRepository(db)
	stepGroups = gr.NewStepGroupRepository(db)
	steps = sr.NewStepRepository(db)
	approvers = ar.NewApproverRepository(db)
	findByFilterRepo := NewFindByFilterRouteRepository(db)

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(fmt.Sprintf("Recovered from panic: %s", r))
			cleaner.ClearDb()
		}
	}()

	t.Run("find a route by status TEMPLATE", func(t *testing.T) {
		createRoutes()
		filter := rm.FilterRouteRequest{
			Guid:   "",
			Status: status1,
			Text:   "",
			PageRequest: cm.PageRequest{
				PageNumber: 1,
				PageSize:   2,
			},
		}

		got, total, err := findByFilterRepo.FindByfilter(filter)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(int64(1), total)
		a.Equal(1, len(got))
		a.Equal(status1, got[0].Status)
		a.Equal(name1, got[0].Name)
		a.Equal(description1, got[0].Description)
		a.False(got[0].IsApproved)
		cleaner.ClearDb()
	})

	t.Run("find a route by status STARTED", func(t *testing.T) {
		createRoutes()
		filter := rm.FilterRouteRequest{
			Guid:   "",
			Status: status2,
			Text:   "",
			PageRequest: cm.PageRequest{
				PageNumber: 1,
				PageSize:   2,
			},
		}

		got, total, err := findByFilterRepo.FindByfilter(filter)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(int64(1), total)
		a.Equal(1, len(got))
		a.Equal(status2, got[0].Status)
		a.Equal(name2, got[0].Name)
		a.Equal(description2, got[0].Description)
		a.False(got[0].IsApproved)
		cleaner.ClearDb()
	})

	t.Run("there is no any route with status FINISHED", func(t *testing.T) {
		createRoutes()
		filter := rm.FilterRouteRequest{
			Guid:   "",
			Status: cm.FINISHED,
			Text:   "",
			PageRequest: cm.PageRequest{
				PageNumber: 1,
				PageSize:   2,
			},
		}

		got, total, err := findByFilterRepo.FindByfilter(filter)

		a.Nil(err)
		a.Nil(got)
		a.Equal(int64(0), total)
		a.Equal(0, len(got))
		cleaner.ClearDb()
	})

	t.Run("find a route by name 'ute1'", func(t *testing.T) {
		createRoutes()
		filter := rm.FilterRouteRequest{
			Guid:   "",
			Status: "",
			Text:   "ute1",
			PageRequest: cm.PageRequest{
				PageNumber: 1,
				PageSize:   2,
			},
		}

		got, total, err := findByFilterRepo.FindByfilter(filter)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(int64(1), total)
		a.Equal(1, len(got))
		a.Equal(status1, got[0].Status)
		a.Equal(name1, got[0].Name)
		a.Equal(description1, got[0].Description)
		a.False(got[0].IsApproved)
		cleaner.ClearDb()
	})

	t.Run("find a route by description", func(t *testing.T) {
		createRoutes()
		filter := rm.FilterRouteRequest{
			Guid:   "",
			Status: "",
			Text:   "description 2",
			PageRequest: cm.PageRequest{
				PageNumber: 1,
				PageSize:   2,
			},
		}

		got, total, err := findByFilterRepo.FindByfilter(filter)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(int64(1), total)
		a.Equal(1, len(got))
		a.Equal(status2, got[0].Status)
		a.Equal(name2, got[0].Name)
		a.Equal(description2, got[0].Description)
		a.False(got[0].IsApproved)
		cleaner.ClearDb()
	})

	t.Run("find a route by guid", func(t *testing.T) {
		createRoutes()
		filter := rm.FilterRouteRequest{
			Guid:   guid + "-1111",
			Status: "",
			Text:   "",
			PageRequest: cm.PageRequest{
				PageNumber: 1,
				PageSize:   2,
			},
		}

		got, total, err := findByFilterRepo.FindByfilter(filter)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(int64(1), total)
		a.Equal(1, len(got))
		a.Equal(status1, got[0].Status)
		a.Equal(name1, got[0].Name)
		a.Equal(description1, got[0].Description)
		a.False(got[0].IsApproved)
		cleaner.ClearDb()
	})

}

func createRoutes() {
	route1 := rm.RouteEntity{
		Name:        name1,
		Description: description1,
		Status:      status1,
	}
	route1.Id, _ = routes.Save(route1)
	group11 := gm.StepGroupEntity{
		RouteId:   route1.Id,
		Name:      name1 + "-group11",
		Number:    1,
		Status:    status1,
		StepOrder: cm.SERIAL,
	}
	group11.Id, _ = stepGroups.Save(group11)
	step111 := sm.StepEntity{
		StepGroupId:   group11.Id,
		Name:          group11.Name + "-step111",
		Number:        1,
		Status:        status1,
		ApproverOrder: cm.SERIAL,
	}
	step111.Id, _ = steps.Save(step111)
	approver1111 := am.ApproverEntity{
		StepId:   step111.Id,
		Guid:     guid + "-1111",
		Name:     step111.Name + "-approver1111",
		Position: "position-1111",
		Email:    "email-1111",
		Number:   1,
		Status:   status1,
	}
	approver1111.Id, _ = approvers.Save(approver1111)
	approver1112 := am.ApproverEntity{
		StepId:   step111.Id,
		Guid:     guid + "-1111",
		Name:     step111.Name + "-approver1112",
		Position: "position-1112",
		Email:    "email-1112",
		Number:   2,
		Status:   status1,
	}
	approver1112.Id, _ = approvers.Save(approver1112)

	group12 := gm.StepGroupEntity{
		RouteId:   route1.Id,
		Name:      name1 + "group12",
		Number:    2,
		Status:    status1,
		StepOrder: cm.PARALLEL_ALL_OF,
	}
	group12.Id, _ = stepGroups.Save(group12)
	step121 := sm.StepEntity{
		StepGroupId:   group12.Id,
		Name:          group12.Name + "-step121",
		Number:        1,
		Status:        status1,
		ApproverOrder: cm.PARALLEL_ALL_OF,
	}
	step121.Id, _ = steps.Save(step121)
	approver1211 := am.ApproverEntity{
		StepId:   step121.Id,
		Guid:     guid + "-1211",
		Name:     step121.Name + "-approver1211",
		Position: "position-1211",
		Email:    "email-1211",
		Number:   1,
		Status:   status1,
	}
	approver1211.Id, _ = approvers.Save(approver1211)

	route2 := rm.RouteEntity{
		Name:        name2,
		Description: description2,
		Status:      status2,
	}
	route2.Id, _ = routes.Save(route2)

	group21 := gm.StepGroupEntity{
		RouteId:   route2.Id,
		Name:      name2 + "-group21",
		Number:    1,
		Status:    status2,
		StepOrder: cm.SERIAL,
	}
	group21.Id, _ = stepGroups.Save(group21)
	step211 := sm.StepEntity{
		StepGroupId:   group21.Id,
		Name:          group21.Name + "-step211",
		Number:        1,
		Status:        status2,
		ApproverOrder: cm.SERIAL,
	}
	step211.Id, _ = steps.Save(step211)
	approver2111 := am.ApproverEntity{
		StepId:   step211.Id,
		Guid:     guid + "-2111",
		Name:     step211.Name + "-approver2111",
		Position: "position-2111",
		Email:    "email-2111",
		Number:   1,
		Status:   status2,
	}
	approver2111.Id, _ = approvers.Save(approver2111)

	group22 := gm.StepGroupEntity{
		RouteId:   route2.Id,
		Name:      name2 + "-group22",
		Number:    2,
		Status:    status2,
		StepOrder: cm.PARALLEL_ALL_OF,
	}
	group22.Id, _ = stepGroups.Save(group22)
	step221 := sm.StepEntity{
		StepGroupId:   group22.Id,
		Name:          group22.Name + "-step221",
		Number:        1,
		Status:        status2,
		ApproverOrder: cm.PARALLEL_ALL_OF,
	}
	step221.Id, _ = steps.Save(step211)
	approver2211 := am.ApproverEntity{
		StepId:   step221.Id,
		Guid:     guid + "-2211",
		Name:     step211.Name + "-approver2211",
		Position: "position-2211",
		Email:    "email-2211",
		Number:   1,
		Status:   status2,
	}
	approver2211.Id, _ = approvers.Save(approver2211)
}
