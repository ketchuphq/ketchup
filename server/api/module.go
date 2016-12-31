package api

import (
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	users2 "github.com/octavore/nagax/users"

	"github.com/octavore/press/db"
	"github.com/octavore/press/db/bolt"
	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
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
			{"/api/v1/pages/:uuid/routes", methodPost, m.Auth.MustWithAuth(m.UpdateRoutesByPage)},
			{"/api/v1/pages/:uuid/publish", methodPost, m.Auth.MustWithAuth(m.PublishPage)},
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

func (m *Module) GetUser(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	userUUID, ok := req.Context().Value(users2.UserTokenKey{}).(string)
	if !ok {
		m.Router.EmptyJSON(rw, http.StatusNotFound)
		return nil
	}
	user, err := m.DB.GetUser(userUUID)
	if err != nil {
		return err
	}
	user.HashedPassword = nil
	return router.Proto(rw, user)
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

func (m *Module) ListPages(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	pages, err := m.DB.ListPages()
	if err != nil {
		return err
	}
	db.SortPagesByUpdatedAt(pages, false)
	return router.Proto(rw, &api.ListPageResponse{
		Pages: pages,
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

func (m *Module) PublishPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
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
	now := time.Now().Unix()
	page.PublishedAt = &now
	err = m.DB.UpdatePage(page)
	if err != nil {
		return err
	}
	err = m.Content.ReloadRouter()
	if err != nil {
		m.Router.InternalError(rw, err)
	}
	return router.Proto(rw, page)
}

func (m *Module) Logout(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) error {
	m.Auth.Auth.Logout(rw, req)
	return nil
}
