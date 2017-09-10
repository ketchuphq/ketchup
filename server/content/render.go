package content

import (
	"io"

	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/content"
	"github.com/ketchuphq/ketchup/server/content/context"
)

func (m *Module) render(w io.Writer, page *models.Page) error {
	contents, err := m.CreateContentMap(page)
	if err != nil {
		return err
	}
	context := context.NewContext(m.Logger, page, m.DB.Backend, contents)
	return m.Templates.RenderPage(w, page, context)
}

type contentMap map[string]interface{}

func (m *Module) CreateContentMap(page *models.Page) (contentMap, error) {
	contents := map[string]interface{}{}
	contents["title"] = page.GetTitle()
	var err error
	for _, c := range page.Contents {
		contents[c.GetKey()], err = content.RenderContent(c)
		if err != nil {
			if errors.IsType(err, content.ErrUnknownContentType{}) {
				m.Logger.Error(err)
			}
			return nil, err
		}
	}
	return contents, err
}
