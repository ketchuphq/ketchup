package users

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/users"
	"github.com/octavore/nagax/users/databaseauth"
	"github.com/octavore/nagax/users/session"

	"github.com/octavore/ketchup/db"
	"github.com/octavore/ketchup/server/config"
)

// Module users is a largely a wrapper around nagax/users/databaseauth
// to enable user logins.
type Module struct {
	Auth     *users.Module
	DBAuth   *databaseauth.Module // todo: make this pluggable?
	Sessions *session.Module
	Config   *config.Module
	Router   *router.Module
	DB       *db.Module
	Logger   *logger.Module
}

// Handle is similar to httprouter.Handle, except it returns an error which can be
// handled separately.
type Handle func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.AddCommand(registerSetPassword(m))
	c.AddCommand(registerUserAdd(m))

	c.Setup = func() error {
		m.DBAuth.Configure(
			databaseauth.WithUserStore(&userStore{m.DB}),
			databaseauth.WithLoginPath("/api/v1/login"),
			databaseauth.WithErrorHandler(m.ErrorHandler),
		)
		m.Sessions.KeyFile = m.Config.DataPath("session.key", "session.key")
		m.Auth.ErrorHandler = m.ErrorHandler
		m.Auth.RegisterAuthenticator(m.Sessions)
		// todo: allow to be set in a config file
		m.Sessions.CookieDomain = ""
		return nil
	}
}

// ErrorHandler to handle auth errors.
// todo: return more fine-grained error codes.
func (m *Module) ErrorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	m.Router.SimpleError(rw, http.StatusInternalServerError, err)
}

// MustWithAuth wraps a Handle function and requires a logged in user.
// The new handler will authenticate using the auth module, creating
// another wrapper for the original function to make it compatible with
// http.HandlerFunc (while ensuring it closes over httprouter.Params).
func (m *Module) MustWithAuth(delegate Handle) Handle {
	return func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
		var err error
		h := func(x http.ResponseWriter, y *http.Request) {
			err = delegate(x, y, par)
		}
		m.Auth.MustWithAuth(h)(rw, req)
		return err
	}
}
