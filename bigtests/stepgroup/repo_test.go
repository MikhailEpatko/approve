package route

import (
	cm "approve/internal/common"
	conf "approve/internal/config"
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	log           = cm.Logger
	routeName     = "route1"
	description   = "test route description 1"
	status        = cm.STARTED
	routeRepo     rr.RouteRepository
	stepGroupRepo gr.StepGroupRepository
	stepRepo      sr.StepRepository
)

func TestFindByFilterRouteRepository(t *testing.T) {
	a := assert.New(t)
	cfg := conf.NewAppConfig()
	db, err := conf.NewDB(cfg)
	a.Nil(err)
	routeRepo = rr.NewRouteRepository(db)
	stepGroupRepo = gr.NewStepGroupRepository(db)
	stepRepo = sr.NewStepRepository(db)
	deleteRoute := func() {
		db.MustExec("delete from route")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(fmt.Sprintf("Recovered from panic: %s", r))
			deleteRoute()
		}
	}()

	t.Run("find a step group by route id", func(t *testing.T) {
		route1 := route(routeName, status)
		group11 := group(route1, 1, status, cm.SERIAL)
		route2 := route(routeName, status)
		_ = group(route2, 1, status, cm.PARALLEL_ALL_OF)

		got, err := stepGroupRepo.FindByRouteId(route1.Id)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(1, len(got))
		a.Equal(group11, got[0])
		deleteRoute()
	})

	t.Run("stat first group", func(t *testing.T) {
		route1 := route(routeName, status)
		want := group(route1, 1, cm.NEW, cm.SERIAL)
		group2Before := group(route1, 2, cm.NEW, cm.PARALLEL_ALL_OF)

		tx := db.MustBegin()
		got, err := stepGroupRepo.StartFirstGroup(tx, route1.Id)
		a.Nil(tx.Commit())

		a.Nil(err)
		a.NotNil(got)
		a.Equal(want.Id, got.Id)
		a.Equal(want.RouteId, got.RouteId)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Number, got.Number)
		a.NotEqual(want.Status, got.Status)
		a.Equal(cm.STARTED, got.Status)
		a.Equal(want.StepOrder, got.StepOrder)
		a.Equal(want.IsApproved, got.IsApproved)

		group2After, err := stepGroupRepo.FindById(group2Before.Id)

		a.Nil(err)
		a.Equal(group2Before, group2After)
		deleteRoute()
	})

	t.Run("should update group", func(t *testing.T) {
		route := route(routeName, cm.NEW)
		groupBefore := group(route, 1, cm.NEW, cm.SERIAL)
		toUpdate := gm.StepGroupEntity{
			Id:        groupBefore.Id,
			Name:      "new name",
			Number:    2,
			StepOrder: cm.PARALLEL_ALL_OF,
		}

		got, err := stepGroupRepo.Update(toUpdate)
		a.Nil(err)
		a.NotNil(got)
		a.Equal(groupBefore.Id, got)

		groupAfter, err := stepGroupRepo.FindById(groupBefore.Id)

		a.Nil(err)
		a.NotEmpty(groupAfter)
		a.Equal(groupBefore.Id, groupAfter.Id)
		a.Equal(groupBefore.RouteId, groupAfter.RouteId)
		a.Equal(groupBefore.Status, groupAfter.Status)
		a.Equal(groupBefore.IsApproved, groupAfter.IsApproved)
		a.Equal(toUpdate.Name, groupAfter.Name)
		a.Equal(toUpdate.Number, groupAfter.Number)
		a.Equal(toUpdate.StepOrder, groupAfter.StepOrder)
		deleteRoute()
	})

	t.Run("IsRouteProcessing", func(t *testing.T) {
		for routeStatus, want := range map[cm.Status]bool{
			cm.TEMPLATE: false,
			cm.NEW:      false,
			cm.STARTED:  true,
			cm.FINISHED: true,
		} {
			route := route(routeName, routeStatus)
			group := group(route, 1, routeStatus, cm.SERIAL)

			got, err := stepGroupRepo.IsRouteProcessing(group.Id)

			a.Nil(err)
			a.Equal(want, got)
			deleteRoute()
		}
	})

	t.Run("finish step group", func(t *testing.T) {
		route := route(routeName, status)
		group := group(route, 1, status, cm.SERIAL)

		tx := db.MustBegin()
		a.Nil(stepGroupRepo.FinishGroup(tx, group.Id))
		a.Nil(tx.Commit())
		got, err := stepGroupRepo.FindById(group.Id)

		a.Nil(err)
		a.Equal(cm.FINISHED, got.Status)
		deleteRoute()
	})

	t.Run("start next step group", func(t *testing.T) {
		route := route(routeName, status)
		group1 := group(route, 1, status, cm.SERIAL)
		group2 := group(route, 2, cm.NEW, cm.SERIAL)

		tx := db.MustBegin()
		nextGroupId, err := stepGroupRepo.StartNextGroup(tx, route.Id, group1.Id)
		a.Nil(err)
		a.Equal(nextGroupId, group2.Id)
		a.Nil(tx.Commit())

		got, err := stepGroupRepo.FindById(group2.Id)

		a.Nil(err)
		a.Equal(cm.STARTED, got.Status)
		deleteRoute()
	})

	t.Run("calculate and set IsApproved when step order is SERIAL and all steps are approved (true)", func(t *testing.T) {
		route := route(routeName, status)
		group := group(route, 1, status, cm.SERIAL)
		_ = step(group, 1, cm.FINISHED, cm.SERIAL, true)
		lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

		tx := db.MustBegin()
		got, err := stepGroupRepo.CalculateAndSetIsApproved(
			tx,
			group.Id,
			group.StepOrder,
			lastStep.IsApproved,
		)
		a.Nil(err)
		a.True(got)
		a.Nil(tx.Commit())

		groupAfter, err := stepGroupRepo.FindById(group.Id)
		a.Nil(err)
		a.True(groupAfter.IsApproved)
		deleteRoute()
	})

	t.Run("calculate and set IsApproved when step order is SERIAL and not all steps are approved (false)", func(t *testing.T) {
		route := route(routeName, status)
		group := group(route, 1, status, cm.SERIAL)
		_ = step(group, 1, cm.FINISHED, cm.SERIAL, true)
		lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

		tx := db.MustBegin()
		got, err := stepGroupRepo.CalculateAndSetIsApproved(
			tx,
			group.Id,
			group.StepOrder,
			lastStep.IsApproved,
		)
		a.Nil(err)
		a.False(got)
		a.Nil(tx.Commit())

		groupAfter, err := stepGroupRepo.FindById(group.Id)
		a.Nil(err)
		a.False(groupAfter.IsApproved)
		deleteRoute()
	})
	//=================
	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ALL_OF and all steps are approved (true)",
		func(t *testing.T) {
			route := route(routeName, status)
			group := group(route, 1, status, cm.PARALLEL_ALL_OF)
			_ = step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

			tx := db.MustBegin()
			got, err := stepGroupRepo.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			groupAfter, err := stepGroupRepo.FindById(group.Id)
			a.Nil(err)
			a.True(groupAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ALL_OF and not all steps are approved (false)",
		func(t *testing.T) {
			route := route(routeName, status)
			group := group(route, 1, status, cm.PARALLEL_ALL_OF)
			_ = step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

			tx := db.MustBegin()
			got, err := stepGroupRepo.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			groupAfter, err := stepGroupRepo.FindById(group.Id)
			a.Nil(err)
			a.False(groupAfter.IsApproved)
			deleteRoute()
		})
	//=================
	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ANY_OF and all steps are approved (true)",
		func(t *testing.T) {
			route := route(routeName, status)
			group := group(route, 1, status, cm.PARALLEL_ANY_OF)
			_ = step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

			tx := db.MustBegin()
			got, err := stepGroupRepo.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			groupAfter, err := stepGroupRepo.FindById(group.Id)
			a.Nil(err)
			a.True(groupAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ANY_OF and one of steps is approved (true)",
		func(t *testing.T) {
			route := route(routeName, status)
			group := group(route, 1, status, cm.PARALLEL_ANY_OF)
			_ = step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

			tx := db.MustBegin()
			got, err := stepGroupRepo.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			groupAfter, err := stepGroupRepo.FindById(group.Id)
			a.Nil(err)
			a.True(groupAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ANY_OF and no one step is approved (false)",
		func(t *testing.T) {
			route := route(routeName, status)
			group := group(route, 1, status, cm.PARALLEL_ANY_OF)
			_ = step(group, 1, cm.FINISHED, cm.SERIAL, false)
			lastStep := step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

			tx := db.MustBegin()
			got, err := stepGroupRepo.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			groupAfter, err := stepGroupRepo.FindById(group.Id)
			a.Nil(err)
			a.False(groupAfter.IsApproved)
			deleteRoute()
		})

}

func route(
	name string,
	routeStatus cm.Status,
) rm.RouteEntity {
	route := rm.RouteEntity{
		Name:        name,
		Description: description,
		Status:      routeStatus,
	}
	id, err := routeRepo.Save(route)
	if err != nil {
		panic(err)
	}
	route.Id = id
	return route
}

func group(
	route rm.RouteEntity,
	number int,
	groupStatus cm.Status,
	stepOrder cm.OrderType,
) gm.StepGroupEntity {
	group := gm.StepGroupEntity{
		RouteId:   route.Id,
		Name:      fmt.Sprintf("%s-group-%d", routeName, number),
		Number:    number,
		Status:    groupStatus,
		StepOrder: stepOrder,
	}
	id, err := stepGroupRepo.Save(group)
	if err != nil {
		panic(err)
	}
	group.Id = id
	return group
}

func step(
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
