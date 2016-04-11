package static

import (
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/octavore/naga/service"
)

type Module struct {
	Router *http.ServeMux
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.Router = http.NewServeMux()
		m.Router.Handle("/", m)
		return nil
	}
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := strings.TrimPrefix(req.URL.Path, "/")
	ext := path.Ext(p)
	if ext == "" {
		p = path.Join(p, "index.html")
		ext = path.Ext(p)
	}

	b, err := Asset(p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", mime.TypeByExtension(p))
	rw.Write(b)
}
