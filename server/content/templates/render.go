package templates

import (
	"fmt"
	"io"

	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/content/engines"
)

// Render a page using the theme and template specified therein.
func (m *Module) Render(w io.Writer, page *models.Page, contents map[string]interface{}) error {
	store, theme, err := m.getTheme(page.GetTheme())
	if err != nil {
		return err
	}
	if theme == nil {
		return fmt.Errorf("content: theme %q not found", page.GetTheme())
	}

	tmpl, err := store.GetTemplate(theme, page.GetTemplate())
	if err != nil {
		return err
	}
	if tmpl == nil {
		return fmt.Errorf("content: template %q not found for theme %q", page.GetTemplate(), page.GetTheme())
	}
	return engines.Render(w, tmpl, contents)
}
