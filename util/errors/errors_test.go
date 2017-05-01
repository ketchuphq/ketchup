package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestError struct{}

func (t *TestError) Error() string {
	return "error"
}

func TestTypeName(t *testing.T) {
	err := &TestError{}
	assert.Equal(t, "*errors.TestError", TypeName(err))

	wrappedErr := Wrap(err)
	assert.Equal(t, "*errors.TestError", TypeName(wrappedErr))
}

func TestIsType(t *testing.T) {
	err := &TestError{}
	wrappedErr := Wrap(err)
	assert.True(t, IsType(err, wrappedErr))
	assert.True(t, IsType(wrappedErr, err))
}
