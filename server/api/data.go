package api

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"

	"github.com/ketchuphq/ketchup/db/bolt"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/router"
)

func (m *Module) ListData(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	data, err := m.DB.ListData()
	if _, ok := err.(bolt.ErrNoKey); ok {
		return router.ErrNotFound
	}
	if err != nil {
		return err
	}
	return router.Proto(rw, &api.ListDataResponse{
		Data: data,
	})
}

func (m *Module) GetData(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	key := par.ByName("key")
	if key == "" {
		return router.ErrNotFound
	}
	data, err := m.DB.GetData(key)
	if _, ok := err.(bolt.ErrNoKey); ok || data == nil {
		return router.ErrNotFound
	}
	if err != nil {
		return err
	}
	return router.Proto(rw, data)
}

func (m *Module) UpdateData(rw http.ResponseWriter, req *http.Request, par httprouter.Params) error {
	data := &models.Data{}
	err := jsonpb.Unmarshal(req.Body, data)
	if err != nil {
		return err
	}
	err = m.DB.UpdateData(data)
	if err != nil {
		return err
	}
	return router.Proto(rw, data)
}
