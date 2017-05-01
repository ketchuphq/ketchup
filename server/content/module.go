package content

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/templates"
	"github.com/ketchuphq/ketchup/server/router"
)

// Module server is responsible for serving published content
type Module struct {
	Router    *router.Module
	DB        *db.Module
	Logger    *logger.Module
	Templates *templates.Module

	contentRouter http.Handler
}

var _ service.Module = &Module{}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.Start = func() {
		var err error
		m.contentRouter, _, err = m.buildRouter()
		if err != nil {
			panic(err)
		}
		m.Router.Handle("/", m)
	}
}

// ServeHTTP is a layer of indirection to allow us
// to replace the router at runtime.
func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if m.contentRouter != nil {
		m.contentRouter.ServeHTTP(rw, req)
	}
}

// ReloadRouter recreates the router to add/remove routes that have change
// and replaces the existing router.
func (m *Module) ReloadRouter() error {
	newRouter, _, err := m.buildRouter()
	if err != nil {
		return err
	}
	m.contentRouter = newRouter
	return nil
}

func (m *Module) serveFile(file string) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		http.ServeFile(rw, req, file)
	}
}

func (m *Module) servePage(pageUUID string) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		page, err := m.DB.GetPage(pageUUID)
		if page.PublishedAt == nil {
			m.Router.NotFound(rw)
			return
		}
		if err == nil {
			err = m.render(rw, page)
		}
		if err != nil {
			m.Logger.Errorf("error serving page %s [%s]: %+v", req.URL.Path, pageUUID, err)
			m.Router.InternalError(rw, err)
		}
	}
}

// buildRouter returns a handler configured to serve content.
func (m *Module) buildRouter() (http.Handler, map[string]bool, error) {
	rt := httprouter.New()
	routes, err := m.DB.ListRoutes()
	if err != nil {
		return nil, nil, err
	}
	activeRoutes := map[string]bool{}
	rt.NotFound = m.Templates
	for _, route := range routes {
		m.Logger.Info("found route:", route)
		if route.GetPath() == "" {
			m.Logger.Warningf("no path for route: %s", route.GetUuid())
			continue
		}
		if route.GetPath()[0] != '/' {
			m.Logger.Warningf("invalid path %q for route: %s", route.GetPath(), route.GetUuid())
			continue
		}
		if _, ok := activeRoutes[route.GetPath()]; ok {
			m.Logger.Warningf("failed to register duplicate route %q", route.GetPath())
			continue
		}
		activeRoutes[route.GetPath()] = true

		switch tgt := route.GetTarget().(type) {
		case *models.Route_File:
			m.Logger.Info("registered file route:", route.GetPath())
			rt.Handle("GET", route.GetPath(), m.serveFile(tgt.File))
		case *models.Route_PageUuid:
			m.Logger.Info("registered uuid route:", route.GetPath())
			rt.Handle("GET", route.GetPath(), m.servePage(tgt.PageUuid))
		default:
			m.Logger.Errorf("unable to register %s", route.GetUuid())
		}
	}
	m.Logger.Info("done rebuilding router")
	return rt, activeRoutes, nil
}
