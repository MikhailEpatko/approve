package big

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	conf "approve/internal/config"
	"approve/tests/big/fixtures"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	approverRepo *ar.ApproverRepository
)

func TestApproverRepository(t *testing.T) {
	a := assert.New(t)
	cfg := conf.NewAppConfig()
	db, err := conf.NewDB(cfg)
	a.Nil(err)
	fx := fixtures.New(db)
	approverRepo = ar.NewApproverRepository(db)
	deleteRoute := func() {
		db.MustExec("delete from route")
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

		got, err := approverRepo.FindByStepId(step1.Id)

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

		tx := db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(approverRepo.StartStepApprovers(tx, step))
		a.Nil(tx.Commit())

		got, err := approverRepo.FindById(approver1.Id)
		a.Nil(err)
		a.Equal(cm.STARTED, got.Status)

		got, err = approverRepo.FindById(approver2.Id)
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

			tx := db.MustBegin()
			defer func() { _ = tx.Rollback() }()
			a.Nil(approverRepo.StartStepApprovers(tx, step))
			a.Nil(tx.Commit())

			got, err := approverRepo.FindById(approver1.Id)
			a.Nil(err)
			a.Equal(cm.STARTED, got.Status)

			got, err = approverRepo.FindById(approver2.Id)
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

		got, err := approverRepo.Update(toUpdate)
		a.Nil(err)
		a.NotNil(got)
		a.Equal(approver.Id, got)

		approverAfter, err := approverRepo.FindById(approver.Id)
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

		tx := db.MustBegin()
		a.Nil(approverRepo.FinishApprover(tx, approver.Id))
		a.Nil(tx.Commit())

		got, err := approverRepo.FindById(approver.Id)
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

		tx := db.MustBegin()
		res, err := approverRepo.ExistNotFinishedApproversInStep(tx, step.Id)
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

		tx := db.MustBegin()
		res, err := approverRepo.ExistNotFinishedApproversInStep(tx, step.Id)
		a.Nil(err)
		a.False(res)
		a.Nil(tx.Commit())
		deleteRoute()
	})

	t.Run("should start nex approver in step", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver1 := fx.Approver(step, 1, cm.FINISHED)
		approver2 := fx.Approver(step, 2, cm.NEW)

		tx := db.MustBegin()
		err := approverRepo.StartNextApprover(tx, step.Id, approver1.Id)
		a.Nil(err)
		a.Nil(tx.Commit())

		approver2After, err := approverRepo.FindById(approver2.Id)
		a.Nil(err)
		a.NotEmpty(approver2After)
		a.Equal(cm.STARTED, approver2After.Status)
		deleteRoute()
	})

	t.Run("IsRouteStarted (true)", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.NEW)

		res, err := approverRepo.IsRouteStarted(approver.Id)
		a.Nil(err)
		a.True(res)
		deleteRoute()
	})

	t.Run("IsRouteStarted (false)", func(t *testing.T) {
		route := fx.Route("route", cm.NEW)
		group := fx.Group(route, 1, cm.NEW, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.NEW, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.NEW)

		res, err := approverRepo.IsRouteStarted(approver.Id)
		a.Nil(err)
		a.False(res)
		deleteRoute()
	})

}
