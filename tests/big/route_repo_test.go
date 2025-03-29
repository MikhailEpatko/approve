package big

import (
	cm "approve/internal/common"
	"approve/internal/database"
	rm "approve/internal/route/model"
	"approve/internal/route/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() rm.RouteEntity {
	route := rm.RouteEntity{
		Name:        "test name",
		Description: "test description",
		Status:      cm.TEMPLATE,
	}
	route.Id, _ = repository.Save(route)
	return route
}

func TestRouteRepository(t *testing.T) {
	a := assert.New(t)
	database.Connect()
	deleteRoutes := func() {
		database.DB.MustExec("delete from route")
	}

	t.Run("save route", func(t *testing.T) {
		route := rm.RouteEntity{
			Name:        "test name",
			Description: "test description",
			Status:      cm.TEMPLATE,
		}
		routeId, err := repository.Save(route)

		a.Nil(err)
		a.NotEmpty(routeId)
		deleteRoutes()
	})

	t.Run("get route by id", func(t *testing.T) {
		want := setup()

		got, err := repository.FindById(want.Id)

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
		want := setup()

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.StartRoute(tx, want.Id))
		_ = tx.Commit()
		got, err := repository.FindById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(want.IsApproved, got.IsApproved)
		a.Equal(cm.STARTED, got.Status)
		deleteRoutes()
	})

	t.Run("update route", func(t *testing.T) {
		want := setup()
		toUpdate := rm.RouteEntity{
			Id:          want.Id,
			Name:        "new test name",
			Description: "new test description",
			Status:      cm.FINISHED,
			IsApproved:  true,
		}

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		_, err := repository.Update(tx, toUpdate)
		a.Nil(err)
		a.Nil(tx.Commit())
		got, err := repository.FindById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(toUpdate.Name, got.Name)
		a.Equal(toUpdate.Description, got.Description)
		a.Equal(want.Status, got.Status)
		a.Equal(want.IsApproved, got.IsApproved)
		deleteRoutes()
	})

	t.Run("is route started (false)", func(t *testing.T) {
		route := setup()

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		got, err := repository.IsRouteStarted(tx, route.Id)
		a.Nil(err)
		a.Nil(tx.Commit())

		a.Nil(err)
		a.False(got)
		deleteRoutes()
	})

	t.Run("is route started (true)", func(t *testing.T) {
		route := setup()
		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.StartRoute(tx, route.Id))

		got, err := repository.IsRouteStarted(tx, route.Id)
		a.Nil(err)
		a.Nil(tx.Commit())

		a.Nil(err)
		a.True(got)
		deleteRoutes()
	})

	t.Run("finish route (isApproved should be true)", func(t *testing.T) {
		want := setup()
		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.StartRoute(tx, want.Id))
		isStarted, _ := repository.IsRouteStarted(tx, want.Id)
		a.Nil(tx.Commit())

		a.True(isStarted)

		tx = database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.FinishRoute(tx, want.Id, true))
		_ = tx.Commit()
		got, err := repository.FindById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(cm.FINISHED, got.Status)
		a.True(got.IsApproved)
		deleteRoutes()
	})

	t.Run("finish route (isApproved should be false)", func(t *testing.T) {
		want := setup()
		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.StartRoute(tx, want.Id))
		isStarted, _ := repository.IsRouteStarted(tx, want.Id)
		a.Nil(tx.Commit())

		a.True(isStarted)

		tx = database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		a.Nil(repository.FinishRoute(tx, want.Id, false))
		_ = tx.Commit()
		got, err := repository.FindById(want.Id)

		a.Nil(err)
		a.Equal(want.Id, got.Id)
		a.Equal(want.Name, got.Name)
		a.Equal(want.Description, got.Description)
		a.Equal(cm.FINISHED, got.Status)
		a.False(got.IsApproved)
		deleteRoutes()
	})

	t.Run("SaveTx should save route and return its id)", func(t *testing.T) {
		toSave := rm.RouteEntity{
			Name:        "name",
			Description: "description",
			Status:      cm.NEW,
		}

		tx := database.DB.MustBegin()
		defer func() { _ = tx.Rollback() }()
		id, err := repository.SaveTx(tx, toSave)
		a.Nil(err)
		a.Nil(tx.Commit())

		got, err := repository.FindById(id)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(toSave.Name, got.Name)
		a.Equal(toSave.Description, got.Description)
		a.Equal(toSave.Status, got.Status)
		a.Equal(toSave.IsApproved, got.IsApproved)
		deleteRoutes()
	})
}
