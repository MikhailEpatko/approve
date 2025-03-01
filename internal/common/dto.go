package common

type PageRequest struct {
	// Count of items per page
	PageSize int `json:"page_size" validate:"positive"`
	// Number of page
	PageNumber int `json:"page_number" validate:"positive"`
}

type PageResponse struct {
	// Total count of items
	Total int `json:"total"`
	PageRequest
}

type Id struct {
	Id int64 `json:"id"`
}
