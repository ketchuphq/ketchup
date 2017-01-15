package users

import (
	"fmt"
	"net/http"

	"github.com/howeyc/gopass"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/users"
	"github.com/octavore/nagax/users/databaseauth"
	"github.com/octavore/nagax/users/session"
	"github.com/octavore/nagax/util/token"

	"github.com/octavore/press/db"
	"github.com/octavore/press/server/config"
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
	c.AddCommand(&service.Command{
		Keyword:    "users:add <email>",
		Usage:      "Add a new user.",
		ShortUsage: "Add a new user",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			email := ctx.Args[0]
			fmt.Println("enter a password:")
			pass, err := gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
			_, err = m.DBAuth.Create(email, string(pass))
			if err != nil {
				panic(err)
			}
		},
	})

	c.AddCommand(&service.Command{
		Keyword:    "users:password <email>",
		Usage:      "Set user password.",
		ShortUsage: "Set user password",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			email := ctx.Args[0]
			u, err := m.DB.GetUserByEmail(email)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("enter a password:")
			pass, err := gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
			hashedPass := databaseauth.HashPassword(string(pass), token.New32())
			u.SetHashedPassword(&hashedPass)
			err = m.DB.UpdateUser(u)
			if err != nil {
				panic(err)
			}
		},
	})
	c.Setup = func() error {
		m.DBAuth.Configure(
			databaseauth.WithUserStore(m),
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

// MustWithAuth wraps a http
func (m *Module) MustWithAuth(delegate Handle) Handle {
	return func(rw http.ResponseWriter, req *http.Request, par httprouter.Params) (err error) {
		// authenticate with the auth module, with a wrapper for the delegate to
		// make it compatible with http.HandlerFunc (while ensuring it closes over httprouter.Params).
		h := func(x http.ResponseWriter, y *http.Request) {
			err = delegate(x, y, par)
		}
		m.Auth.MustWithAuth(h)(rw, req)
		return
	}
}
