package dummy

import "github.com/ketchuphq/ketchup/proto/ketchup/models"

type Theme struct {
	*models.Theme
	ThemeRef string
}

func (t *Theme) Ref() (string, bool) {
	if t.ThemeRef == "" {
		return "", false
	}
	return t.ThemeRef, true
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
