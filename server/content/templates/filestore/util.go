package filestore

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/nagax/util/errors"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

type themeFile interface {
	SetData(*string)
	GetData() string
}

// readConfig unmarshals a theme.json file to a model.Theme.
// if there is no file, return an error.
func readConfig(themeConfigPath string) (*models.Theme, error) {
	_, err := os.Stat(themeConfigPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}
	b, err := ioutil.ReadFile(themeConfigPath)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	t := &models.Theme{
		Templates: map[string]*models.ThemeTemplate{},
		Assets:    map[string]*models.ThemeAsset{},
	}
	err = jsonpb.Unmarshal(bytes.NewBuffer(b), t)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return t, nil
}

// themeIterator iterates over theme files stored in a models.Theme struct
func themeIterator(theme *models.Theme, iterFn func(fn string, el themeFile) error) error {
	for fn, tmpl := range theme.Templates {
		err := iterFn(path.Join(fileStoreTemplateDir, fn), tmpl)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	for fn, asset := range theme.Assets {
		err := iterFn(path.Join(fileStoreAssetsDir, fn), asset)
		if err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}

// getLatestRef returns the latest ref for repo at repoPath
func getLatestRef(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", errors.Wrap(err)
	}
	head, err := repo.Head()
	if err != nil {
		return "", errors.Wrap(err)
	}
	return head.Hash().String(), nil
}

func themeNameFromPath(p string) string {
	dir, _ := path.Split(p)
	return path.Base(dir)
}
