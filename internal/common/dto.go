package common

type PageRequest struct {
	PageSize   int `json:"page_size"   validate:"min=1,max=100"`
	PageNumber int `json:"page_number" validate:"min=1"`
}

type PageResponse struct {
	Total int `json:"total"`
	PageRequest
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"error"`
	Data    interface{} `json:"data"`
}
