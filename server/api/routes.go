package api

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

var (
	re1 = regexp.MustCompile(`[^a-zA-Z0-9\/]`)
	re2 = regexp.MustCompile(`^-+`)
	re3 = regexp.MustCompile(`-+$`)
	re4 = regexp.MustCompile(`-*\/\/*-*\/*`)
)

func formatRoute(r *models.Route) *models.Route {
	if r.Path == nil {
		return r
	}
	p := strings.ToLower(r.GetPath())
	p = re1.ReplaceAllString(p, "-")
	p = re2.ReplaceAllString(p, "")
	p = re3.ReplaceAllString(p, "")
	p = re4.ReplaceAllString(p, "/")

	p = "/" + strings.Trim(p, "/")
	r.Path = &p
	return r
}

// ListRoutes returns all routes.
func (m *Module) ListRoutes(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	opts := &api.ListRouteRequest{}
	err := req.ParseForm()
	if err != nil {
		return err
	}
	err = m.decoder.Decode(opts, req.Form)
	if err != nil {
		return err
	}
	routes, err := m.DB.ListRoutes(opts)
	if err != nil {
		return err
	}
	return router.ProtoOK(rw, &api.ListRouteResponse{
		Routes: db.SortRoutesByPath(routes),
	})
}

// ListRoutesByPage returns all routes for a given page, identified by uuid in the parameter.
func (m *Module) ListRoutesByPage(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	pageUUID := par.ByName("uuid")
	routes, err := m.DB.ListRoutes(&api.ListRouteRequest{
		Options: &api.ListRouteRequest_ListRouteOptions{
			PageUuid: &pageUUID,
		},
	})
	if err != nil {
		return err
	}
	return router.ProtoOK(rw, &api.ListRouteResponse{
		Routes: db.SortRoutesByPath(routes),
	})
}

// UpdateRoute updates the given route. The path is sanitized, and the content router is
// reloaded after the route is saved.
func (m *Module) UpdateRoute(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	route := &models.Route{}
	err := jsonpb.Unmarshal(req.Body, route)
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
	return router.ProtoOK(rw, route)
}

// DeleteRoute deletes the given route.
func (m *Module) DeleteRoute(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	routeUUID := par.ByName("uuid")
	r, err := m.DB.GetRoute(routeUUID)
	if err != nil {
		return err
	}
	return m.DB.DeleteRoute(r)
}

// UpdateRoutesByPage takes a list of routes and sets it for the given page, deleting
// any existing routes that aren't in the new list.
func (m *Module) UpdateRoutesByPage(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	pageUUID := par.ByName("uuid")
	pb := &api.UpdatePageRoutesRequest{}
	err := jsonpb.Unmarshal(req.Body, pb)
	if err != nil {
		return errors.Wrap(err)
	}

	oldRoutes, err := m.DB.ListRoutes(&api.ListRouteRequest{
		Options: &api.ListRouteRequest_ListRouteOptions{
			PageUuid: &pageUUID,
		},
	})
	if err != nil {
		return err
	}

	// newList contains a map of uuid to routes for routes to add.
	newList := map[string]*models.Route{}
	for _, route := range pb.GetRoutes() {
		route.Target = &models.Route_PageUuid{PageUuid: pageUUID}

		// if no uuid, then just save the route
		if route.GetUuid() == "" {
			err = m.DB.UpdateRoute(formatRoute(route))
			if err != nil {
				return err
			}
		} else {
			newList[route.GetUuid()] = route
		}
	}

	// loop over all routes
	for _, route := range oldRoutes {
		if newList[route.GetUuid()] == nil {
			// if we're not adding this existing route, then delete it
			err = m.DB.DeleteRoute(route)
			if err != nil {
				return err
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

	return m.Content.ReloadRouter()
}
