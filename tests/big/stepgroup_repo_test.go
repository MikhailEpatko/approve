package big

import (
	cm "approve/internal/common"
	cfg "approve/internal/database"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	fx "approve/tests/big/fixtures"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	routeName = "route1"
	status    = cm.STARTED
)

func TestStepGroupRepository(t *testing.T) {
	a := assert.New(t)
	cfg.Connect()
	deleteRoute := func() {
		cfg.DB.MustExec("delete from route")
	}

	defer func() {
		if r := recover(); r != nil {
			cm.Logger.Fatal(fmt.Sprintf("Recovered from panic: %s", r))
			deleteRoute()
		}
	}()

	t.Run("find a step group by route id", func(t *testing.T) {
		route1 := fx.Route(routeName, status)
		group11 := fx.Group(route1, 1, status, cm.SERIAL, false)
		route2 := fx.Route(routeName, status)
		_ = fx.Group(route2, 1, status, cm.PARALLEL_ALL_OF, false)

		got, err := stepGroupRepo.FindByRouteId(route1.Id)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(1, len(got))
		a.Equal(group11, got[0])
		deleteRoute()
	})

	t.Run("stat first group", func(t *testing.T) {
		route1 := fx.Route(routeName, status)
		want := fx.Group(route1, 1, cm.NEW, cm.SERIAL, false)
		group2Before := fx.Group(route1, 2, cm.NEW, cm.PARALLEL_ALL_OF, false)

		tx := cfg.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
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
		route := fx.Route(routeName, cm.NEW)
		groupBefore := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		toUpdate := gm.StepGroupEntity{
			Id:        groupBefore.Id,
			Name:      "new name",
			Number:    2,
			StepOrder: cm.PARALLEL_ALL_OF,
		}

		tx := cfg.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepGroupRepo.Update(tx, toUpdate)
		a.Nil(err)
		a.NotNil(got)
		a.Equal(groupBefore.Id, got)
		a.Nil(tx.Commit())

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
			route := fx.Route(routeName, routeStatus)
			group := fx.Group(route, 1, routeStatus, cm.SERIAL, false)

			tx := cfg.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepGroupRepo.IsRouteProcessing(tx, group.Id)
			a.Nil(tx.Commit())

			a.Nil(err)
			a.Equal(want, got)
			deleteRoute()
		}
	})

	t.Run("finish step group", func(t *testing.T) {
		route := fx.Route(routeName, status)
		group := fx.Group(route, 1, status, cm.SERIAL, false)

		tx := cfg.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(stepGroupRepo.FinishGroup(tx, group.Id))
		a.Nil(tx.Commit())
		got, err := stepGroupRepo.FindById(group.Id)

		a.Nil(err)
		a.Equal(cm.FINISHED, got.Status)
		deleteRoute()
	})

	t.Run("start next step group", func(t *testing.T) {
		route := fx.Route(routeName, status)
		group1 := fx.Group(route, 1, status, cm.SERIAL, false)
		group2 := fx.Group(route, 2, cm.NEW, cm.SERIAL, false)

		tx := cfg.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
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
		route := fx.Route(routeName, status)
		group := fx.Group(route, 1, status, cm.SERIAL, false)
		_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
		lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

		tx := cfg.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
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
		route := fx.Route(routeName, status)
		group := fx.Group(route, 1, status, cm.SERIAL, false)
		_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
		lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

		tx := cfg.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
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

	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ALL_OF and all steps are approved (true)",
		func(t *testing.T) {
			route := fx.Route(routeName, status)
			group := fx.Group(route, 1, status, cm.PARALLEL_ALL_OF, false)
			_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

			tx := cfg.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
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
			route := fx.Route(routeName, status)
			group := fx.Group(route, 1, status, cm.PARALLEL_ALL_OF, false)
			_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

			tx := cfg.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
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

	t.Run(
		"calculate and set IsApproved when step order is PARALLEL_ANY_OF and all steps are approved (true)",
		func(t *testing.T) {
			route := fx.Route(routeName, status)
			group := fx.Group(route, 1, status, cm.PARALLEL_ANY_OF, false)
			_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

			tx := cfg.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
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
			route := fx.Route(routeName, status)
			group := fx.Group(route, 1, status, cm.PARALLEL_ANY_OF, false)
			_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
			lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

			tx := cfg.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
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
			route := fx.Route(routeName, status)
			group := fx.Group(route, 1, status, cm.PARALLEL_ANY_OF, false)
			_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, false)
			lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

			tx := cfg.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
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
