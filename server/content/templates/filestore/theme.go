package filestore

import "github.com/ketchuphq/ketchup/proto/ketchup/models"

type Theme struct {
	*models.Theme
	store *FileStore
	ref   string
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
	return t.store.getTemplate(t.Theme, templateName)
}

func (t *Theme) GetAsset(assetName string) (*models.ThemeAsset, error) {
	return t.store.getAsset(t.Theme, assetName)
}
