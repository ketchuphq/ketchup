package api

import (
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"

	"github.com/octavore/press/db"
	"github.com/octavore/press/db/bolt"
	"github.com/octavore/press/proto/press/api"
	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/router"
)

func (m *Module) getPage(par httprouter.Params, fn func(*models.Page) error) (*models.Page, error) {
	uuid := par.ByName("uuid")
	if uuid == "" {
		return nil, router.ErrNotFound
	}
	page, err := m.DB.GetPage(uuid)
	if _, ok := err.(bolt.ErrNoKey); ok {
		return nil, router.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	err = fn(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (m *Module) GetPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page, err := m.getPage(par, func(*models.Page) error { return nil })
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) ListPages(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	pages, err := m.DB.ListPages()
	if err != nil {
		return err
	}
	db.SortPagesByUpdatedAt(pages, false)
	return router.Proto(rw, &api.ListPageResponse{
		Pages: pages,
	})
}

func (m *Module) UpdatePage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page := &models.Page{}
	err := jsonpb.Unmarshal(req.Body, page)
	if err != nil {
		return err
	}
	err = m.DB.UpdatePage(page)
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) PublishPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page, err := m.getPage(par, func(page *models.Page) error {
		now := time.Now().Unix()
		page.PublishedAt = &now
		err := m.DB.UpdatePage(page)
		if err != nil {
			return err
		}
		return m.Content.ReloadRouter()
	})
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) UnpublishPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page, err := m.getPage(par, func(page *models.Page) error {
		page.PublishedAt = nil
		err := m.DB.UpdatePage(page)
		if err != nil {
			return err
		}
		return m.Content.ReloadRouter()
	})
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) DeletePage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	uuid := par.ByName("uuid")
	page, err := m.DB.GetPage(uuid)
	if _, ok := err.(bolt.ErrNoKey); ok {
		return router.ErrNotFound
	}
	if err != nil {
		return err
	}
	return m.DB.DeletePage(page)
}
