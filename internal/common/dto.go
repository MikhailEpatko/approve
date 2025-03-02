package common

type PageRequest struct {
	// Count of items per page
	PageSize int `json:"page_size"     validate:"min=1,max=100"`
	// Number of page
	PageNumber int `json:"page_number" validate:"min=1"`
}

type PageResponse struct {
	// Total count of items
	Total int `json:"total"`
	PageRequest
}

type IdDto struct {
	Id int64 `json:"id"`
}

type Response struct {
	Data interface{} `json:"data"`
	Err  string      `json:"error"`
	PageResponse
}
