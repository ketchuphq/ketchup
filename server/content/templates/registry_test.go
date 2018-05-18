package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrySync(t *testing.T) {
	m, stop := setup(false)
	defer stop()
	reg, err := m.Registry()
	assert.Nil(t, err)
	assert.Empty(t, reg.Packages)
}

func TestRegistrySearch(t *testing.T) {
	m, stop := setup(false)
	defer stop()
	theme, err := m.SearchRegistry("fake-theme")
	assert.Nil(t, err)
	assert.Nil(t, theme)
}
