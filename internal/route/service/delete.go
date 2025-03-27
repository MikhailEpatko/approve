package service

import routeRepo "approve/internal/route/repository"

func DeleteRouteById(routeId int64) (err error) {
	return routeRepo.DeleteById(routeId)
}
