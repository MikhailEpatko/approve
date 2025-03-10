package common

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
