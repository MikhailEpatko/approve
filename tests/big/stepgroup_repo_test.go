package big

import (
	cm "approve/internal/common"
	"approve/internal/database"
	gm "approve/internal/stepgroup/model"
	"approve/internal/stepgroup/repository"
	fx "approve/tests/big/fixtures"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

var (
	routeName = "route1"
	status    = cm.STARTED
)

func TestStepGroupRepository(t *testing.T) {
	a := assert.New(t)
	database.Connect()
	deleteRoute := func() {
		database.DB.MustExec("delete from route")
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

		got, err := repository.FindByRouteId(route1.Id)

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

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := repository.StartFirstGroup(tx, route1.Id)
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

		group2After, err := repository.FindById(group2Before.Id)

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

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := repository.Update(tx, toUpdate)
		a.Nil(err)
		a.NotNil(got)
		a.Equal(groupBefore.Id, got)
		a.Nil(tx.Commit())

		groupAfter, err := repository.FindById(groupBefore.Id)

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

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := repository.IsRouteProcessing(tx, group.Id)
			a.Nil(tx.Commit())

			a.Nil(err)
			a.Equal(want, got)
			deleteRoute()
		}
	})

	t.Run("finish step group", func(t *testing.T) {
		route := fx.Route(routeName, status)
		group := fx.Group(route, 1, status, cm.SERIAL, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.FinishGroup(tx, group.Id))
		a.Nil(tx.Commit())
		got, err := repository.FindById(group.Id)

		a.Nil(err)
		a.Equal(cm.FINISHED, got.Status)
		deleteRoute()
	})

	t.Run("start next step group", func(t *testing.T) {
		route := fx.Route(routeName, status)
		group1 := fx.Group(route, 1, status, cm.SERIAL, false)
		group2 := fx.Group(route, 2, cm.NEW, cm.SERIAL, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		nextGroupId, err := repository.StartNextGroup(tx, route.Id, group1.Id)
		a.Nil(err)
		a.Equal(nextGroupId, group2.Id)
		a.Nil(tx.Commit())

		got, err := repository.FindById(group2.Id)

		a.Nil(err)
		a.Equal(cm.STARTED, got.Status)
		deleteRoute()
	})

	t.Run("calculate and set IsApproved when step order is SERIAL and all steps are approved (true)", func(t *testing.T) {
		route := fx.Route(routeName, status)
		group := fx.Group(route, 1, status, cm.SERIAL, false)
		_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
		lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, true)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := repository.CalculateAndSetIsApproved(
			tx,
			group.Id,
			group.StepOrder,
			lastStep.IsApproved,
		)
		a.Nil(err)
		a.True(got)
		a.Nil(tx.Commit())

		groupAfter, err := repository.FindById(group.Id)
		a.Nil(err)
		a.True(groupAfter.IsApproved)
		deleteRoute()
	})

	t.Run("calculate and set IsApproved when step order is SERIAL and not all steps are approved (false)", func(t *testing.T) {
		route := fx.Route(routeName, status)
		group := fx.Group(route, 1, status, cm.SERIAL, false)
		_ = fx.Step(group, 1, cm.FINISHED, cm.SERIAL, true)
		lastStep := fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ALL_OF, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := repository.CalculateAndSetIsApproved(
			tx,
			group.Id,
			group.StepOrder,
			lastStep.IsApproved,
		)
		a.Nil(err)
		a.False(got)
		a.Nil(tx.Commit())

		groupAfter, err := repository.FindById(group.Id)
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

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := repository.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			groupAfter, err := repository.FindById(group.Id)
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

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := repository.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			groupAfter, err := repository.FindById(group.Id)
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

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := repository.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			groupAfter, err := repository.FindById(group.Id)
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

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := repository.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			groupAfter, err := repository.FindById(group.Id)
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

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := repository.CalculateAndSetIsApproved(
				tx,
				group.Id,
				group.StepOrder,
				lastStep.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			groupAfter, err := repository.FindById(group.Id)
			a.Nil(err)
			a.False(groupAfter.IsApproved)
			deleteRoute()
		})

	t.Run("SaveAll should save all step groups and return saved", func(t *testing.T) {
		route := fx.Route("route", cm.NEW)
		group1 := gm.StepGroupEntity{
			RouteId:   route.Id,
			Name:      "name1",
			Number:    1,
			StepOrder: cm.PARALLEL_ANY_OF,
			Status:    cm.NEW,
		}
		group2 := gm.StepGroupEntity{
			RouteId:   route.Id,
			Name:      "name2",
			Number:    2,
			StepOrder: cm.SERIAL,
			Status:    cm.NEW,
		}

		toSave := []gm.StepGroupEntity{group1, group2}

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		saved, err := repository.SaveAll(tx, toSave)
		a.Nil(err)
		a.Nil(tx.Commit())
		a.NotEmpty(saved)
		a.Equal(2, len(saved))

		idx1 := slices.IndexFunc(toSave, func(s gm.StepGroupEntity) bool { return s.Number == 1 })
		idx2 := slices.IndexFunc(toSave, func(s gm.StepGroupEntity) bool { return s.Number == 2 })
		saved1 := saved[idx1]
		saved2 := saved[idx2]

		a.Equal(group1.RouteId, saved1.RouteId)
		a.Equal(group1.Number, saved1.Number)
		a.Equal(group1.Name, saved1.Name)
		a.Equal(group1.Status, saved1.Status)
		a.Equal(group1.StepOrder, saved1.StepOrder)

		a.Equal(group2.RouteId, saved2.RouteId)
		a.Equal(group2.Number, saved2.Number)
		a.Equal(group2.Name, saved2.Name)
		a.Equal(group2.Status, saved2.Status)
		a.Equal(group2.StepOrder, saved2.StepOrder)

		deleteRoute()
	})
}
