package bigtests

import (
	cm "approve/internal/common"
	conf "approve/internal/config"
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup(repo rr.RouteRepository) rm.RouteEntity {
	route := rm.RouteEntity{
		Name:        "test name",
		Description: "test description",
		Status:      cm.TEMPLATE,
	}
	route.Id, _ = repo.Save(route)
	return route
}

func TestRouteRepository(t *testing.T) {
	a := assert.New(t)
	cfg := conf.NewAppConfig()
	db, err := conf.NewDB(cfg)
	a.Nil(err)
	routeRepo := rr.NewRouteRepository(db)
	deleteRoutes := func() {
		db.MustExec("delete from route")
	}

	t.Run("save route", func(t *testing.T) {
		route := rm.RouteEntity{
			Name:        "test name",
			Description: "test description",
			Status:      cm.TEMPLATE,
		}
		routeId, err := routeRepo.Save(route)

		a.Nil(err)
		a.NotEmpty(routeId)
		deleteRoutes()
	})

	t.Run("get route by id", func(t *testing.T) {
		want := setup(routeRepo)

		got, err := routeRepo.GetById(want.Id)

		a.Nil(err)
		a.NotEmpty(got.Id)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(want.Status, got.Status)
		a.Equal(want.IsApproved, got.IsApproved)
		deleteRoutes()
	})

	t.Run("start route", func(t *testing.T) {
		want := setup(routeRepo)

		tx := db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(routeRepo.StartRoute(tx, want.Id))
		_ = tx.Commit()
		got, err := routeRepo.GetById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(want.IsApproved, got.IsApproved)
		a.Equal(cm.STARTED, got.Status)
		deleteRoutes()
	})

	t.Run("update route", func(t *testing.T) {
		want := setup(routeRepo)
		toUpdate := rm.RouteEntity{
			Id:          want.Id,
			Name:        "new test name",
			Description: "new test description",
			Status:      cm.FINISHED,
			IsApproved:  true,
		}

		_, err := routeRepo.Update(toUpdate)
		a.Nil(err)
		got, err := routeRepo.GetById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(toUpdate.Name, got.Name)
		a.Equal(toUpdate.Description, got.Description)
		a.Equal(want.Status, got.Status)
		a.Equal(want.IsApproved, got.IsApproved)
		deleteRoutes()
	})

	t.Run("is route started (false)", func(t *testing.T) {
		route := setup(routeRepo)

		got, err := routeRepo.IsRouteStarted(route.Id)

		a.Nil(err)
		a.False(got)
		deleteRoutes()
	})

	t.Run("is route started (true)", func(t *testing.T) {
		route := setup(routeRepo)
		tx := db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(routeRepo.StartRoute(tx, route.Id))
		_ = tx.Commit()

		got, err := routeRepo.IsRouteStarted(route.Id)

		a.Nil(err)
		a.True(got)
		deleteRoutes()
	})

	t.Run("finish route (isApproved should be true)", func(t *testing.T) {
		want := setup(routeRepo)
		tx := db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(routeRepo.StartRoute(tx, want.Id))
		_ = tx.Commit()
		isStarted, _ := routeRepo.IsRouteStarted(want.Id)

		a.True(isStarted)

		tx = db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(routeRepo.FinishRoute(tx, want.Id, true))
		_ = tx.Commit()
		got, err := routeRepo.GetById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(cm.FINISHED, got.Status)
		a.True(got.IsApproved)
		deleteRoutes()
	})

	t.Run("finish route (isApproved should be false)", func(t *testing.T) {
		want := setup(routeRepo)
		tx := db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(routeRepo.StartRoute(tx, want.Id))
		_ = tx.Commit()
		isStarted, _ := routeRepo.IsRouteStarted(want.Id)

		a.True(isStarted)

		tx = db.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(routeRepo.FinishRoute(tx, want.Id, false))
		_ = tx.Commit()
		got, err := routeRepo.GetById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(cm.FINISHED, got.Status)
		a.False(got.IsApproved)
		deleteRoutes()
	})
}
