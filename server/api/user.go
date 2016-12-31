package api

import (
	"net/http"

	"github.com/octavore/nagax/users"
	"github.com/octavore/press/server/router"

	"github.com/julienschmidt/httprouter"
)

func (m *Module) GetUser(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	userUUID, ok := req.Context().Value(users.UserTokenKey{}).(string)
	if !ok {
		m.Router.EmptyJSON(rw, http.StatusNotFound)
		return nil
	}
	user, err := m.DB.GetUser(userUUID)
	if err != nil {
		return err
	}
	user.HashedPassword = nil
	return router.Proto(rw, user)
}

func (m *Module) Logout(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) error {
	m.Auth.Auth.Logout(rw, req)
	return nil
}
