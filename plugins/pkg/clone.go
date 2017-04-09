package pkg

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/octavore/ketchup/proto/ketchup/models"
)

// Clone the given repo to {data_dir}/{dir}/{name}
func (m *Module) Clone(packageName, packageURL, dir string) error {
	packagePath := m.Config.DataPath(path.Join(dir, packageName), "")
	return m.clone(packagePath, packageURL)
}

func getRepoMasterIterator(repo *git.Repository) (*object.FileIter, error) {
	ref, err := repo.Reference("refs/remotes/origin/master", true)
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	iter, err := commit.Files()
	if err != nil {
		return nil, err
	}
	// defer iter.Close()

	return iter, nil
}

type setDatar interface {
	SetData(v *string)
}

// CloneToTheme fetches the theme to memory returns the model and theme
func (m *Module) CloneToTheme(url string) (*models.Theme, error) {
	memStorage := memory.NewStorage()
	r, err := git.Clone(memStorage, nil, &git.CloneOptions{URL: url})
	if err != nil {
		return nil, err
	}
	iter, err := getRepoMasterIterator(r)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	theme := &models.Theme{}
	assets := map[string]*models.ThemeAsset{}
	templates := map[string]*models.ThemeTemplate{}

	err = iter.ForEach(func(f *object.File) error {
		pth := path.Clean(f.Name)
		if pth == "theme.json" {
			rdr, err := f.Reader()
			if err != nil {
				return err
			}
			defer rdr.Close()
			return jsonpb.Unmarshal(rdr, theme)
		}

		var dest setDatar
		if strings.HasPrefix(pth, "assets/") {
			asset := &models.ThemeAsset{
				Name: proto.String(strings.TrimPrefix(pth, "assets/")),
			}
			assets[*asset.Name] = asset
			dest = asset
		} else if strings.HasPrefix(pth, "templates/") {
			template := &models.ThemeTemplate{
				Name: proto.String(strings.TrimPrefix(pth, "templates/")),
			}
			templates[*template.Name] = template
			dest = template
		} else {
			// log unknown
			return nil
		}

		rdr, err := f.Reader()
		if err != nil {
			return err
		}
		defer rdr.Close()

		bytes, err := ioutil.ReadAll(rdr)
		if err != nil {
			return err
		}
		data := string(bytes)
		dest.SetData(&data)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// set data for templates in theme.json
	for templateName, template := range theme.Templates {
		if t, ok := templates[templateName]; ok {
			template.Data = t.Data
			delete(templates, templateName)
		}
	}

	// set data for assets in theme.json
	for assetName, asset := range theme.Assets {
		if a, ok := assets[assetName]; ok {
			asset.Data = a.Data
			delete(assets, assetName)
		}
	}

	// copy templates not in theme.json
	for templateName, template := range templates {
		theme.Templates[templateName] = template
	}

	// copy assets not in theme.json
	for assetName, asset := range assets {
		theme.Assets[assetName] = asset
	}

	return theme, nil
}

// Clone the given url to the given dest
// todo: support cloning into s3/github store?
func (m *Module) clone(dest, url string) error {
	r, err := git.PlainClone(dest, false, &git.CloneOptions{URL: url})
	if err != nil {
		return err
	}

	iter, err := getRepoMasterIterator(r)
	if err != nil {
		return err
	}
	defer iter.Close()

	return iter.ForEach(func(f *object.File) error {
		pth := path.Join(dest, f.Name)
		err := os.MkdirAll(path.Dir(pth), 0700)
		if err != nil {
			return err
		}
		mode, err := f.Mode.ToOSFileMode()
		if err != nil {
			return err
		}
		g, err := os.OpenFile(pth, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
		defer g.Close()
		rdr, err := f.Reader()
		if err != nil {
			return err
		}
		defer rdr.Close()
		_, err = io.Copy(g, rdr)
		return err
	})
}
