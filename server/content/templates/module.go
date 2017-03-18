package templates

import (
	"time"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/octavore/ketchup/server/config"
	"github.com/octavore/ketchup/server/content/templates/defaultstore"
	"github.com/octavore/ketchup/server/content/templates/filestore"
)

const (
	themeDir         = "themes"
	internalThemeDir = "internal_themes"
)

type ThemesConfig struct {
	Themes struct {
		Path        string `json:"dir"`
		RegistryURL string `json:"registry_url"`
	} `json:"themes"`
}

// Module templates provides support for looking up and using themes and
// their corresponding templates.
type Module struct {
	ConfigModule *config.Module
	Logger       *logger.Module

	internalStore *filestore.FileStore
	Stores        []ThemeStore

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
		m.internalStore = filestore.New(
			m.ConfigModule.DataPath(m.config.Themes.Path, internalThemeDir),
			time.Second*10,
			m.Logger.Error,
		)
		m.Stores = []ThemeStore{
			&defaultstore.DefaultStore{},
			filestore.New(m.config.Themes.Path, time.Second*10, m.Logger.Error),
			m.internalStore,
		}
		return nil
	}
}
