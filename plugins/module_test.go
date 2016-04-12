package plugins

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/press/proto/press/models"
)

func TestModule(t *testing.T) {
	app := &Module{
		config: &RouteConfig{Plugins: "test_plugins"},
	}
	err := app.loadPlugins()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(app.Plugins, []*Plugin{
		{Path: "test_plugins/test_plugin/plugin.lua"},
	}) {
		t.Error(app.Plugins)
	}
}

func TestPlugin(t *testing.T) {
	p := &Plugin{"test_plugins/test_plugin/plugin.lua"}
	loader := &PluginRouteLoader{}
	err := p.Routes(loader)
	if err != nil {
		t.Error(err)
	}
	exp := []*models.Route{{
		Path:   proto.String("hello"),
		Target: &models.Route_PageUuid{PageUuid: "world"},
	}, {
		Path:   proto.String("goodbye"),
		Target: &models.Route_PageUuid{PageUuid: "moon"},
	}}
	if !reflect.DeepEqual(exp, loader.routes) {
		t.Error(loader.routes)
	}
}
