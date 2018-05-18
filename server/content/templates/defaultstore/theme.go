package defaultstore

import "github.com/ketchuphq/ketchup/proto/ketchup/models"

type Theme struct {
}

func (t *Theme) Ref() (string, bool) {
	return "", false
}

func (t *Theme) Proto() *models.Theme {
	return noneThemeNoData
}

func (t *Theme) GetTemplate(templateName string) (*models.ThemeTemplate, error) {
	return noneTheme.Templates[templateName], nil
}

func (t *Theme) GetAsset(assetName string) (*models.ThemeAsset, error) {
	return nil, nil
}
