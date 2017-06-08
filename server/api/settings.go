package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ketchuphq/ketchup/server/version"
)

func (m *Module) GetInfo(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	return m.Router.JSON(rw, http.StatusOK, map[string]string{
		"version":      version.Get(),
		"registry_url": m.Templates.GetRegistryURL(),
	})
}
