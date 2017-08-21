package dummy

import (
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

type DummyTemplateStore struct {
	Themes map[string]*models.Theme
}

func New() *DummyTemplateStore {
	return &DummyTemplateStore{
		Themes: map[string]*models.Theme{},
	}
}

func (d *DummyTemplateStore) List() ([]*models.Theme, error) {
	themes := []*models.Theme{}
	for _, t := range d.Themes {
		themes = append(themes, t)
	}
	return themes, nil
}

func (d *DummyTemplateStore) Add(theme *models.Theme) error {
	d.Themes[theme.GetName()] = theme
	return nil
}

func (d *DummyTemplateStore) AddPackage(p *packages.Package) error {
	panic("not implemented")
}

func (d *DummyTemplateStore) Get(themeName string) (store.Theme, error) {
	theme := d.Themes[themeName]
	if theme == nil {
		return nil, nil
	}
	return &Theme{Theme: theme}, nil
}

func (d *DummyTemplateStore) GetAsset(asset string) (*models.ThemeAsset, error) {
	for _, theme := range d.Themes {
		if asset, ok := theme.GetAssets()[asset]; ok {
			return asset, nil
		}
	}
	return nil, nil
}
