package model

import "approve/internal/common"

type RouteResponse struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      common.Status `json:"status"`
}
