package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/server/router"
)

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
	// todo: fetch remote theme info if available
	if theme == nil {
		return router.ErrNotFound
	}
	return router.Proto(rw, theme)
}

func (m *Module) GetTemplate(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	name := par.ByName("name")
	template := par.ByName("template")
	theme, err := m.Templates.GetTemplate(name, template)
	if err != nil {
		return err
	}

	// todo: fetch remote theme info if available
	if theme == nil {
		return router.ErrNotFound
	}
	return router.Proto(rw, theme)
}
