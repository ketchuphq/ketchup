package db

import (
	"fmt"
	"io"

	"github.com/octavore/naga/service"

	"github.com/octavore/ketchup/proto/ketchup/models"
	"github.com/octavore/ketchup/util/errors"
)

// Backend interface for models
type Backend interface {
	GetUser(uuid string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(*models.User) error

	GetPage(uuid string) (*models.Page, error)
	UpdatePage(*models.Page) error
	DeletePage(*models.Page) error
	ListPages() ([]*models.Page, error)

	GetRoute(uuid string) (*models.Route, error)
	UpdateRoute(*models.Route) error
	DeleteRoute(*models.Route) error
	ListRoutes() ([]*models.Route, error)

	Debug(w io.Writer) error // print debug info
}

type Module struct {
	Backend
}

func (m *Module) Init(c *service.Config) {
	c.Start = func() {
		if m.Backend == nil {
			panic("backend not configured")
		}
	}
}

func (m *Module) Register(b Backend) error {
	if m.Backend != nil {
		return errors.Wrap(fmt.Errorf("backend already configured"))
	}
	m.Backend = b
	return nil
}
