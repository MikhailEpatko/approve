package big

import (
	cm "approve/internal/common"
	"approve/internal/database"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	fx "approve/tests/big/fixtures"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepRepository(t *testing.T) {
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

	t.Run("find a step by step_group_id", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group1 := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step11 := fx.Step(group1, 1, cm.STARTED, cm.SERIAL, false)
		group2 := fx.Group(route, 2, cm.NEW, cm.PARALLEL_ANY_OF, false)
		_ = fx.Step(group2, 2, cm.STARTED, cm.PARALLEL_ANY_OF, false)

		got, err := stepRepo.FindByGroupId(group1.Id)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(1, len(got))
		a.Equal(step11, got[0])
		deleteRoute()
	})

	t.Run("start SERIAL step", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step1 := fx.Step(group, 1, cm.NEW, cm.PARALLEL_ALL_OF, false)
		step2 := fx.Step(group, 2, cm.NEW, cm.PARALLEL_ANY_OF, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.StartSteps(tx, group)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(1, len(got))
		a.Equal(step1.Id, got[0].Id)
		a.Equal(step1.StepGroupId, got[0].StepGroupId)
		a.Equal(step1.Name, got[0].Name)
		a.Equal(step1.Number, got[0].Number)
		a.NotEqual(step1.Status, got[0].Status)
		a.Equal(cm.STARTED, got[0].Status)
		a.Equal(step1.ApproverOrder, got[0].ApproverOrder)
		a.Equal(step1.IsApproved, got[0].IsApproved)
		a.Nil(tx.Commit())

		step2After, err := stepRepo.FindById(step2.Id)
		a.Nil(err)
		a.NotEmpty(step2After)
		a.Equal(step2, step2After)
		deleteRoute()
	})

	t.Run("start PARALLEL_ALL_OF step", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.PARALLEL_ALL_OF, false)
		step1 := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		step2 := fx.Step(group, 2, cm.NEW, cm.SERIAL, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.StartSteps(tx, group)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(2, len(got))
		a.Equal(step1.Id, got[0].Id)
		a.Equal(step1.StepGroupId, got[0].StepGroupId)
		a.Equal(step1.Name, got[0].Name)
		a.Equal(step1.Number, got[0].Number)
		a.NotEqual(step1.Status, got[0].Status)
		a.Equal(cm.STARTED, got[0].Status)
		a.Equal(step1.ApproverOrder, got[0].ApproverOrder)
		a.Equal(step1.IsApproved, got[0].IsApproved)
		a.Equal(step2.Id, got[1].Id)
		a.Equal(step2.StepGroupId, got[1].StepGroupId)
		a.Equal(step2.Name, got[1].Name)
		a.Equal(step2.Number, got[1].Number)
		a.NotEqual(step2.Status, got[1].Status)
		a.Equal(cm.STARTED, got[1].Status)
		a.Equal(step2.ApproverOrder, got[1].ApproverOrder)
		a.Equal(step2.IsApproved, got[1].IsApproved)
		a.Nil(tx.Commit())

		deleteRoute()
	})

	t.Run("start PARALLEL_ANY_OF step", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.PARALLEL_ANY_OF, false)
		step1 := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		step2 := fx.Step(group, 2, cm.NEW, cm.SERIAL, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.StartSteps(tx, group)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(2, len(got))
		a.Equal(step1.Id, got[0].Id)
		a.Equal(step1.StepGroupId, got[0].StepGroupId)
		a.Equal(step1.Name, got[0].Name)
		a.Equal(step1.Number, got[0].Number)
		a.NotEqual(step1.Status, got[0].Status)
		a.Equal(cm.STARTED, got[0].Status)
		a.Equal(step1.ApproverOrder, got[0].ApproverOrder)
		a.Equal(step1.IsApproved, got[0].IsApproved)
		a.Equal(step2.Id, got[1].Id)
		a.Equal(step2.StepGroupId, got[1].StepGroupId)
		a.Equal(step2.Name, got[1].Name)
		a.Equal(step2.Number, got[1].Number)
		a.NotEqual(step2.Status, got[1].Status)
		a.Equal(cm.STARTED, got[1].Status)
		a.Equal(step2.ApproverOrder, got[1].ApproverOrder)
		a.Equal(step2.IsApproved, got[1].IsApproved)
		a.Nil(tx.Commit())

		deleteRoute()
	})

	t.Run("should update step", func(t *testing.T) {
		route := fx.Route("route", cm.NEW)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		toUpdate := sm.StepEntity{
			Id:            step.Id,
			Name:          "new name",
			Number:        2,
			ApproverOrder: cm.PARALLEL_ALL_OF,
		}

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.Update(tx, toUpdate)
		a.Nil(err)
		a.Nil(tx.Commit())

		a.NotEmpty(got)
		a.Equal(step.Id, got)

		stepAfter, err := stepRepo.FindById(step.Id)
		a.Nil(err)
		a.NotEmpty(stepAfter)
		a.Equal(step.Id, stepAfter.Id)
		a.Equal(step.StepGroupId, stepAfter.StepGroupId)
		a.Equal(step.Status, stepAfter.Status)
		a.Equal(step.IsApproved, stepAfter.IsApproved)
		a.Equal(toUpdate.Name, stepAfter.Name)
		a.Equal(toUpdate.Number, stepAfter.Number)
		a.Equal(toUpdate.ApproverOrder, stepAfter.ApproverOrder)
		deleteRoute()
	})

	t.Run("IsRouteProcessing", func(t *testing.T) {
		for status, want := range map[cm.Status]bool{
			cm.STARTED:  true,
			cm.FINISHED: true,
			cm.NEW:      false,
			cm.TEMPLATE: false,
		} {
			route := fx.Route("route", status)
			group := fx.Group(route, 1, status, cm.SERIAL, false)
			step := fx.Step(group, 1, status, cm.SERIAL, false)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.IsRouteProcessing(tx, step.Id)
			a.Nil(err)
			a.Nil(tx.Commit())

			a.Equal(want, got)
			deleteRoute()
		}
	})

	t.Run("FinishStep", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		err := stepRepo.FinishStep(tx, step.Id)
		a.Nil(err)
		a.Nil(tx.Commit())

		stepAfter, err := stepRepo.FindById(step.Id)
		a.Nil(err)
		a.NotEmpty(stepAfter)
		a.Equal(step.Id, stepAfter.Id)
		a.Equal(step.StepGroupId, stepAfter.StepGroupId)
		a.Equal(step.IsApproved, stepAfter.IsApproved)
		a.Equal(step.Name, stepAfter.Name)
		a.Equal(step.Number, stepAfter.Number)
		a.Equal(step.ApproverOrder, stepAfter.ApproverOrder)
		a.NotEqual(step.Status, stepAfter.Status)
		a.Equal(cm.FINISHED, stepAfter.Status)
		deleteRoute()
	})

	t.Run(
		"calculate and set IsApproved when order is SERIAL and all resolutions are approved (true)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, true)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, true)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.True(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when order is SERIAL and not all resolutions are approved (false)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, true)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, false)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.False(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when order is PARALLEL_ALL_OF and all resolutions are approved (true)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.PARALLEL_ALL_OF, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, true)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, true)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.True(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when order is PARALLEL_ALL_OF and not all resolutions are approved (false)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.PARALLEL_ALL_OF, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, true)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, false)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.False(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when order is PARALLEL_ANY_OF and all resolutions are approved (true)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.PARALLEL_ANY_OF, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, true)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, true)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.True(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when order is PARALLEL_ANY_OF and not all resolutions are approved (true)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.PARALLEL_ANY_OF, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, true)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, false)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.True(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.True(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run(
		"calculate and set IsApproved when order is PARALLEL_ANY_OF and no one resolution is approved (false)",
		func(t *testing.T) {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, cm.PARALLEL_ANY_OF, false)
			approver1 := fx.Approver(step, 1, cm.FINISHED)
			_ = fx.Resolution(approver1, false)
			approver2 := fx.Approver(step, 2, cm.FINISHED)
			resolution2 := fx.Resolution(approver2, false)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			got, err := stepRepo.CalculateAndSetIsApproved(
				tx,
				step.Id,
				step.ApproverOrder,
				resolution2.IsApproved,
			)
			a.Nil(err)
			a.False(got)
			a.Nil(tx.Commit())

			stepAfter, err := stepRepo.FindById(step.Id)
			a.Nil(err)
			a.NotEmpty(stepAfter)
			a.False(stepAfter.IsApproved)
			deleteRoute()
		})

	t.Run("ExistsNotFinishedStepsInGroup (true)", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		_ = fx.Step(group, 1, cm.FINISHED, cm.PARALLEL_ANY_OF, false)
		_ = fx.Step(group, 2, cm.NEW, cm.PARALLEL_ANY_OF, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.ExistsNotFinishedStepsInGroup(tx, group.Id)
		a.Nil(err)
		a.True(got)
		a.Nil(tx.Commit())
		deleteRoute()
	})

	t.Run("ExistsNotFinishedStepsInGroup (false)", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		_ = fx.Step(group, 1, cm.FINISHED, cm.PARALLEL_ANY_OF, false)
		_ = fx.Step(group, 2, cm.FINISHED, cm.PARALLEL_ANY_OF, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.ExistsNotFinishedStepsInGroup(tx, group.Id)
		a.Nil(err)
		a.False(got)
		a.Nil(tx.Commit())
		deleteRoute()
	})

	t.Run("StartNextStep should start one step if order is SERIAL", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step1 := fx.Step(group, 1, cm.FINISHED, cm.SERIAL, false)
		step2 := fx.Step(group, 2, cm.NEW, cm.SERIAL, false)
		step3 := fx.Step(group, 3, cm.NEW, cm.SERIAL, false)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := stepRepo.StartNextStep(tx, group.Id, step1.Id)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(step2.Id, got)
		a.Nil(tx.Commit())

		step2After, err := stepRepo.FindById(step2.Id)
		a.Nil(err)
		a.NotEmpty(step2After)
		a.Equal(cm.STARTED, step2After.Status)

		step3After, err := stepRepo.FindById(step3.Id)
		a.Nil(err)
		a.NotEmpty(step3After)
		a.Equal(cm.NEW, step3After.Status)
		deleteRoute()
	})
}
