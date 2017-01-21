package errors

import (
	"fmt"
	"reflect"

	"github.com/go-errors/errors"
)

// New returns a new wrapped error with m as message
func New(m string, a ...interface{}) error {
	return errors.Wrap(fmt.Errorf(m, a...), 1)
}

// Wrap an error e if it is not nil
func Wrap(e error) error {
	if e == nil {
		return nil
	}
	return errors.Wrap(e, 1)
}

// WrapS is like Wrap but skips 'skip' lines of trace
func WrapS(e error, skip int) error {
	if e == nil {
		return nil
	}
	return errors.Wrap(e, skip+1)
}

func IsType(e1, e2 error) bool {
	if err, ok := e1.(*errors.Error); ok {
		return IsType(err.Err, e2)
	}
	if err, ok := e2.(*errors.Error); ok {
		return IsType(e1, err.Err)
	}
	return reflect.TypeOf(e1) == reflect.TypeOf(e2)
}

func TypeName(e error) string {
	if err, ok := e.(*errors.Error); ok {
		return err.TypeName()
	}
	return fmt.Sprintf("%T", e)
}
