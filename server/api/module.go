package api

import (
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	router2 "github.com/octavore/nagax/router"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/server/config"
	"github.com/ketchuphq/ketchup/server/content"
	"github.com/ketchuphq/ketchup/server/content/templates"
	"github.com/ketchuphq/ketchup/server/router"
	"github.com/ketchuphq/ketchup/server/tls"
	"github.com/ketchuphq/ketchup/server/users"
)

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
		auth := m.Auth.Auth.MustWithAuth
		routes := []struct {
			path, method string
			handle       router2.Handle
		}{
			{"/api/v1/pages/:uuid", methodGet, m.GetPage},
			{"/api/v1/pages/:uuid/contents", methodGet, m.GetRenderedPage},
			{"/api/v1/pages/:uuid/routes", methodGet, m.ListRoutesByPage},
			{"/api/v1/data/:key", methodGet, m.GetData},
			{"/api/v1/pages", methodGet, m.ListPages},
			{"/api/v1/routes", methodGet, m.ListRoutes},
			{"/api/v1/data", methodGet, auth(m.ListData)},

			{"/api/v1/user", methodGet, auth(m.GetUser)},
			{"/api/v1/settings/info", methodGet, auth(m.GetInfo)},
			{"/api/v1/settings/tls", methodGet, auth(m.GetTLS)},
			{"/api/v1/settings/tls", methodPost, auth(m.EnableTLS)},

			{"/api/v1/themes", methodGet, auth(m.ListThemes)},
			{"/api/v1/themes/:name", methodGet, auth(m.GetTheme)},
			{"/api/v1/themes/:name/updates", methodGet, auth(m.CheckThemeForUpdate)},
			{"/api/v1/themes/:name/update", methodPost, auth(m.UpdateTheme)},
			{"/api/v1/themes/:name/templates/:template", methodGet, auth(m.GetTemplate)},
			{"/api/v1/theme-registry", methodGet, auth(m.ThemeRegistry)},
			{"/api/v1/theme-install", methodPost, auth(m.InstallTheme)},

			{"/api/v1/pages", methodPost, auth(m.UpdatePage)},
			{"/api/v1/pages/:uuid", methodDelete, auth(m.DeletePage)},
			{"/api/v1/pages/:uuid/routes", methodPost, auth(m.UpdateRoutesByPage)},
			{"/api/v1/pages/:uuid/publish", methodPost, auth(m.PublishPage)},
			{"/api/v1/pages/:uuid/unpublish", methodPost, auth(m.UnpublishPage)},

			{"/api/v1/routes", methodPost, auth(m.UpdateRoute)},
			{"/api/v1/routes/:uuid", methodDelete, auth(m.DeleteRoute)},

			{"/api/v1/data", methodPost, auth(m.UpdateData)},
			{"/api/v1/data/:key", methodDelete, auth(m.DeleteData)},

			{"/api/v1/download-backup", methodGet, auth(m.GetBackup)},
			{"/api/v1/debug", methodGet, auth(m.Debug)},
			{"/api/v1/logout", methodGet, auth(m.Logout)},
		}
		for _, route := range routes {
			m.Router.WrappedHandle(route.method, route.path, route.handle)
		}
		return nil
	}
}

func (m *Module) Debug(rw http.ResponseWriter, req *http.Request, par router2.Params) error {
	return m.DB.Debug(rw)
}
