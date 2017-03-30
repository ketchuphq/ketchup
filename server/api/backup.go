package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (m *Module) GetBackup(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	exp, err := m.DB.Export()
	if err != nil {
		return err
	}

	date := time.Now().Format("20060102-030405")

	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s.bak", date))
	return m.Router.JSON(rw, http.StatusOK, exp)
}
