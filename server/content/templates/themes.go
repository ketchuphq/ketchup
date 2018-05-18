package templates

import (
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/plugins/pkg"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

// getTheme checks all stores for the theme with the given name.
// returns:
// - theme store the theme was found in
// - theme object
// - ref of the theme
// - (error)
// theme returns the thing
func (m *Module) getTheme(name string) (store.Theme, error) {
	for i := len(m.Stores) - 1; i > -1; i-- {
		store := m.Stores[i]
		theme, err := store.Get(name)
		if err != nil {
			return theme, errors.Wrap(err)
		}
		if theme != nil {
			theme.Proto().Name = &name
			return theme, nil
		}
	}
	return nil, nil
}

// GetTheme returns an installed theme with all related assets populated.
func (m *Module) GetTheme(name string) (*models.Theme, string, error) {
	theme, err := m.getTheme(name)
	if err != nil {
		return nil, "", errors.Wrap(err)
	}
	if theme == nil {
		return nil, "", nil
	}
	ref, _ := theme.Ref()
	return theme.Proto(), ref, nil
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
		asset, err := store.GetAsset(name)
		if err != nil {
			return nil, err
		}
		if asset != nil {
			return asset, nil
		}
	}
	return nil, nil
}

// for mocking
var getLatestRef = pkg.GetLatestRef

// CheckThemeForUpdate checks the given theme for updates,
// and if it has updates, returns true, and the current ref and the latest ref.
func (m *Module) CheckThemeForUpdate(name string) (bool, string, string, error) {
	theme, err := m.getTheme(name)
	if err != nil {
		return false, "", "", err
	}
	if theme == nil {
		return false, "", "", errors.New("theme does not exist")
	}
	ref, ok := theme.Ref()
	if !ok {
		return false, "", "", nil
	}
	vcsURL := theme.Proto().GetPackage().GetVcsUrl()
	if vcsURL == "" {
		return false, ref, "", nil
	}
	remoteRef, err := getLatestRef(vcsURL)
	if err != nil {
		return false, ref, "", err
	}
	if remoteRef == "" {
		return false, ref, "", nil
	}
	return remoteRef != ref, ref, remoteRef, nil
}

type updatable interface {
	UpdateThemeToRef(name, ref string) error
}

func (m *Module) UpdateTheme(name, ref string) error {
	if ts, ok := m.themeStore.(updatable); ok {
		return ts.UpdateThemeToRef(name, ref)
	}
	return errors.New("theme store is not updatable")
}
