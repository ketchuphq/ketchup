package db

import (
	"io"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

// Backend interface for models.
// todo: split update/create?
// todo: filter page list?
// for SQL dbs, how to determine if we need to upgrade database/run migration?
// https://github.com/mgutz/dat
type Backend interface {
	// GetUser looks up a user by uuid, and returns nil, nil if no user.
	GetUser(uuid string) (*models.User, error)

	// GetUserByEmail TODO
	GetUserByEmail(email string) (*models.User, error)

	// GetUserByToken looks up a user by token, and returns nil, nil if no user.
	GetUserByToken(token string) (*models.User, error)

	// UpdateUser
	UpdateUser(*models.User) error

	// GetPage looks up a page by uuid, and returns nil, nil if no page.
	GetPage(uuid string) (*models.Page, error)

	// UpdatePage updates (creating if necessary) a new page.
	// New pages and new content blocks will be assigned UUIDs
	UpdatePage(*models.Page) error

	// DeletePage deletes the page and should also delete
	// related routes.
	DeletePage(*models.Page) error

	// ListPages lists all the existing pages. May be unsorted.
	ListPages(*api.ListPageRequest) ([]*models.Page, error)

	// GetRoute fetches a route from the DB by UUID
	GetRoute(uuid string) (*models.Route, error)

	// UpdateRoute updates (or creates if necessary) an existing route.
	// The route uuid is set if it is blank.
	UpdateRoute(*models.Route) error

	// DeleteRoute deletes a route from the DB (but not any related pages).
	DeleteRoute(*models.Route) error

	// ListRoutes returns all existing routes. May be unsorted.
	ListRoutes() ([]*models.Route, error)

	GetData(key string) (*models.Data, error)

	UpdateData(*models.Data) error

	DeleteData(*models.Data) error

	ListData() ([]*models.Data, error)

	Debug(w io.Writer) error // print debug info
}
