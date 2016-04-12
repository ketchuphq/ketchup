package plugins

import (
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/protobuf/proto"
	"github.com/yuin/gopher-lua"
)

type PluginRouteLoader struct {
	routes []*models.Route
}

func (l *PluginRouteLoader) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"setRoute": l.set,
	})
	L.Push(mod)
	return 1
}

// setRoute(path, uuid)
func (l *PluginRouteLoader) set(L *lua.LState) int {
	path := L.CheckString(1)
	uuid := L.CheckString(2)
	l.routes = append(l.routes, &models.Route{
		Path:   proto.String(path),
		Target: &models.Route_PageUuid{PageUuid: uuid},
	})
	return 0
}
