package plugins

import (
	"os"
	"path/filepath"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/config"
	"github.com/yuin/gopher-lua"
)

const pluginFileName = "plugin.lua"

type RouteConfig struct {
	Plugins string `json:"plugins"`
}

type Module struct {
	Config *config.Module

	Plugins []*Plugin
	config  *RouteConfig
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.config = &RouteConfig{}
		m.Config.ReadConfig(m.config)
		return nil
	}
	c.SetupTest = func() {
		m.config.Plugins = "test_plugins"
	}
	c.Start = func() {
		m.loadPlugins()
	}
}

// load all plugins from the plugin folder
func (m *Module) loadPlugins() error {
	if m.config.Plugins == "" {
		return nil
	}
	f, err := os.Open(m.config.Plugins)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return nil
	}
	lst, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, dir := range lst {
		if !dir.IsDir() {
			continue
		}
		p := filepath.Join(m.config.Plugins, dir.Name(), pluginFileName)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		m.Plugins = append(m.Plugins, &Plugin{Path: p})
	}
	return nil
}

type Plugin struct {
	Path string
}

func (p *Plugin) Routes(loader *PluginRouteLoader) error {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("press", loader.Loader)
	err := L.DoFile(p.Path)
	if err != nil {
		return err
	}
	fn := L.GetGlobal("routes")
	if _, ok := fn.(*lua.LNilType); ok {
		return nil
	}
	err = L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	})
	if err != nil {
		return err
	}
	return nil
}
