package templates

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplate(t *testing.T) {
	m, stop := setup(false, testTheme)
	defer stop()

	tmpl, err := m.GetTemplate("fake-theme", "fake-template")
	assert.Nil(t, tmpl)
	assert.EqualError(t, err, `content: theme "fake-theme" not found`)

	tmpl, err = m.GetTemplate(testTheme.GetName(), "fake-template")
	assert.Nil(t, tmpl)
	assert.EqualError(t, err, `content: template "fake-template" not found for theme "test-theme"`)

	expected := proto.Clone(testTemplate).(*models.ThemeTemplate)
	expected.SetTheme(testTheme.Name)
	tmpl, err = m.GetTemplate(testTheme.GetName(), "test-template")
	assert.NotNil(t, tmpl)
	assert.Equal(t, expected, tmpl)
}
