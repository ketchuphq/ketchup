package api

import (
	"net/http"

	_bolt "github.com/boltdb/bolt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/press/bolt"
	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/router"
)

type Handle func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error

type Module struct {
	Router *router.Module
	Bolt   *bolt.Module
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		r := m.Router.Subrouter("/api/v1/")
		routes := []struct {
			path, method string
			handle       Handle
		}{
			{"/api/v1/pages/:uuid", "GET", m.GetPage},
			{"/api/v1/pages", "POST", m.UpdatePage},
			{"/api/v1/routes", "GET", m.ListRoutes},
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
	return m.Bolt.DB.View(func(tx *_bolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *_bolt.Bucket) error {
			return bucket.ForEach(func(key, value []byte) error {
				rw.Write(key)
				rw.Write([]byte("\n"))
				rw.Write(value)
				rw.Write([]byte("\n"))
				return nil
			})
		})
	})
}

func (m *Module) GetPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	uuid := par.ByName("uuid")
	if uuid == "" {
		return router.ErrNotFound
	}
	page, err := m.Bolt.GetPage(uuid)
	if _, ok := err.(bolt.ErrNoKey); ok {
		return router.ErrNotFound
	}
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) ListRoutes(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	routes, err := m.Bolt.ListRoutes()
	if err != nil {
		return err
	}
	return router.Proto(rw, &api.ListRouteResponse{
		Routes: routes,
	})
}

func (m *Module) UpdatePage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page := &models.Page{}
	err := jsonpb.Unmarshal(req.Body, page)
	if err != nil {
		return err
	}
	err = m.Bolt.UpdatePage(page)
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}
