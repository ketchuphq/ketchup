package api

import (
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"

	"github.com/octavore/ketchup/db"
	"github.com/octavore/ketchup/db/bolt"
	"github.com/octavore/ketchup/proto/ketchup/api"
	"github.com/octavore/ketchup/proto/ketchup/models"
	"github.com/octavore/ketchup/server/router"
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

// GetPage gets a page by UUID.
// todo: nest response?
func (m *Module) GetPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page, err := m.getPage(par, func(*models.Page) error { return nil })
	if err != nil {
		return err
	}
	return router.Proto(rw, page)
}

func (m *Module) GetRenderedPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page, err := m.getPage(par, func(*models.Page) error { return nil })
	if err != nil {
		return err
	}

	contents, err := m.Content.CreateContentMap(page)
	if err != nil {
		return err
	}

	return m.Router.JSON(rw, http.StatusOK, contents)
}

// ListPages returns all pages, sorted by updated at.
// todo: pagination, filtering
// todo: error handling?
func (m *Module) ListPages(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	opts := &api.ListPageRequest{}
	err := req.ParseForm()
	if err != nil {
		return err
	}
	err = m.decoder.Decode(opts, req.Form)
	if err != nil {
		return err
	}
	pages, err := m.DB.ListPages(opts)
	if err != nil {
		return err
	}
	db.SortPagesByUpdatedAt(pages, false)
	return router.Proto(rw, &api.ListPageResponse{
		Pages: pages,
	})
}

// UpdatePage saves the given page to the DB.
// todo: nest response?
func (m *Module) UpdatePage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page := &models.Page{}
	// use jsonpb.unmarshal to correct unmarshal int64 e.g. PublishedAt
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

// PublishPage sets the published time on a page to the current time.
func (m *Module) PublishPage(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	page, err := m.getPage(par, func(page *models.Page) error {
		// already published
		if page.PublishedAt != nil {
			return nil
		}

		// set published at to current time
		nowMillis := time.Now().UnixNano() / 1e6
		page.PublishedAt = &nowMillis
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

// UnpublishPage sets published at to null, effectively unpublishing the page.
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

// DeletePage deletes the given page.
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
