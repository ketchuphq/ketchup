package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"

	"github.com/octavore/press/proto/press/models"
)

const PAGE_BUCKET = "pages"

func (m *Module) GetPage(uuid string) (*models.Page, error) {
	page := &models.Page{}
	err := m.Get(PAGE_BUCKET, uuid, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (m *Module) UpdatePage(page *models.Page) error {
	if page.GetUuid() == "" {
		page.Uuid = proto.String(uuid.NewV4().String())
	}
	return m.Update(PAGE_BUCKET, page)
}

func (m *Module) ListPages() ([]*models.Page, error) {
	lst := []*models.Page{}
	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PAGE_BUCKET))
		return b.ForEach(func(_, v []byte) error {
			pb := &models.Page{}
			err := proto.Unmarshal(v, pb)
			if err != nil {
				return err
			}
			lst = append(lst, pb)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return lst, nil
}
