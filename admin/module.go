package admin

import (
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/users"
)

const basePath = "/admin/"

var staticDirs = []string{
	"/admin/js/",
	"/admin/css/",
	"/admin/vendor/",
	"/admin/images/",
}

type Module struct {
	Router *router.Module
	Auth   *users.Module
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.Router.Root.Handle(basePath, m)
		m.Router.Root.HandleFunc("/admin/logout", m.handleLogout)
		m.Router.Root.HandleFunc("/admin/preview", func(rw http.ResponseWriter, req *http.Request) {
			err := m.Auth.MustWithAuth(m.handlePreview)(rw, req, nil)
			if err != nil {
				m.Router.HandleError(rw, req, err)
			}
		})
		return nil
	}
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	isStaticAsset := false
	for _, dir := range staticDirs {
		isStaticAsset = isStaticAsset || strings.HasPrefix(req.URL.Path, dir)
	}

	p := strings.TrimPrefix(req.URL.Path, basePath)
	ext := path.Ext(p)

	if !isStaticAsset {
		p = "index.html"
		ext = path.Ext(p)
	}

	b, err := Asset(p)
	if err != nil && strings.Contains(err.Error(), "not found") {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", mime.TypeByExtension(ext))
	rw.Write(b)
}

func (m *Module) handleLogout(rw http.ResponseWriter, req *http.Request) {
	m.Auth.Logout(rw, req)
	http.Redirect(rw, req, "/admin", http.StatusTemporaryRedirect)
}
