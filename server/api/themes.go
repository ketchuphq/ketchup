package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/octavore/ketchup/proto/ketchup/api"
	"github.com/octavore/ketchup/server/router"
	"github.com/octavore/ketchup/util/errors"
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
	// m.DB.GetThemeByName(name
	// m.Templates
	// additional data for a theme: which store it is in?
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

func (m *Module) ThemeRegistry(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	r, err := m.Templates.Registry()
	if err != nil {
		return err
	}
	return m.Router.ProtoOK(rw, r)
}

// InstallTheme installs a theme from a registry
func (m *Module) InstallTheme(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	r := &api.InstallThemeRequest{}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrap(err)
	}
	err = json.Unmarshal(b, r)
	if err != nil {
		return errors.Wrap(err)
	}

	if r.GetName() == "" {
		return errors.New("Theme name is required.")
	}

	// search the registry for the theme package
	// install the package
	p, err := m.Templates.SearchRegistry(r.GetName())
	if err != nil {
		return errors.New("error searching registry: %s", err)
	}
	if p == nil {
		return errors.New("Theme %s not found", r.GetName())
	}

	m.Logger.Infof("cloning package %s from %s", p.GetName(), p.GetVcsUrl())
	err = m.Templates.InstallThemeFromPackage(p)
	if err != nil {
		return errors.New("error searching registry: %s", err)
	}

	return nil
}
