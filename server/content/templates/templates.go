package templates

import (
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/util/errors"
)

// getTemplate returns the desired template. If the theme or template
// does not exist, and an error is returned.
// themeName (e.g. 'san-marzano')
// template (e.g. 'index.html')
func (m *Module) getTemplate(themeName string, template string) (*models.ThemeTemplate, error) {
	theme, err := m.getTheme(themeName)
	if err != nil {
		return nil, err
	}
	if theme == nil {
		return nil, errors.New("content: theme %q not found", themeName)
	}

	tmpl, err := theme.GetTemplate(template)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, errors.New("content: template %q not found for theme %q", template, theme.Proto().GetName())
	}
	return tmpl, nil
}

// GetTemplate returns the template for the given theme and template.
func (m *Module) GetTemplate(theme, template string) (*models.ThemeTemplate, error) {
	tmpl, err := m.getTemplate(theme, template)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
