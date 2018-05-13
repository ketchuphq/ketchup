package filestore

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/stretchr/testify/assert"
)

func TestThemeGetTemplate(t *testing.T) {
	fs := newForTest(t)
	err := fs.Add(testTheme)
	assert.Empty(t, err)

	theme, err := fs.Get(testTheme.GetName())
	assert.Empty(t, err)

	// get non existent template
	template, err := theme.GetTemplate("fake")
	assert.Nil(t, err)
	assert.Nil(t, template)

	// get actual template
	expectedTemplate := proto.Clone(testTheme.Templates["test-template.foo"]).(*models.ThemeTemplate)
	expectedTemplate.SetTheme(testTheme.Name)

	template, err = theme.GetTemplate("test-template.foo")
	assert.Nil(t, err)
	assert.Equal(t, expectedTemplate, template)

	// test setting engine from proto default
	template, err = theme.GetTemplate("test-no-engine.foo")
	assert.Nil(t, err)
	assert.Equal(t, "html", template.GetEngine())
}

func TestThemeGetAsset(t *testing.T) {
	fs := newForTest(t)
	err := fs.Add(testTheme)
	assert.Empty(t, err)

	theme, err := fs.Get(testTheme.GetName())
	assert.Empty(t, err)

	// get non existent asset
	asset, err := theme.GetAsset("fake")
	assert.Nil(t, err)
	assert.Nil(t, asset)

	// get actual asset
	expectedAsset := proto.Clone(testTheme.Assets["app.js"]).(*models.ThemeAsset)
	expectedAsset.SetTheme(testTheme.Name)

	asset, err = theme.GetAsset("app.js")
	assert.Nil(t, err)
	assert.Equal(t, expectedAsset, asset)
}

func TestThemeRef(t *testing.T) {
	theme := &Theme{}
	ref, valid := theme.Ref()
	assert.Equal(t, "", ref)
	assert.Equal(t, false, valid)

	expected := "e3076fc3ef90c41d42d7927f4302a7fd5b95a8a8"
	theme.ref = expected
	ref, valid = theme.Ref()
	assert.Equal(t, expected, ref)
	assert.Equal(t, true, valid)
}
