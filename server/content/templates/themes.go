package templates

import (
	"github.com/ketchuphq/ketchup/plugins/pkg"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/util/errors"
)

// getTheme checks all stores for the theme with the given name.
// returns:
// - theme store the theme was found in
// - theme object
// - ref of the theme
// - (error)
func (m *Module) getTheme(name string) (ThemeStore, *models.Theme, string, error) {
	for i := len(m.Stores) - 1; i > -1; i-- {
		store := m.Stores[i]
		theme, ref, err := store.Get(name)
		if err != nil {
			return nil, nil, "", errors.Wrap(err)
		}
		if theme != nil {
			theme.Name = &name
			return store, theme, ref, nil
		}
	}
	return nil, nil, "", nil
}

// GetTheme returns an installed theme with all related assets populated.
func (m *Module) GetTheme(name string) (*models.Theme, string, error) {
	_, theme, ref, err := m.getTheme(name)
	if err != nil {
		return nil, "", errors.Wrap(err)
	}
	return theme, ref, nil
}

// ListThemes returns a list of all installed themes.
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

// GetAsset searches all installed themes for the named asset
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

// CheckThemeForUpdate checks the given theme for updates,
// and if true, returns the current ref and the latest ref.
func (m *Module) CheckThemeForUpdate(name string) (bool, string, string, error) {
	_, theme, ref, err := m.getTheme(name)
	if err != nil {
		return false, "", "", err
	}
	if ref == "" {
		return false, "", "", nil
	}
	vcsURL := theme.GetPackage().GetVcsUrl()
	if vcsURL == "" {
		return false, ref, "", nil
	}
	remoteRef, err := pkg.GetLatestRef(vcsURL)
	if err != nil {
		return false, ref, "", err
	}
	if remoteRef == "" {
		return false, ref, "", nil
	}
	return remoteRef == ref, ref, remoteRef, nil
}

func (m *Module) UpdateTheme(name, ref string) error {
	return m.themeStore.UpdateThemeToRef(name, ref)
}
