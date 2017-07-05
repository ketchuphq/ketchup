package templates

import (
	"time"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/ketchuphq/ketchup/plugins/pkg"
	"github.com/ketchuphq/ketchup/server/config"
	"github.com/ketchuphq/ketchup/server/content/templates/defaultstore"
	"github.com/ketchuphq/ketchup/server/content/templates/filestore"
)

const (
	themeDir           = "themes"
	internalThemeDir   = "internal_themes"
	defaultRegistryURL = "http://themes.ketchuphq.com/registry.json"
	devRegistryURL     = "http://localhost:8000/registry.json"
)

type ThemesConfig struct {
	Themes struct {
		Path        string `json:"path"`
		RegistryURL string `json:"registry_url"`
	} `json:"themes"`
}

// Module templates provides support for looking up and using themes and
// their corresponding templates.
type Module struct {
	ConfigModule *config.Module
	Logger       *logger.Module
	Pkg          *pkg.Module

	themeRegistry    *pkg.Registry
	themeRegistryURL string
	themeStore       *filestore.FileStore
	internalStore    *filestore.FileStore
	Stores           []ThemeStore

	config ThemesConfig
}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		err := m.ConfigModule.ReadConfig(&m.config)
		if err != nil {
			return err
		}

		m.config.Themes.Path = m.ConfigModule.DataPath(m.config.Themes.Path, themeDir)
		m.themeStore, err = filestore.New(m.config.Themes.Path, time.Second*10, m.Logger.Error)
		if err != nil {
			return err
		}

		m.internalStore, err = filestore.New(
			m.ConfigModule.DataPath(internalThemeDir, ""),
			time.Second*10,
			m.Logger.Error,
		)
		if err != nil {
			return err
		}

		m.Stores = []ThemeStore{
			&defaultstore.DefaultStore{},
			m.themeStore,
			m.internalStore,
		}
		registryURL := defaultRegistryURL
		if c.Env().IsDevelopment() {
			registryURL = devRegistryURL
		}
		m.themeRegistryURL = registryURL
		m.themeRegistry = m.Pkg.Registry(registryURL)
		err = m.themeRegistry.Sync()
		if err != nil {
			m.Logger.Warning(err)
		}
		return nil
	}
}

// GetRegistryURL returns the configured registry URL
func (m *Module) GetRegistryURL() string {
	return m.themeRegistryURL
}
