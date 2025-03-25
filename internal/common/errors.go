package common

type RequestValidationError struct {
	Message string
}

func (err RequestValidationError) Error() string {
	return err.Message
}
