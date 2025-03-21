package big

import (
	cm "approve/internal/common"
	conf "approve/internal/config"
	resr "approve/internal/resolution/repository"
	"approve/tests/big/fixtures"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	resolutionRepo *resr.ResolutionRepository
)

func TestResolutionRepository(t *testing.T) {
	a := assert.New(t)
	cfg := conf.NewAppConfig()
	db, err := conf.NewDB(cfg)
	a.Nil(err)
	fx := fixtures.New(db)
	resolutionRepo = resr.NewResolutionRepository(db)
	deleteRoute := func() {
		db.MustExec("delete from route")
	}

	defer func() {
		if r := recover(); r != nil {
			cm.Logger.Fatal(fmt.Sprintf("Recovered from panic: %s", r))
			deleteRoute()
		}
	}()

	t.Run("should find resolution by approver_id", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.SERIAL, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.FINISHED)
		want := fx.Resolution(approver, true)
		approver2 := fx.Approver(step, 2, cm.STARTED)
		_ = fx.Resolution(approver2, true)

		got, err := resolutionRepo.FindByApproverId(approver.Id)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(1, len(got))
		a.Equal(want, got[0])
		deleteRoute()
	})

	t.Run("should find approving info", func(t *testing.T) {
		route := fx.Route("route", cm.STARTED)
		group := fx.Group(route, 1, cm.STARTED, cm.PARALLEL_ANY_OF, false)
		step := fx.Step(group, 1, cm.STARTED, cm.SERIAL, false)
		approver := fx.Approver(step, 1, cm.STARTED)
		group2 := fx.Group(route, 2, cm.NEW, cm.SERIAL, false)
		step2 := fx.Step(group2, 1, cm.NEW, cm.PARALLEL_ALL_OF, false)
		_ = fx.Approver(step2, 1, cm.NEW)

		tx := db.MustBegin()
		got, err := resolutionRepo.ApprovingInfoTx(tx, approver.Id)
		a.Nil(tx.Commit())

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(route.Id, got.RouteId)
		a.Equal(group.Id, got.StepGroupId)
		a.Equal(step.Id, got.StepId)
		a.Equal(approver.Id, got.ApproverId)
		a.Equal(cm.PARALLEL_ANY_OF, got.StepOrder)
		a.Equal(cm.SERIAL, got.ApproverOrder)
		deleteRoute()
	})
}
