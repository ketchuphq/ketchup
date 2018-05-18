package filestore

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/engines/html"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

// todo: test key not the same as Name
var testTheme = &models.Theme{
	Name: proto.String("test-theme"),
	Templates: map[string]*models.ThemeTemplate{
		"test-template.foo": {
			Name:   proto.String("test-template.foo"),
			Engine: proto.String(html.EngineTypeHTML),
			Data:   proto.String(`<div>{{.Page.Content}}</div>`),
			Placeholders: []*models.ThemePlaceholder{
				{
					Key: proto.String("content"),
					Type: &models.ThemePlaceholder_Text{
						Text: &models.ContentText{
							Type: models.ContentTextType_markdown.Enum(),
						},
					},
				},
			},
		},
		"test-no-engine.foo": {
			Name: proto.String("test-no-engine.foo"),
			Data: proto.String(`<div>{{.Page.Content}}</div>`),
			Placeholders: []*models.ThemePlaceholder{
				{
					Key: proto.String("content"),
					Type: &models.ThemePlaceholder_Text{
						Text: &models.ContentText{
							Type: models.ContentTextType_markdown.Enum(),
						},
					},
				},
			},
		},
		"test-empty-engine.foo": {
			Name:   proto.String("test-empty-engine.foo"),
			Engine: proto.String(""),
			Data:   proto.String(`<div>{{.Page.Content}}</div>`),
			Placeholders: []*models.ThemePlaceholder{
				{
					Key: proto.String("content"),
					Type: &models.ThemePlaceholder_Text{
						Text: &models.ContentText{
							Type: models.ContentTextType_markdown.Enum(),
						},
					},
				},
			},
		},
	},
	Assets: map[string]*models.ThemeAsset{
		"app.js": {
			Name: proto.String("app.js"),
			Data: proto.String("var foo = 1;"),
		},
	},
}

func newForTest(t *testing.T) *FileStore {
	dir, err := ioutil.TempDir("", "ketchup-filestore-test-")
	assert.Nil(t, err)

	fs, err := New(dir, time.Hour, log.Println)
	assert.Nil(t, err)
	return fs
}

func TestGet(t *testing.T) {
	fs := newForTest(t)
	err := fs.Add(testTheme)
	assert.Nil(t, err)

	expected := store.StripData(testTheme)
	expected.Templates["test-empty-engine.foo"].Engine = proto.String("foo")
	theme, err := fs.Get(expected.GetName())
	assert.Nil(t, err)
	assert.Equal(t, expected, theme.Proto())
}

func TestAddAndList(t *testing.T) {
	fs := newForTest(t)

	themes, err := fs.List()
	assert.Nil(t, err)
	assert.Empty(t, themes)

	err = fs.Add(testTheme)
	assert.Nil(t, err)

	themes, err = fs.List()
	assert.Nil(t, err)
	assert.Equal(t, []*models.Theme{store.StripData(testTheme)}, themes)
}

func TestGetAsset(t *testing.T) {
	fs := newForTest(t)
	err := fs.Add(testTheme)
	assert.Nil(t, err)

	asset, err := fs.GetAsset("fake")
	assert.Nil(t, err)
	assert.Nil(t, asset)

	expectedAsset := testTheme.Assets["app.js"]
	// expectedAsset.SetTheme(testTheme.Name)

	asset, err = fs.GetAsset("app.js")
	assert.Nil(t, err)
	assert.Equal(t, expectedAsset, asset)
}

func TestUpdateThemeDirMap(t *testing.T) {
	fs := newForTest(t)
	err := fs.Add(testTheme)
	assert.Nil(t, err)
	assert.Empty(t, fs.themeDirMap)

	altName := "alt-theme-name"
	oldPath := path.Join(fs.baseDir, testTheme.GetName())
	newPath := path.Join(fs.baseDir, altName)
	err = os.Rename(oldPath, newPath)
	assert.Nil(t, err)

	// check that rename worked
	expected := map[string]string{
		testTheme.GetName(): altName,
	}
	err = fs.updateThemeDirMap()
	assert.Nil(t, err)
	assert.Equal(t, expected, fs.themeDirMap)

	// check that rename back works
	err = os.Rename(newPath, oldPath)
	assert.Nil(t, err)

	err = fs.updateThemeDirMap()
	assert.Nil(t, err)
	assert.Empty(t, fs.themeDirMap)
}
