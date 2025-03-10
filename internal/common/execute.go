package common

import "fmt"

func SafeExecute(err error, f func() error) error {
	if err != nil {
		return err
	}
	return f()
}

func SafeExecuteBool(err error, f func() (bool, error)) (bool, error) {
	if err != nil {
		return false, err
	}
	return f()
}

func ErrorOrNil(msg string, err error) error {
	if err != nil {
		err = fmt.Errorf("%s: %w", msg, err)
	}
	return err
}

func SafeExecuteInt64(err error, f func() (int64, error)) (int64, error) {
	if err != nil {
		return 0, err
	}
	return f()
}
