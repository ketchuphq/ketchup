package templates

import "github.com/ketchuphq/ketchup/proto/ketchup/packages"

func (m *Module) Registry() (*packages.Registry, error) {
	err := m.themeRegistry.Sync()
	if err != nil {
		return nil, err
	}
	return m.themeRegistry.Proto(), nil
}

func (m *Module) SearchRegistry(themeName string) (*packages.Package, error) {
	return m.themeRegistry.Search(themeName)
}

// InstallThemeFromPackage to the default store from a remote git.
// todo: check for existing package
func (m *Module) InstallThemeFromPackage(p *packages.Package) error {
	return m.themeStore.AddPackage(p)
}
