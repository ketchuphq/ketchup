package templates

import (
	"errors"
	"io"

	"github.com/golang/protobuf/proto"

	"github.com/octavore/ketchup/proto/ketchup/models"
	"github.com/octavore/ketchup/server/content/engines"
)

var noneTheme = &models.Theme{
	Name: proto.String("none"),
	Templates: map[string]*models.ThemeTemplate{
		"html": &models.ThemeTemplate{
			Name:   proto.String("html"),
			Engine: proto.String(engines.EngineTypeHTML),
			Data:   proto.String("<html>{{.content}}</html>"),
			Placeholders: []*models.ThemePlaceholder{
				{
					Key: proto.String("content"),
					Type: &models.ThemePlaceholder_Text{
						Text: &models.ContentText{
							Type: models.ContentTextType_html.Enum(),
						},
					},
				},
			},
		},
		"markdown": &models.ThemeTemplate{
			Name:   proto.String("markdown"),
			Engine: proto.String(engines.EngineTypeHTML),
			Data:   proto.String("<html>{{.content}}</html>"),
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
	Assets: map[string]*models.ThemeAsset{},
}

type defaultStore struct{}

func (d *defaultStore) List() ([]*models.Theme, error) {
	return []*models.Theme{noneTheme}, nil
}

func (d *defaultStore) Add(io.Reader) (*models.Theme, error) {
	return nil, errors.New("templates: cannot add to default store")
}

func (d *defaultStore) Get(themeName string) (*models.Theme, error) {
	if themeName == "" || themeName == noneTheme.GetName() {
		return noneTheme, nil
	}
	return nil, nil
}

func (d *defaultStore) GetTemplate(t *models.Theme, template string) (*models.ThemeTemplate, error) {
	if t.GetName() != "" && t.GetName() != noneTheme.GetName() {
		return nil, nil
	}
	return t.Templates["default"], nil
}

func (d *defaultStore) GetAsset(t *models.Theme, asset string) (*models.ThemeAsset, error) {
	return nil, nil
}
