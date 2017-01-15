package templates

import (
	"io"

	"github.com/octavore/press/proto/press/models"
	"github.com/octavore/press/server/content/engines"
	"github.com/octavore/press/util/errors"
)

// getTemplate returns the desired template. If the theme or template
// does not exist, and an error is returned.
func (m *Module) getTemplate(themeName string, template string) (*models.ThemeTemplate, error) {
	store, theme, err := m.getTheme(themeName)
	if err != nil {
		return nil, err
	}
	if theme == nil {
		return nil, errors.New("content: theme %q not found", themeName)
	}

	tmpl, err := store.GetTemplate(theme, template)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, errors.New("content: template %q not found for theme %q", template, theme.GetName())
	}
	return tmpl, nil
}

// Render a page using the theme and template specified therein.
func (m *Module) Render(w io.Writer, page *models.Page, contents map[string]interface{}) error {
	tmpl, err := m.getTemplate(page.GetTheme(), page.GetTemplate())
	if err != nil {
		return err
	}
	return engines.Render(w, tmpl, contents)
}
