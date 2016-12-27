package errors

import "github.com/go-errors/errors"

func Wrap(e interface{}) error {
	if e == nil {
		return nil
	}
	return errors.Wrap(e, 1)
}

func WrapS(e interface{}, skip int) error {
	if e == nil {
		return nil
	}
	return errors.Wrap(e, skip+1)
}
