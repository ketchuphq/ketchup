package defaultstore

import (
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/content/engines"
)

var noneTemplate = `<html>
	<style>#c{max-width:600px;margin:0 auto;font-family:helvetica, sans-serif;}</style>
	<div id='c'>{{.Page.Content}}</div>
</html>`

var noneTheme = &models.Theme{
	Name: proto.String("none"),
	Templates: map[string]*models.ThemeTemplate{
		"html": {
			Name:   proto.String("html"),
			Engine: proto.String(engines.EngineTypeHTML),
			Data:   &noneTemplate,
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
		"markdown": {
			Name:   proto.String("markdown"),
			Engine: proto.String(engines.EngineTypeHTML),
			Data:   &noneTemplate,
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

type DefaultStore struct{}

func (d *DefaultStore) List() ([]*models.Theme, error) {
	return []*models.Theme{noneTheme}, nil
}

func (d *DefaultStore) Add(*models.Theme) error {
	return errors.New("templates: cannot add to default store")
}

func (d *DefaultStore) AddPackage(p *packages.Package) error {
	return errors.New("templates: cannot add to default store")
}

func (d *DefaultStore) Get(themeName string) (*models.Theme, string, error) {
	if themeName == "" || themeName == noneTheme.GetName() {
		return noneTheme, "", nil
	}
	return nil, "", nil
}

func (d *DefaultStore) GetTemplate(t *models.Theme, template string) (*models.ThemeTemplate, error) {
	if t.GetName() != "" && t.GetName() != noneTheme.GetName() {
		return nil, nil
	}
	return t.Templates[template], nil
}

func (d *DefaultStore) GetAsset(t *models.Theme, asset string) (*models.ThemeAsset, error) {
	return nil, nil
}
