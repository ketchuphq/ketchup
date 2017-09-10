package api

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/nagax/router"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/util/errors"
)

func (m *Module) ListThemes(rw http.ResponseWriter, req *http.Request, _ router.Params) error {
	themes, err := m.Templates.ListThemes()
	if err != nil {
		return err
	}
	return router.ProtoOK(rw, &api.ListThemeResponse{Themes: themes})
}

func (m *Module) GetTheme(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	name := par.ByName("name")
	theme, ref, err := m.Templates.GetTheme(name)
	if err != nil {
		return err
	}
	// todo: fetch remote theme info if available
	if theme == nil {
		return router.ErrNotFound
	}
	return router.ProtoOK(rw, &api.GetThemeResponse{
		Theme: theme,
		Ref:   &ref,
	})
}

func (m *Module) GetTemplate(rw http.ResponseWriter, req *http.Request, par router.Params) error {
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
	return router.ProtoOK(rw, theme)
}

func (m *Module) ThemeRegistry(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	r, err := m.Templates.Registry()
	if err != nil {
		return err
	}
	return router.ProtoOK(rw, r)
}

// InstallTheme installs a theme from a registry
func (m *Module) InstallTheme(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	r := &api.InstallThemeRequest{}
	err := jsonpb.Unmarshal(req.Body, r)
	if err != nil {
		return errors.Wrap(err)
	}

	if r.GetName() == "" {
		return errors.New("Theme name is required.")
	}

	// search the registry for the theme package
	p, err := m.Templates.SearchRegistry(r.GetName())
	if err != nil {
		return errors.New("error searching registry: %s", err)
	}
	if p == nil || p.GetVcsUrl() != r.GetVcsUrl() {
		return errors.New("Theme %s not found %s, %s", r.GetName(),
			p.GetVcsUrl(), r.GetVcsUrl())
	}

	m.Logger.Infof("cloning package %s from %s", p.GetName(), p.GetVcsUrl())

	// install the package
	return m.Templates.InstallThemeFromPackage(p)
}

func (m *Module) CheckThemeForUpdate(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	name := par.ByName("name")
	_, oldRef, currentRef, err := m.Templates.CheckThemeForUpdate(name)
	if err != nil {
		return errors.Wrap(err)
	}

	return router.JSON(rw, http.StatusOK, &api.CheckThemeForUpdateResponse{
		OldRef:     &oldRef,
		CurrentRef: &currentRef,
	})
}

func (m *Module) UpdateTheme(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	themeName := par.ByName("name")
	r := &api.UpdateThemeRequest{}
	err := jsonpb.Unmarshal(req.Body, r)
	if err != nil {
		return errors.Wrap(err)
	}
	if r.GetName() != "" && r.GetName() != themeName {
		return errors.New("theme name mismatch")
	}
	return m.Templates.UpdateTheme(r.GetName(), r.GetRef())
}
