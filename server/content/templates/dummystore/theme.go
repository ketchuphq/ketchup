package dummy

import "github.com/ketchuphq/ketchup/proto/ketchup/models"

type Theme struct {
	*models.Theme
	ref string
}

func (t *Theme) Render(templateName string) (string, error) {
	return "", nil
}

func (t *Theme) Ref() (string, bool) {
	if t.ref == "" {
		return "", false
	}
	return t.ref, true
}

func (t *Theme) Proto() *models.Theme {
	return t.Theme
}

func (t *Theme) GetTemplate(templateName string) (*models.ThemeTemplate, error) {
	return t.Theme.Templates[templateName], nil
}

func (t *Theme) GetAsset(assetName string) (*models.ThemeAsset, error) {
	return nil, nil
}
