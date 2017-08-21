package templates

import (
	"io"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/context"
	"github.com/ketchuphq/ketchup/server/content/engines"
	_ "github.com/ketchuphq/ketchup/server/content/engines/html" // load html engine
)

// Render a page using the theme and template specified therein.
func (m *Module) RenderPage(w io.Writer, page *models.Page, context *context.EngineContext) error {
	theme, err := m.getTheme(page.GetTheme())
	if err != nil {
		return err
	}
	return engines.Render(w, theme, page.GetTemplate(), context)
}
