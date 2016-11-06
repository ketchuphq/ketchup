package api

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"

	"github.com/octavore/press/db"
	"github.com/octavore/press/db/bolt"
	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/server/users"
)

type Handle func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error

type Module struct {
	Router *router.Module
	DB     *db.Module
	Auth   *users.Module
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		r := m.Router.Subrouter("/api/v1/")
		routes := []struct {
			path, method string
			handle       Handle
		}{
			{"/api/v1/pages/:uuid", "GET", m.GetPage},
			{"/api/v1/pages/:uuid/routes", "GET", m.ListRoutesByPage},
			{"/api/v1/pages", "POST", m.UpdatePage},
			{"/api/v1/pages", "GET", m.ListPages},
			{"/api/v1/routes", "GET", m.ListRoutes},
			{"/api/v1/routes", "POST", m.UpdateRoute},
			{"/api/v1/debug", "GET", m.Debug},
		}
		for _, route := range routes {
			r.Handle(route.method, route.path, m.wrap(route.handle))
		}
		return nil
	}
}

func (m *Module) wrap(h Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) {
		err := h(rw, req, par)
		if err != nil {
			router.InternalError(rw, err)
		}
	}
}

func (m *Module) Debug(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	return m.DB.Debug(rw)
}

func (m *Module) GetPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	uuid := par.ByName("uuid")
	if uuid == "" {
		return router.ErrNotFound
	}
	page, err := m.DB.GetPage(uuid)
	if _, ok := err.(bolt.ErrNoKey); ok {
		return router.ErrNotFound
	}
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) ListPages(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	pages, err := m.DB.ListPages()
	if err != nil {
		return err
	}
	return router.Proto(rw, &api.ListPageResponse{
		Pages: pages,
	})
}

func (m *Module) ListRoutes(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	routes, err := m.DB.ListRoutes()
	if err != nil {
		return err
	}
	return router.Proto(rw, &api.ListRouteResponse{
		Routes: routes,
	})
}

func (m *Module) ListRoutesByPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	routes, err := m.DB.ListRoutes()
	if err != nil {
		return err
	}
	pageUUID := par.ByName("uuid")
	filteredRoutes := []*models.Route{}
	for _, route := range routes {
		if route.GetPageUuid() == pageUUID {
			filteredRoutes = append(filteredRoutes, route)
		}
	}
	return router.Proto(rw, &api.ListRouteResponse{
		Routes: filteredRoutes,
	})
}

func (m *Module) UpdatePage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page := &models.Page{}
	err := jsonpb.Unmarshal(req.Body, page)
	if err != nil {
		return err
	}
	err = m.DB.UpdatePage(page)
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) UpdateRoute(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	route := &models.Route{}
	err := jsonpb.Unmarshal(req.Body, route)
	if err != nil {
		return err
	}
	err = m.DB.UpdateRoute(route)
	if err != nil {
		return err
	}
	return router.Proto(rw, route)
}
