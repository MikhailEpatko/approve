package big

import (
	am "approve/internal/approver/model"
	"approve/internal/approver/repository"
	cm "approve/internal/common"
	"approve/internal/database"
	fx "approve/tests/big/fixtures"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApproverRepository(t *testing.T) {
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

	t.Run("should find an approver by step_id", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step1 := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		want := fx.Approver(step1, 2, cm.NEW)
		step2 := fx.Step(group, 2, cm.STARTED, cm.SERIAL, false)
		_ = fx.Approver(step2, 1, cm.NEW)

		got, err := repository.FindByStepId(step1.Id)

		a.Nil(err)
		a.NotNil(got)
		a.Equal(1, len(got))
		a.Equal(want, got[0])
		deleteRoute()
	})

	t.Run("start SERIAL approver by step_id", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver1 := fx.Approver(step, 1, cm.NEW)
		approver2 := fx.Approver(step, 2, cm.NEW)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.StartStepApprovers(tx, step))
		a.Nil(tx.Commit())

		got, err := repository.FindById(approver1.Id)
		a.Nil(err)
		a.Equal(cm.STARTED, got.Status)

		got, err = repository.FindById(approver2.Id)
		a.Nil(err)
		a.Equal(cm.NEW, got.Status)
		deleteRoute()
	})

	t.Run("start PARALLEL_ALL_OF and PARALLEL_ANY_OF approvers by step_id", func(t *testing.T) {
		for _, aprroverOrder := range []cm.OrderType{cm.PARALLEL_ALL_OF, cm.PARALLEL_ANY_OF} {
			route := fx.Route("route", cm.STARTED)
			group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
			step := fx.Step(group, 1, cm.STARTED, aprroverOrder, false)
			approver1 := fx.Approver(step, 1, cm.NEW)
			approver2 := fx.Approver(step, 2, cm.NEW)

			tx := database.DB.MustBegin()
			defer func() { _ = tx.Rollback() }()
			defer func() { _ = tx.Rollback() }()
			a.Nil(repository.StartStepApprovers(tx, step))
			a.Nil(tx.Commit())

			got, err := repository.FindById(approver1.Id)
			a.Nil(err)
			a.Equal(cm.STARTED, got.Status)

			got, err = repository.FindById(approver2.Id)
			a.Nil(err)
			a.Equal(cm.STARTED, got.Status)
			deleteRoute()
		}
	})

	t.Run("should update approver", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.NEW)
		toUpdate := am.ApproverEntity{
			Id:       approver.Id,
			StepId:   12345,
			Guid:     "new_guid",
			Name:     "new_name",
			Position: "new_position",
			Email:    "new_email@email.com",
			Number:   3,
			Status:   cm.STARTED,
		}

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := repository.Update(tx, toUpdate)
		a.Nil(err)
		a.Nil(tx.Commit())
		a.NotNil(got)
		a.Equal(approver.Id, got)

		approverAfter, err := repository.FindById(approver.Id)
		a.Nil(err)
		a.NotEmpty(approverAfter)
		a.Equal(approver.Id, approverAfter.Id)
		a.Equal(approver.StepId, approverAfter.StepId)
		a.Equal(approver.Status, approverAfter.Status)
		a.Equal(toUpdate.Guid, approverAfter.Guid)
		a.Equal(toUpdate.Name, approverAfter.Name)
		a.Equal(toUpdate.Position, approverAfter.Position)
		a.Equal(toUpdate.Email, approverAfter.Email)
		a.Equal(toUpdate.Number, approverAfter.Number)
		deleteRoute()
	})

	t.Run("should finish an approver by id", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.STARTED)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.FinishApprover(tx, approver.Id))
		a.Nil(tx.Commit())

		got, err := repository.FindById(approver.Id)
		a.Nil(err)
		a.NotNil(got)
		a.Equal(cm.FINISHED, got.Status)
		deleteRoute()
	})

	t.Run("ExistNotFinishedApproversInStep (true)", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		_ = fx.Approver(step, 1, cm.FINISHED)
		_ = fx.Approver(step, 2, cm.STARTED)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		res, err := repository.ExistNotFinishedApproversInStep(tx, step.Id)
		a.Nil(err)
		a.True(res)
		a.Nil(tx.Commit())
		deleteRoute()
	})

	t.Run("ExistNotFinishedApproversInStep (false)", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		_ = fx.Approver(step, 1, cm.FINISHED)
		_ = fx.Approver(step, 2, cm.FINISHED)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		res, err := repository.ExistNotFinishedApproversInStep(tx, step.Id)
		a.Nil(err)
		a.False(res)
		a.Nil(tx.Commit())
		deleteRoute()
	})

	t.Run("should start next approver in step", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver1 := fx.Approver(step, 1, cm.FINISHED)
		approver2 := fx.Approver(step, 2, cm.NEW)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		err := repository.StartNextApprover(tx, step.Id, approver1.Id)
		a.Nil(err)
		a.Nil(tx.Commit())

		approver2After, err := repository.FindById(approver2.Id)
		a.Nil(err)
		a.NotEmpty(approver2After)
		a.Equal(cm.STARTED, approver2After.Status)
		deleteRoute()
	})

	t.Run("IsRouteProcessing (STARTED - true)", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.NEW)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		res, err := repository.IsRouteProcessing(tx, approver.Id)
		a.Nil(err)
		a.Nil(tx.Commit())
		a.True(res)
		deleteRoute()
	})

	t.Run("IsRouteProcessing (FINISHED - true)", func(t *testing.T) {
		route := fx.Route("route", cm.FINISHED)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.NEW)

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		res, err := repository.IsRouteProcessing(tx, approver.Id)
		a.Nil(err)
		a.Nil(tx.Commit())
		a.True(res)
		deleteRoute()
	})

	t.Run("IsRouteProcessing (NEW - false)", func(t *testing.T) {
		route := fx.Route("route", cm.NEW)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.NEW)

		tx := database.DB.MustBegin()
		res, err := repository.IsRouteProcessing(tx, approver.Id)
		a.Nil(err)
		a.Nil(tx.Commit())
		a.False(res)
		deleteRoute()
	})

	t.Run("SaveAll sholuld save and return list of approvers", func(t *testing.T) {
		route := fx.Route("route", cm.NEW)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		approver1 := am.ApproverEntity{
			StepId:   step.Id,
			Guid:     "guid1",
			Name:     "name1",
			Email:    "email1",
			Position: "position1",
			Number:   1,
			Status:   cm.NEW,
		}
		approver2 := am.ApproverEntity{
			StepId:   step.Id,
			Guid:     "guid2",
			Name:     "name2",
			Email:    "email2",
			Position: "position2",
			Number:   2,
			Status:   cm.NEW,
		}
		approver3 := am.ApproverEntity{
			StepId:   step.Id,
			Guid:     "guid3",
			Name:     "name3",
			Email:    "email3",
			Position: "position3",
			Number:   3,
			Status:   cm.NEW,
		}
		approversToSave := []am.ApproverEntity{approver1, approver2, approver3}

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		err := repository.SaveAll(tx, approversToSave)
		a.Nil(err)
		a.Nil(tx.Commit())

		savedApprovers, err := repository.FindByStepId(step.Id)

		a.Nil(err)
		a.NotNil(savedApprovers)
		a.Equal(3, len(savedApprovers))

		idx1 := slices.IndexFunc(savedApprovers, func(a am.ApproverEntity) bool { return a.Number == 1 })
		idx2 := slices.IndexFunc(savedApprovers, func(a am.ApproverEntity) bool { return a.Number == 2 })
		idx3 := slices.IndexFunc(savedApprovers, func(a am.ApproverEntity) bool { return a.Number == 3 })
		saved1 := savedApprovers[idx1]
		saved2 := savedApprovers[idx2]
		saved3 := savedApprovers[idx3]

		a.NotEmpty(saved1)
		a.NotEmpty(saved1.Id)
		a.Equal(approver1.StepId, saved1.StepId)
		a.Equal(approver1.Status, saved1.Status)
		a.Equal(approver1.Guid, saved1.Guid)
		a.Equal(approver1.Name, saved1.Name)
		a.Equal(approver1.Position, saved1.Position)
		a.Equal(approver1.Email, saved1.Email)
		a.Equal(approver1.Number, saved1.Number)

		a.NotEmpty(saved2)
		a.NotEmpty(saved2.Id)
		a.Equal(approver2.StepId, saved2.StepId)
		a.Equal(approver2.Status, saved2.Status)
		a.Equal(approver2.Guid, saved2.Guid)
		a.Equal(approver2.Name, saved2.Name)
		a.Equal(approver2.Position, saved2.Position)
		a.Equal(approver2.Email, saved2.Email)
		a.Equal(approver2.Number, saved2.Number)

		a.NotEmpty(saved3)
		a.NotEmpty(saved3.Id)
		a.Equal(approver3.StepId, saved3.StepId)
		a.Equal(approver3.Status, saved3.Status)
		a.Equal(approver3.Guid, saved3.Guid)
		a.Equal(approver3.Name, saved3.Name)
		a.Equal(approver3.Position, saved3.Position)
		a.Equal(approver3.Email, saved3.Email)
		a.Equal(approver3.Number, saved3.Number)

		deleteRoute()
	})
}
