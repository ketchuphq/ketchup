package api

import (
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	router2 "github.com/octavore/nagax/router"
	"github.com/xenolf/lego/acme"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/config"
	"github.com/ketchuphq/ketchup/server/content"
	"github.com/ketchuphq/ketchup/server/content/templates"
	"github.com/ketchuphq/ketchup/server/files"
	"github.com/ketchuphq/ketchup/server/router"
	"github.com/ketchuphq/ketchup/server/tls"
	"github.com/ketchuphq/ketchup/server/users"
)

type mockableTemplateModule interface {
	CheckThemeForUpdate(name string) (bool, string, string, error)
	GetTemplate(theme, template string) (*models.ThemeTemplate, error)
	GetTheme(name string) (*models.Theme, string, error)
	InstallThemeFromPackage(p *packages.Package) error
	ListThemes() ([]*models.Theme, error)
	Registry() (*packages.Registry, error)
	SearchRegistry(themeName string) (*packages.Package, error)
	UpdateTheme(name, ref string) error
	GetRegistryURL() string
}

type mockableTLSModule interface {
	// CleanUp(domain, token, keyAuth string) error
	GetAllRegisteredDomains() ([]string, error)
	GetRegistration(domain string, withPrivateKey bool) (*tls.Registration, error)
	LoadCertResource(domain string) (*acme.CertificateResource, error)
	ObtainCert(email, domain string) error
	// Present(domain, token, keyAuth string) error
	// SaveRegistration(r *Registration) error
	// ServeHTTP(rw http.ResponseWriter, req *http.Request)
}

type Module struct {
	Router    *router.Module
	DB        *db.Module
	Auth      *users.Module
	Templates *templates.Module
	Content   *content.Module
	Config    *config.Module
	TLS       *tls.Module
	Logger    *logger.Module
	Files     *files.Module

	decoder *schema.Decoder
	version string

	// for tests
	templates mockableTemplateModule
	tls       mockableTLSModule
}

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"
)

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.templates = m.Templates
		m.tls = m.TLS
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

		// THESE ROUTES ARE PUBLIC
		// todo: fine-grained toggle
		publicRoutes := []struct {
			path, method string
			handle       router2.Handle
		}{
			{"/api/v1/pages/:uuid", methodGet, m.GetPage},
			{"/api/v1/pages/:uuid/contents", methodGet, m.GetRenderedPage},
			{"/api/v1/pages/:uuid/routes", methodGet, m.ListRoutesByPage},
			{"/api/v1/data/:key", methodGet, m.GetData},
			{"/api/v1/pages", methodGet, m.ListPages},
			{"/api/v1/routes", methodGet, m.ListRoutes},
		}
		for _, route := range publicRoutes {
			m.Router.WrappedHandle(route.method, route.path, route.handle)
		}

		privateRoutes := []struct {
			path, method string
			handle       router2.Handle
		}{
			{"/api/v1/data", methodGet, m.ListData},

			{"/api/v1/user", methodGet, m.GetUser},
			{"/api/v1/settings/info", methodGet, m.GetInfo},
			{"/api/v1/settings/tls", methodGet, m.GetTLS},
			{"/api/v1/settings/tls", methodPost, m.EnableTLS},

			{"/api/v1/themes", methodGet, m.ListThemes},
			{"/api/v1/themes/:name", methodGet, m.GetTheme},
			{"/api/v1/themes/:name/updates", methodGet, m.CheckThemeForUpdate},
			{"/api/v1/themes/:name/update", methodPost, m.UpdateTheme},
			{"/api/v1/themes/:name/templates/:template", methodGet, m.GetTemplate},
			{"/api/v1/theme-registry", methodGet, m.ThemeRegistry},
			{"/api/v1/theme-install", methodPost, m.InstallTheme},

			{"/api/v1/pages", methodPost, m.UpdatePage},
			{"/api/v1/pages/:uuid", methodDelete, m.DeletePage},
			{"/api/v1/pages/:uuid/routes", methodPost, m.UpdateRoutesByPage},
			{"/api/v1/pages/:uuid/publish", methodPost, m.PublishPage},
			{"/api/v1/pages/:uuid/unpublish", methodPost, m.UnpublishPage},

			{"/api/v1/routes", methodPost, m.UpdateRoute},
			{"/api/v1/routes/:uuid", methodDelete, m.DeleteRoute},

			{"/api/v1/data", methodPost, m.UpdateData},
			{"/api/v1/data/:key", methodDelete, m.DeleteData},

			{"/api/v1/files", methodGet, m.ListFiles},
			{"/api/v1/files", methodPost, m.UploadFile},
			{"/api/v1/files/:uuid", methodGet, m.GetFile},
			{"/api/v1/files/:uuid", methodDelete, m.DeleteFile},

			{"/api/v1/download-backup", methodGet, m.GetBackup},
			{"/api/v1/debug", methodGet, m.Debug},
			{"/api/v1/logout", methodGet, m.Logout},
		}
		for _, route := range privateRoutes {
			m.Router.WrappedHandle(route.method, route.path, auth(route.handle))
		}
		return nil
	}
}

func (m *Module) Debug(rw http.ResponseWriter, req *http.Request, par router2.Params) error {
	return m.DB.Debug(rw)
}
