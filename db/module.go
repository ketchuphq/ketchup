package db

import (
	"io"

	"github.com/octavore/naga/service"

	"github.com/octavore/press/proto/press/models"
)

// Backend interface for models
type Backend interface {
	GetUser(email string) (*models.User, error)

	GetPage(uuid string) (*models.Page, error)
	UpdatePage(page *models.Page) error
	ListPages() ([]*models.Page, error)

	GetRoute(key string) (*models.Route, error)
	UpdateRoute(route *models.Route) error
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
