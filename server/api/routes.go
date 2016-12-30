package api

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"

	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/util/errors"
)

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

func (m *Module) UpdateRoute(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	route := &models.Route{}
	err := jsonpb.Unmarshal(req.Body, route)
	if err != nil {
		return errors.Wrap(err)
	}
	err = m.DB.UpdateRoute(route)
	if err != nil {
		return err
	}
	err = m.Content.ReloadRouter()
	if err != nil {
		return nil
	}
	return router.Proto(rw, route)
}

func (m *Module) DeleteRoute(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	routeUUID := par.ByName("uuid")
	r, err := m.DB.GetRoute(routeUUID)
	if err != nil {
		return err
	}
	return m.DB.DeleteRoute(r)
}

func (m *Module) UpdateRoutesByPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	pageUUID := par.ByName("uuid")
	pb := &api.UpdatePageRoutesRequest{}
	err := jsonpb.Unmarshal(req.Body, pb)
	if err != nil {
		return errors.Wrap(err)
	}

	newList := map[string]*models.Route{}
	for _, route := range pb.GetRoutes() {
		route.Target = &models.Route_PageUuid{PageUuid: pageUUID}
		newList[route.GetUuid()] = route
	}

	routes, err := m.DB.ListRoutes()
	if err != nil {
		return err
	}

	filteredRoutes := []*models.Route{}
	for _, route := range routes {
		if route.GetPageUuid() == pageUUID {
			if newList[route.GetUuid()] == nil {
				err = m.DB.DeleteRoute(route)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, route := range newList {
		err = m.DB.UpdateRoute(route)
		if err != nil {
			return err
		}
	}

	err = m.Content.ReloadRouter()
	if err != nil {
		return err
	}

	return router.Proto(rw, &api.ListRouteResponse{
		Routes: filteredRoutes,
	})
}
