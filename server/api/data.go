package api

import (
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/router"
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/db/bolt"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

var defaultPlaceholders = []*models.Data{
	{
		Key: proto.String("title"),
		Type: &models.Data_Short{
			Short: &models.ContentString{
				Type: models.ContentTextType_text.Enum(),
			},
		},
	},
}

func (m *Module) addAllPlaceholders(data []*models.Data) ([]*models.Data, error) {
	keys := map[string]*models.Data{}
	for _, d := range data {
		keys[d.GetKey()] = d
	}

	for _, p := range defaultPlaceholders {
		if _, ok := keys[p.GetKey()]; !ok {
			data = append(data, proto.Clone(p).(*models.Data))
		}
	}

	themes, err := m.Templates.ListThemes()
	if err != nil {
		return nil, err
	}

	for _, theme := range themes {
		for _, p := range theme.GetPlaceholders() {
			key := fmt.Sprintf("%s.%s", theme.GetName(), p.GetKey())
			if _, ok := keys[key]; !ok {
				keys[key] = proto.Clone(p).(*models.Data)
				data = append(data, keys[key])
			}
		}
	}
	return data, nil
}

func (m *Module) ListData(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	data, err := m.DB.ListData()
	if err != nil {
		return err
	}
	data, err = m.addAllPlaceholders(data)
	if err != nil {
		return err
	}
	return router.ProtoOK(rw, &api.ListDataResponse{
		Data: data,
	})
}

func (m *Module) GetData(rw http.ResponseWriter, req *http.Request, par router.Params) error {
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
	return router.ProtoOK(rw, data)
}

func (m *Module) DeleteData(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	key := par.ByName("key")
	if key == "" {
		return router.ErrNotFound
	}
	return m.DB.DeleteData(&models.Data{Key: &key})
}

func (m *Module) UpdateData(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	rpb := &api.UpdateDataRequest{}
	err := jsonpb.Unmarshal(req.Body, rpb)
	if err != nil {
		return errors.Wrap(err)
	}

	for _, data := range rpb.Data {
		if data.Type == nil {
			data.Type = &models.Data_Short{
				Short: &models.ContentString{
					Type: models.ContentTextType_text.Enum(),
				},
			}
		}
	}
	return m.DB.UpdateDataBatch(rpb.Data)
}
