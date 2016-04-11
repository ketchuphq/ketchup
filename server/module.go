package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/press/proto/press/models"
)

type Module struct {
	*http.ServeMux
}

func (m *Module) Setup(c *service.Config) {
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//req.URL.Path
}

func (m *Module) BuildRouter() {
	router := httprouter.New()
	router.Handle("GET", "/foo/bar", m.Route(nil))
}

func (m *Module) Route(r *models.Route) httprouter.Handle {
	return func(http.ResponseWriter, *http.Request, httprouter.Params) {
		// lookup FILE
		// lookup PAGE
		// or REDIRECT
	}
}
