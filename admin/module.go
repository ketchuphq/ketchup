package admin

import (
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/router"
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
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.Router.Handle(basePath, m)
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
