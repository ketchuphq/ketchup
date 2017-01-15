package errors

import (
	"fmt"

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
