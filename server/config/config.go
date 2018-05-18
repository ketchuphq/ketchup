package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/config"
)

const defaultDataDir = "data"

type Config struct {
	DataDir string `json:"data_dir"` // themes, plugins
}

type Module struct {
	*config.Module
	Config Config
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		err := m.ReadConfig(&m.Config)
		if err != nil {
			return err
		}
		if m.Config.DataDir == "" {
			m.Config.DataDir = defaultDataDir
		}
		return nil
	}

	var tmpDir string
	var err error
	c.SetupTest = func() {
		tmpDir, err = ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}
		m.Config.DataDir = tmpDir
	}

	c.Stop = func() {
		if tmpDir != "" {
			err = os.RemoveAll(tmpDir)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func (m *Module) DataPath(p, backup string) string {
	if p == "" {
		p = backup
	}
	return path.Join(m.Config.DataDir, p)
}
