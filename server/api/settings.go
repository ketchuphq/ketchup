package api

import (
	"net/http"

	"github.com/octavore/nagax/router"

	"github.com/ketchuphq/ketchup/server/version"
)

func (m *Module) GetInfo(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	return router.JSON(rw, http.StatusOK, map[string]string{
		"version":      version.Get(),
		"registry_url": m.templates.GetRegistryURL(),
	})
}
