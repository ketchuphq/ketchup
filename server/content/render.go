package content

import (
	"io"

	"github.com/octavore/press/proto/press/models"
)

func (m *Module) render(w io.Writer, page *models.Page) error {
	contents, err := renderPageContents(page)
	if err != nil {
		return err
	}
	return m.Templates.Render(w, page, contents)
}
