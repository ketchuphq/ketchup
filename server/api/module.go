package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"

	"github.com/octavore/press/db"
	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/server/content"
	"github.com/octavore/press/server/content/templates"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/server/users"
)

type Module struct {
	Router    *router.Module
	DB        *db.Module
	Auth      *users.Module
	Templates *templates.Module
	Content   *content.Module
}

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"
)

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		r := m.Router.Subrouter("/api/v1/")
		routes := []struct {
			path, method string
			handle       users.Handle
		}{
			{"/api/v1/pages/:uuid", methodGet, m.GetPage},
			{"/api/v1/pages/:uuid/routes", methodGet, m.ListRoutesByPage},
			{"/api/v1/pages", methodGet, m.ListPages},
			{"/api/v1/routes", methodGet, m.ListRoutes},

			{"/api/v1/user", methodGet, m.Auth.MustWithAuth(m.GetUser)},
			{"/api/v1/themes", methodGet, m.Auth.MustWithAuth(m.ListThemes)},
			{"/api/v1/themes/:name", methodGet, m.Auth.MustWithAuth(m.GetTheme)},

			{"/api/v1/pages", methodPost, m.Auth.MustWithAuth(m.UpdatePage)},
			{"/api/v1/pages/:uuid", methodDelete, m.Auth.MustWithAuth(m.DeletePage)},
			{"/api/v1/pages/:uuid/routes", methodPost, m.Auth.MustWithAuth(m.UpdateRoutesByPage)},
			{"/api/v1/pages/:uuid/publish", methodPost, m.Auth.MustWithAuth(m.PublishPage)},
			{"/api/v1/pages/:uuid/unpublish", methodPost, m.Auth.MustWithAuth(m.UnpublishPage)},

			{"/api/v1/routes", methodPost, m.Auth.MustWithAuth(m.UpdateRoute)},
			{"/api/v1/routes/:uuid", methodDelete, m.Auth.MustWithAuth(m.DeleteRoute)},
			{"/api/v1/debug", methodGet, m.Auth.MustWithAuth(m.Debug)},
			{"/api/v1/logout", methodGet, m.Auth.MustWithAuth(m.Logout)},
		}
		for _, route := range routes {
			r.Handle(route.method, route.path, m.wrap(route.handle))
		}
		return nil
	}
}

func (m *Module) wrap(h users.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) {
		err := h(rw, req, par)
		if err != nil {
			m.Router.InternalError(rw, err)
		}
	}
}

func (m *Module) Debug(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	return m.DB.Debug(rw)
}

func (m *Module) ListThemes(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) error {
	themes, err := m.Templates.ListThemes()
	if err != nil {
		return err
	}
	return router.Proto(rw, &api.ListThemeResponse{Themes: themes})
}

func (m *Module) GetTheme(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	name := par.ByName("name")
	theme, err := m.Templates.GetTheme(name)
	if err != nil {
		return err
	}
	if theme == nil {
		return router.ErrNotFound
	}
	return router.Proto(rw, theme)
}
