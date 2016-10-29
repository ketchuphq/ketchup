package bolt

import (
	"github.com/golang/protobuf/proto"
	"github.com/octavore/press/proto/press/models"
	"github.com/satori/go.uuid"
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
	if page.Uuid == nil {
		page.Uuid = proto.String(uuid.NewV4().String())
	}
	return m.Update(PAGE_BUCKET, page)
}
