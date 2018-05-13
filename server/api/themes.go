package api

import (
	"net/http"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
)

func (m *Module) ListThemes(rw http.ResponseWriter, req *http.Request, _ router.Params) error {
	themes, err := m.templates.ListThemes()
	if err != nil {
		return err
	}
	return router.ProtoOK(rw, &api.ListThemeResponse{Themes: themes})
}

func (m *Module) GetTheme(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	name := par.ByName("name")
	theme, ref, err := m.templates.GetTheme(name)
	if err != nil {
		return err
	}
	// todo: fetch remote theme info if available
	if theme == nil {
		return router.ErrNotFound
	}
	var refP *string
	if ref != "" {
		refP = &ref
	}
	return router.ProtoOK(rw, &api.GetThemeResponse{
		Theme: theme,
		Ref:   refP,
	})
}

func (m *Module) GetTemplate(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	name := par.ByName("name")
	templateName := par.ByName("template")
	template, err := m.templates.GetTemplate(name, templateName)
	if err != nil {
		// todo: better errors from m.templates
		if strings.Contains(err.Error(), "not found") {
			return router.ErrNotFound
		}
		return err
	}
	// todo: fetch remote template info if available
	if template == nil {
		return router.ErrNotFound
	}
	// todo: wrap in {"template": ...}?
	return router.ProtoOK(rw, template)
}

func (m *Module) ThemeRegistry(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	r, err := m.templates.Registry()
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
	p, err := m.templates.SearchRegistry(r.GetName())
	if err != nil {
		return errors.New("error searching registry: %s", err)
	}
	if p == nil || p.GetVcsUrl() != r.GetVcsUrl() {
		return router.ErrNotFound
	}

	m.Logger.Infof("cloning package %s from %s", p.GetName(), p.GetVcsUrl())

	// install the package
	return m.templates.InstallThemeFromPackage(p)
}

func (m *Module) CheckThemeForUpdate(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	name := par.ByName("name")
	_, oldRef, currentRef, err := m.templates.CheckThemeForUpdate(name)
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
	return m.templates.UpdateTheme(r.GetName(), r.GetRef())
}
