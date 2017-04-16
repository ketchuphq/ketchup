package api

import (
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/octavore/ketchup/db"
	"github.com/octavore/ketchup/proto/ketchup/api"
	"github.com/octavore/ketchup/server/config"
	"github.com/octavore/ketchup/server/content"
	"github.com/octavore/ketchup/server/content/templates"
	"github.com/octavore/ketchup/server/router"
	"github.com/octavore/ketchup/server/tls"
	"github.com/octavore/ketchup/server/users"
)

var KetchupVersion = ""

type Module struct {
	Router    *router.Module
	DB        *db.Module
	Auth      *users.Module
	Templates *templates.Module
	Content   *content.Module
	Config    *config.Module
	TLS       *tls.Module
	Logger    *logger.Module

	decoder *schema.Decoder
	version string
}

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"
)

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.decoder = schema.NewDecoder()
		m.decoder.SetAliasTag("json")
		m.decoder.IgnoreUnknownKeys(true)
		m.decoder.RegisterConverter(
			api.ListPageRequest_all,
			func(val string) reflect.Value {
				v := api.ListPageRequest_ListPageFilter_value[val]
				return reflect.ValueOf(v)
			})

		r := m.Router.Subrouter("/api/v1/")
		routes := []struct {
			path, method string
			handle       users.Handle
		}{
			{"/api/v1/pages/:uuid", methodGet, m.GetPage},
			{"/api/v1/pages/:uuid/contents", methodGet, m.GetRenderedPage},
			{"/api/v1/pages/:uuid/routes", methodGet, m.ListRoutesByPage},
			{"/api/v1/pages", methodGet, m.ListPages},
			{"/api/v1/routes", methodGet, m.ListRoutes},

			{"/api/v1/user", methodGet, m.Auth.MustWithAuth(m.GetUser)},
			{"/api/v1/settings/info", methodGet, m.Auth.MustWithAuth(m.GetInfo)},
			{"/api/v1/settings/tls", methodGet, m.Auth.MustWithAuth(m.GetTLS)},
			{"/api/v1/settings/tls", methodPost, m.Auth.MustWithAuth(m.EnableTLS)},

			{"/api/v1/themes", methodGet, m.Auth.MustWithAuth(m.ListThemes)},
			{"/api/v1/themes/:name", methodGet, m.Auth.MustWithAuth(m.GetTheme)},
			{"/api/v1/themes/:name/templates/:template", methodGet, m.Auth.MustWithAuth(m.GetTemplate)},
			{"/api/v1/theme-registry", methodGet, m.Auth.MustWithAuth(m.ThemeRegistry)},
			{"/api/v1/theme-install", methodGet, m.Auth.MustWithAuth(m.InstallTheme)},

			{"/api/v1/pages", methodPost, m.Auth.MustWithAuth(m.UpdatePage)},
			{"/api/v1/pages/:uuid", methodDelete, m.Auth.MustWithAuth(m.DeletePage)},
			{"/api/v1/pages/:uuid/routes", methodPost, m.Auth.MustWithAuth(m.UpdateRoutesByPage)},
			{"/api/v1/pages/:uuid/publish", methodPost, m.Auth.MustWithAuth(m.PublishPage)},
			{"/api/v1/pages/:uuid/unpublish", methodPost, m.Auth.MustWithAuth(m.UnpublishPage)},

			{"/api/v1/routes", methodPost, m.Auth.MustWithAuth(m.UpdateRoute)},
			{"/api/v1/routes/:uuid", methodDelete, m.Auth.MustWithAuth(m.DeleteRoute)},

			{"/api/v1/download-backup", methodGet, m.Auth.MustWithAuth(m.GetBackup)},
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
