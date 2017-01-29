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

	"github.com/octavore/ketchup/proto/ketchup/models"
)

const (
	configFileName       = "theme.json"
	fileStoreTemplateDir = "templates"
	fileStoreAssetsDir   = "assets"
)

// FileStore stores and loads templates on the filesystem
type FileStore struct {
	dataDir     string
	themeDirMap map[string]string // maps theme name to dir
}

func (f *FileStore) readConfig(themeConfigPath string) (*models.Theme, error) {
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
	return t, nil
}

func (f *FileStore) updateThemeDirMap() error {
	lst, err := ioutil.ReadDir(f.dataDir)
	if err != nil {
		return err
	}
	m := map[string]string{}
	for _, fi := range lst {
		if !fi.IsDir() {
			continue
		}
		themeConfigPath := path.Join(f.dataDir, fi.Name(), configFileName)
		c, err := f.readConfig(themeConfigPath)
		if err != nil {
			return nil
		}
		if c.GetName() != "" && fi.Name() != c.GetName() {
			m[c.GetName()] = fi.Name()
		}
	}
	f.themeDirMap = m
	return nil
}

// GetTemplate fetches a theme's template from the filesystem. The
// template's Engine is inferred from the extension in templateName
func (f *FileStore) GetTemplate(theme *models.Theme, templateName string) (*models.ThemeTemplate, error) {
	if theme == nil || theme.GetName() == "" {
		return nil, nil
	}
	themeDir := theme.GetName()
	if altDir := f.themeDirMap[theme.GetName()]; altDir != "" {
		themeDir = altDir
	}
	p := path.Join(f.dataDir, themeDir, fileStoreTemplateDir, templateName)
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
	themeDir := theme.GetName()
	if altDir := f.themeDirMap[theme.GetName()]; altDir != "" {
		themeDir = altDir
	}

	p := path.Join(f.dataDir, themeDir, fileStoreAssetsDir, assetName)
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

func (f *FileStore) loadTheme(themeDir string, t *models.Theme) error {
	// get templates (todo: supported subdirs)
	glob := path.Join(f.dataDir, themeDir, fileStoreTemplateDir, "*")
	paths, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, p := range paths {
		q := path.Base(p)
		if strings.HasPrefix(q, ".") {
			continue
		}
		e := strings.TrimLeft(path.Ext(p), ".")
		if t.Templates[q] == nil {
			t.Templates[q] = &models.ThemeTemplate{}
		}
		t.Templates[q].Name = &q
		t.Templates[q].Engine = &e
	}

	// get assets (todo: supported subdirs)
	glob = path.Join(f.dataDir, themeDir, fileStoreAssetsDir, "*")
	paths, err = filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, p := range paths {
		q := path.Base(p)
		if strings.HasPrefix(q, ".") {
			continue
		}
		t.Assets[q] = &models.ThemeAsset{Name: &q}
	}
	return nil
}

// Get a theme from the file store
func (f *FileStore) Get(themeName string) (*models.Theme, error) {
	themeDir := themeName
	if altDir := f.themeDirMap[themeName]; altDir != "" {
		themeDir = altDir
	}
	themeConfigPath := path.Join(f.dataDir, themeDir, configFileName)
	t, err := f.readConfig(themeConfigPath)
	if err != nil || t == nil {
		return nil, nil
	}

	// get templates (todo: supported subdirs)
	glob := path.Join(f.dataDir, themeDir, fileStoreTemplateDir, "*")
	paths, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		q := path.Base(p)
		if strings.HasPrefix(q, ".") {
			continue
		}
		e := strings.TrimLeft(path.Ext(p), ".")
		if t.Templates[q] == nil {
			t.Templates[q] = &models.ThemeTemplate{}
		}
		t.Templates[q].Name = &q
		t.Templates[q].Engine = &e
	}

	// get assets (todo: supported subdirs)
	glob = path.Join(f.dataDir, themeDir, fileStoreAssetsDir, "*")
	paths, err = filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		q := path.Base(p)
		if strings.HasPrefix(q, ".") {
			continue
		}
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
		if theme.GetName() == "" {
			dir := themeNameFromPath(p)
			theme.Name = &dir
		}
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
