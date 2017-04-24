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
