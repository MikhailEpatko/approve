package common

import "fmt"

func SafeExecute(err error, f func() error) error {
	if err != nil {
		return err
	}
	return f()
}

func SafeExecuteG[T any](err error, f func() (T, error)) (T, error) {
	if err != nil {
		return *new(T), err
	}
	return f()
}

func ErrorOrNil(msg string, err error) error {
	if err != nil {
		err = fmt.Errorf("%s: %w", msg, err)
	}
	return err
}
