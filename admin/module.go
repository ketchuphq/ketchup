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
	p := strings.TrimPrefix(req.URL.Path, basePath)
	ext := path.Ext(p)
	if ext == "" {
		p = path.Join(p, "index.html")
		ext = path.Ext(p)
	}

	b, err := Asset(p)
	if err != nil && strings.Contains(err.Error(), "not found") {
		p = "index.html"
		b, err = Asset(p)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", mime.TypeByExtension(p))
	rw.Write(b)
}
