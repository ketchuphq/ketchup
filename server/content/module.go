package content

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/octavore/press/db"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/router"
)

// Module server is responsible for serving published content
type Module struct {
	Router *router.Module
	DB     *db.Module
	Logger *logger.Module

	router http.Handler
}

var _ service.Module = &Module{}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.Start = func() {
		var err error
		m.router, err = m.BuildRouter()
		if err != nil {
			panic(err)
		}
		m.Router.Handle("/", m)
	}
}

// ServeHTTP is a layer of indirection to allow us
// to replace the router at runtime.
func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if m.router != nil {
		m.router.ServeHTTP(rw, req)
	}
}

// BuildRouter returns a handler configured to serve content
// from the DB
func (m *Module) BuildRouter() (http.Handler, error) {
	router := httprouter.New()
	routes, err := m.DB.ListRoutes()
	if err != nil {
		return nil, err
	}
	for _, route := range routes {
		m.Logger.Info("found route:", route)
		switch tgt := route.GetTarget().(type) {
		case *models.Route_File:
			m.Logger.Info("registered file route:", route.GetPath())
			router.Handle("GET", route.GetPath(), func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
				http.ServeFile(rw, req, tgt.File)
			})
		case *models.Route_PageUuid:
			m.Logger.Info("registered uuid route:", route.GetPath())
			router.Handle("GET", route.GetPath(), func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
				page, err := m.DB.GetPage(tgt.PageUuid)
				if err != nil {
					m.Logger.Errorf("error serving page %+v: %+v", route, err)
				} else {
					rw.Write([]byte(page.GetData()))
				}
			})
		default:
			m.Logger.Errorf("unable to register %s", route.GetUuid())
		}
	}
	return router, nil
}
