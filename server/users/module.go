package users

import (
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/users"
	"github.com/octavore/nagax/users/databaseauth"
	"github.com/octavore/nagax/users/session"
	"github.com/octavore/nagax/users/tokenauth"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/server/config"
)

// Module users is a largely a wrapper around nagax/users/databaseauth
// to enable user logins.
type Module struct {
	Auth     *users.Module
	Sessions *session.Module
	Config   *config.Module
	Router   *router.Module
	DB       *db.Module
	Logger   *logger.Module

	// todo: make these pluggable?
	DBAuth    *databaseauth.Module
	TokenAuth *tokenauth.Module
}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.AddCommand(registerSetPassword(m))
	c.AddCommand(registerUserAdd(m))
	c.AddCommand(registerGenerateToken(m))

	c.Setup = func() error {
		m.DBAuth.Configure(
			databaseauth.WithUserStore(&userStore{m.DB}),
			databaseauth.WithLoginPath("/api/v1/login"),
		)
		m.Sessions.KeyFile = m.Config.DataPath("session.key", "session.key")
		// todo: allow to be set in a config file
		m.Sessions.CookieDomain = ""
		m.Auth.RegisterAuthenticator(m.Sessions)

		m.TokenAuth.Configure(
			tokenauth.WithTokenSource(&tokenStore{m.DB, m.Logger}),
		)
		m.Auth.RegisterAuthenticator(m.TokenAuth)
		return nil
	}
}
