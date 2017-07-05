package api

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/server/router"
	"github.com/ketchuphq/ketchup/util/errors"
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
	theme, _, err := m.Templates.GetTheme(name)
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
	err := jsonpb.Unmarshal(req.Body, r)
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
	if p == nil || p.GetVcsUrl() != r.GetVcsUrl() {
		return errors.New("Theme %s not found %s, %s", r.GetName(),
			p.GetVcsUrl(), r.GetVcsUrl())
	}

	m.Logger.Infof("cloning package %s from %s", p.GetName(), p.GetVcsUrl())

	return m.Templates.InstallThemeFromPackage(p)
}
