package dummy

import (
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
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

func (d *DummyTemplateStore) Get(themeName string) (*models.Theme, error) {
	return d.Themes[themeName], nil
}

func (d *DummyTemplateStore) GetTemplate(t *models.Theme, template string) (*models.ThemeTemplate, error) {
	return t.GetTemplates()[template], nil
}

func (d *DummyTemplateStore) GetAsset(t *models.Theme, asset string) (*models.ThemeAsset, error) {
	return t.GetAssets()[asset], nil
}
