package dummy

import (
	"io"

	"github.com/octavore/ketchup/proto/ketchup/api"
	"github.com/octavore/ketchup/proto/ketchup/models"
)

type DummyDB struct {
	Users  map[string]*models.User
	Pages  map[string]*models.Page
	Routes map[string]*models.Route
}

func New() *DummyDB {
	return &DummyDB{
		Users:  map[string]*models.User{},
		Pages:  map[string]*models.Page{},
		Routes: map[string]*models.Route{},
	}
}

func (d *DummyDB) GetUser(uuid string) (*models.User, error) {
	return d.Users[uuid], nil
}

func (d *DummyDB) GetUserByEmail(email string) (*models.User, error) {
	for _, u := range d.Users {
		if u.GetEmail() == email {
			return u, nil
		}
	}
	return nil, nil
}

func (d *DummyDB) GetUserByToken(token string) (*models.User, error) {
	for _, u := range d.Users {
		if u.GetToken() == token {
			return u, nil
		}
	}
	return nil, nil
}

func (d *DummyDB) UpdateUser(u *models.User) error {
	// todo: set uuid, timestamp?
	d.Users[u.GetEmail()] = u
	return nil
}

func (d *DummyDB) GetPage(uuid string) (*models.Page, error) {
	return d.Pages[uuid], nil
}

func (d *DummyDB) UpdatePage(p *models.Page) error {
	d.Pages[p.GetUuid()] = p
	return nil
}

func (d *DummyDB) DeletePage(p *models.Page) error {
	delete(d.Pages, p.GetUuid())
	return nil
}

func (d *DummyDB) ListPages(_ *api.ListPageRequest) ([]*models.Page, error) {
	pages := []*models.Page{}
	for _, p := range d.Pages {
		pages = append(pages, p)
	}
	return pages, nil
}

func (d *DummyDB) GetRoute(uuid string) (*models.Route, error) {
	return d.Routes[uuid], nil
}

func (d *DummyDB) UpdateRoute(r *models.Route) error {
	d.Routes[r.GetUuid()] = r
	return nil
}

func (d *DummyDB) DeleteRoute(r *models.Route) error {
	delete(d.Routes, r.GetUuid())
	return nil
}

func (d *DummyDB) ListRoutes() ([]*models.Route, error) {
	routes := []*models.Route{}
	for _, r := range d.Routes {
		routes = append(routes, r)
	}
	return routes, nil
}

func (d *DummyDB) Debug(w io.Writer) error {
	return nil
}
