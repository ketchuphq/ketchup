package templates

import (
	"io"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/server/content/engines"
)

// Render a page using the theme and template specified therein.
func (m *Module) Render(w io.Writer, page *models.Page, context *context.EngineContext) error {
	tmpl, err := m.getTemplate(page.GetTheme(), page.GetTemplate())
	if err != nil {
		return err
	}
	return engines.Render(w, tmpl, context)
}
