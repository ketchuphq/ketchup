package templates

import (
	"github.com/octavore/ketchup/proto/ketchup/models"
	"github.com/octavore/ketchup/util/errors"
)

func (m *Module) getTheme(name string) (ThemeStore, *models.Theme, error) {
	for i := len(m.Stores) - 1; i > -1; i-- {
		store := m.Stores[i]
		theme, err := store.Get(name)
		if err != nil {
			return nil, nil, errors.Wrap(err)
		}
		if theme != nil {
			theme.Name = &name
			return store, theme, nil
		}
	}
	return nil, nil, nil
}

// GetTheme returns a theme with all related assets populated.
func (m *Module) GetTheme(name string) (*models.Theme, error) {
	_, theme, err := m.getTheme(name)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return theme, nil
}

// ListThemes returns a list of all known themes.
func (m *Module) ListThemes() ([]*models.Theme, error) {
	themes := []*models.Theme{}
	for _, store := range m.Stores {
		t, err := store.List()
		if err != nil {
			return nil, err
		}
		themes = append(themes, t...)
	}
	return themes, nil
}

// GetAsset searches all themes for the named asset
func (m *Module) GetAsset(name string) (*models.ThemeAsset, error) {
	for i := len(m.Stores) - 1; i != 0; i-- {
		store := m.Stores[i]
		themes, err := store.List()
		if err != nil {
			return nil, err
		}
		for _, theme := range themes {
			asset, err := store.GetAsset(theme, name)
			if err != nil {
				return nil, err
			}
			if asset != nil {
				return asset, nil
			}
		}
	}
	return nil, nil
}
