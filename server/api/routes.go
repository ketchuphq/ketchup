package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/util/errors"
)

var (
	re1 = regexp.MustCompile(`[^a-zA-Z0-9\/]`)
	re2 = regexp.MustCompile(`^-+`)
	re3 = regexp.MustCompile(`-+$`)
	re4 = regexp.MustCompile(`\/\/+`)
)

func formatRoute(r *models.Route) *models.Route {
	if r.Path == nil {
		return r
	}
	p := "/" + strings.Trim(r.GetPath(), "/")
	p = strings.ToLower(p)
	p = re1.ReplaceAllString(p, "-")
	p = re2.ReplaceAllString(p, "")
	p = re3.ReplaceAllString(p, "")
	p = re4.ReplaceAllString(p, "/")
	r.Path = &p
	return r
}

// ListRoutes returns all routes.
func (m *Module) ListRoutes(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	routes, err := m.DB.ListRoutes()
	if err != nil {
		return err
	}
	return router.Proto(rw, &api.ListRouteResponse{
		Routes: routes,
	})
}

// ListRoutesByPage returns all routes for a given page, identified by uuid in the parameter.
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

// UpdateRoute updates the given route. The path is sanitized, and the content router is
// reloaded after the route is saved.
func (m *Module) UpdateRoute(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	route := &models.Route{}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrap(err)
	}
	err = json.Unmarshal(b, route)
	if err != nil {
		return errors.Wrap(err)
	}
	err = m.DB.UpdateRoute(formatRoute(route))
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

// UpdateRoutesByPage takes a list of routes and
// /api/v1/pages/:uuid/routes
func (m *Module) UpdateRoutesByPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	pageUUID := par.ByName("uuid")
	pb := &api.UpdatePageRoutesRequest{}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrap(err)
	}
	err = json.Unmarshal(b, pb)
	if err != nil {
		return errors.Wrap(err)
	}

	// newList contains a map of uuid to routes for routes to add.
	newList := map[string]*models.Route{}
	for _, route := range pb.GetRoutes() {
		route.Target = &models.Route_PageUuid{PageUuid: pageUUID}
		// bug: need to handle multiple routes with no uuid
		newList[route.GetUuid()] = route
	}

	routes, err := m.DB.ListRoutes()
	if err != nil {
		return err
	}

	// loop over all routes
	filteredRoutes := []*models.Route{}
	for _, route := range routes {
		if route.GetPageUuid() == pageUUID {
			if newList[route.GetUuid()] == nil {
				// if we're not adding this existing route, then delete it
				err = m.DB.DeleteRoute(route)
				if err != nil {
					return err
				}
			}
		}
	}

	// save all the routes specified.
	for _, route := range newList {
		err = m.DB.UpdateRoute(formatRoute(route))
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
