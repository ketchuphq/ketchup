package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/packages"
	"github.com/octavore/press/server/router"
	"github.com/octavore/press/util/errors"
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

func (m *Module) ThemeRegistry(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	lst, err := m.Pkg.FetchDefaultRegistry()
	if err != nil {
		return err
	}
	return m.Router.ProtoOK(rw, lst)
}

type InstallThemeRequest struct {
	Package string `json:"package"`
}

func (m *Module) InstallTheme(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	r := &InstallThemeRequest{}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrap(err)
	}
	err = json.Unmarshal(b, r)
	if err != nil {
		return errors.Wrap(err)
	}
	if r.Package == "" {
		return fmt.Errorf("Theme name is required.")
	}

	// loop over all registry packages to find the theme package
	registry, err := m.Pkg.FetchDefaultRegistry()
	if err != nil {
		return err
	}
	var pkg *packages.Package
	for _, p := range registry.GetPackages() {
		if p.GetName() == r.Package {
			if p.GetType() != packages.Package_theme {
				return fmt.Errorf("%s is not a theme", r.Package)
			}
			pkg = p
			break
		}
	}
	if pkg == nil {
		return fmt.Errorf("Theme %s not found", r.Package)
	}

	m.Logger.Infof("cloning package %s from %s", pkg.GetName(), pkg.GetVcsUrl())
	err = m.Pkg.Clone(pkg)
	if err != nil {
		return err
	}

	return nil
}
