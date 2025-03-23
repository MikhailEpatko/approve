package common

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate[T any](request T) (err error) {
	err = validate.Struct(request)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return validateErrs
		}
	}
	return err
}
