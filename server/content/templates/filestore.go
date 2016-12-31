package templates

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/nagax/util/token"

	"github.com/octavore/press/proto/press/models"
)

const (
	configFileName       = "theme.json"
	fileStoreTemplateDir = "templates"
	fileStoreAssetsDir   = "assets"
)

// FileStore stores and loads templates on the filesystem
type FileStore struct {
	dataDir string
}

func NewFileStore(dataDir string) *FileStore {
	return &FileStore{dataDir: dataDir}
}

// GetTemplate fetches a theme's template from the filesystem. The
// template's Engine is inferred from the extension in templateName
func (f *FileStore) GetTemplate(theme *models.Theme, templateName string) (*models.ThemeTemplate, error) {
	if theme == nil || theme.GetName() == "" {
		return nil, nil
	}
	p := path.Join(f.dataDir, theme.GetName(), fileStoreTemplateDir, templateName)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	data := string(b)
	ext := strings.TrimLeft(path.Ext(templateName), ".")
	t := &models.ThemeTemplate{
		Theme:  theme.Name,
		Data:   &data,
		Name:   &templateName,
		Engine: &ext,
	}
	return t, nil
}

// GetAsset fetches an asset from the filesystem
func (f *FileStore) GetAsset(theme *models.Theme, assetName string) (*models.ThemeAsset, error) {
	if theme == nil || theme.GetName() == "" {
		return nil, nil
	}
	p := path.Join(f.dataDir, theme.GetName(), fileStoreAssetsDir, assetName)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	data := string(b)
	t := &models.ThemeAsset{
		Theme: theme.Name,
		Data:  &data,
		Name:  &assetName,
	}
	return t, nil
}

// Get a theme from the file store
func (f *FileStore) Get(themeName string) (*models.Theme, error) {
	if themeName == "" {
		return nil, nil
	}
	themeConfigPath := path.Join(f.dataDir, themeName, configFileName)
	_, err := os.Stat(themeConfigPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(themeConfigPath)
	if err != nil {
		return nil, err
	}
	t := &models.Theme{
		Templates: map[string]*models.ThemeTemplate{},
		Assets:    map[string]*models.ThemeAsset{},
	}
	err = jsonpb.Unmarshal(bytes.NewBuffer(b), t)
	if err != nil {
		return nil, err
	}

	// get templates (todo: supported subdirs)
	glob := path.Join(f.dataDir, themeName, fileStoreTemplateDir, "*")
	paths, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		q := path.Base(p)
		e := strings.TrimLeft(path.Ext(p), ".")
		if t.Templates[q] == nil {
			t.Templates[q] = &models.ThemeTemplate{}
		}
		t.Templates[q].Name = &q
		t.Templates[q].Engine = &e
	}

	// get assets (todo: supported subdirs)
	glob = path.Join(f.dataDir, themeName, fileStoreAssetsDir, "*")
	paths, err = filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		q := path.Base(p)
		t.Assets[q] = &models.ThemeAsset{Name: &q}
	}
	return t, nil
}

func themeNameFromPath(p string) string {
	dir, _ := path.Split(p)
	return path.Base(dir)
}

// List all themes in the store
func (f *FileStore) List() ([]*models.Theme, error) {
	glob := path.Join(f.dataDir, "*", configFileName)
	paths, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	themes := []*models.Theme{}
	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		theme := &models.Theme{}
		err = jsonpb.Unmarshal(bytes.NewBuffer(b), theme)
		if err != nil {
			return nil, err
		}
		dir := themeNameFromPath(p)
		theme.Name = &dir
		themes = append(themes, theme)
	}
	return themes, nil
}

// Add a new theme to the store
func (f *FileStore) Add(data io.Reader) (*models.Theme, error) {
	tmpDir := token.New32()
	pth, err := f.extract(data, tmpDir)
	if err != nil {
		return nil, err
	}
	theme, err := f.Get(pth)
	if err != nil {
		return nil, err
	}
	return theme, nil
}

// extract tar.gz from reader into dir
func (f *FileStore) extract(r io.Reader, dir string) (string, error) {
	templateDir := path.Join(f.dataDir, dir)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		// ignore links
		if hdr.Mode == tar.TypeSymlink {
			continue
		}

		p := path.Clean(path.Join(templateDir, hdr.Name))
		if strings.HasPrefix(p, "..") {
			// no relative paths
			continue
		}
		err = os.MkdirAll(path.Dir(p), os.FileMode(0700))
		if err != nil {
			return "", err
		}

		f, err := os.Create(p)
		if err != nil {
			return "", err
		}

		defer f.Close()
		_, err = io.Copy(f, tr)
		if err != nil {
			return "", err
		}
	}
	return templateDir, nil
}
