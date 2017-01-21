package config

import (
	"path"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/config"
)

type Config struct {
	DataDir string `json:"data_dir"` // themes, plugins
}

type Module struct {
	*config.Module
	Config Config
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		return m.ReadConfig(&m.Config)
	}
}

func (m *Module) DataPath(p, backup string) string {
	if p == "" {
		p = backup
	}
	return path.Join(m.Config.DataDir, p)
}
