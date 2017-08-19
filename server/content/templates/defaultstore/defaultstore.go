package defaultstore

import (
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/content/engines/html"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
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
			Engine: proto.String(html.EngineTypeHTML),
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
			Engine: proto.String(html.EngineTypeHTML),
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

func (d *DefaultStore) Get(themeName string) (store.Theme, error) {
	if themeName == "" || themeName == noneTheme.GetName() {
		return &Theme{}, nil
	}
	return nil, nil
}

func (d *DefaultStore) GetAsset(assetName string) (*models.ThemeAsset, error) {
	return nil, nil
}
