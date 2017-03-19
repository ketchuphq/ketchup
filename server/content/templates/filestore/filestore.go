package filestore

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
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/octavore/ketchup/proto/ketchup/models"
	"github.com/octavore/ketchup/util/errors"
)

const (
	configFileName       = "theme.json"
	fileStoreTemplateDir = "templates"
	fileStoreAssetsDir   = "assets"
)

// FileStore stores and loads templates on the filesystem
type FileStore struct {
	baseDir     string
	themeDirMap map[string]string // maps theme name to dir
}

// New returns a new file store which updates periodically
func New(baseDir string, updateInterval time.Duration, log func(args ...interface{})) (*FileStore, error) {
	f := &FileStore{baseDir: baseDir}
	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		return nil, err
	}
	err = f.updateThemeDirMap()
	if err != nil {
		return nil, err
	}
	go func() {
		for range time.Tick(updateInterval) {
			err := f.updateThemeDirMap()
			if err != nil {
				log(err)
			}
		}
	}()
	return f, nil
}

func (f *FileStore) readConfig(themeConfigPath string) (*models.Theme, error) {
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

func (f *FileStore) updateThemeDirMap() error {
	lst, err := ioutil.ReadDir(f.baseDir)
	if err != nil {
		return errors.Wrap(err)
	}
	m := map[string]string{}
	for _, fi := range lst {
		if !fi.IsDir() {
			continue
		}
		themeConfigPath := path.Join(f.baseDir, fi.Name(), configFileName)
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
	p := path.Join(f.baseDir, themeDir, fileStoreTemplateDir, templateName)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	data := string(b)
	ext := strings.TrimLeft(path.Ext(templateName), ".")
	t := proto.Clone(theme.GetTemplates()[templateName]).(*models.ThemeTemplate)
	t.Theme = theme.Name
	t.Data = &data
	t.Name = &templateName
	t.Engine = &ext
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

	p := path.Join(f.baseDir, themeDir, fileStoreAssetsDir, assetName)
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

func (f *FileStore) loadTheme(t *models.Theme) error {
	// get templates (todo: supported subdirs)
	glob := path.Join(f.baseDir, t.GetName(), fileStoreTemplateDir, "*")
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
	glob = path.Join(f.baseDir, t.GetName(), fileStoreAssetsDir, "*")
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
	themeConfigPath := path.Join(f.baseDir, themeDir, configFileName)
	t, err := f.readConfig(themeConfigPath)
	if err != nil || t == nil {
		return nil, nil
	}

	// get templates (todo: supported subdirs)
	glob := path.Join(f.baseDir, themeDir, fileStoreTemplateDir, "*")
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
	glob = path.Join(f.baseDir, themeDir, fileStoreAssetsDir, "*")
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
	glob := path.Join(f.baseDir, "*", configFileName)
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

type themeFile interface {
	SetData(*string)
	GetData() string
}

func themeIterator(theme *models.Theme, iterFn func(fn string, el themeFile) error) error {
	for fn, tmpl := range theme.Templates {
		err := iterFn(fn, tmpl)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	for fn, asset := range theme.Assets {
		err := iterFn(fn, asset)
		if err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}

var jpb = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	Indent:       "  ",
	OrigName:     false,
}

// Add a theme from a theme file.
func (f *FileStore) Add(theme *models.Theme) error {
	theme = proto.Clone(theme).(*models.Theme)
	templateDir := path.Join(f.baseDir, theme.GetName())
	perm := os.FileMode(0700)

	err := themeIterator(theme, func(fn string, el themeFile) error {
		p := path.Clean(path.Join(templateDir, fn))
		if strings.HasPrefix(p, "..") {
			return nil
		}

		err := os.MkdirAll(path.Dir(p), perm)
		if err != nil {
			return errors.Wrap(err)
		}

		err = ioutil.WriteFile(p, []byte(el.GetData()), perm)
		if err != nil {
			return errors.Wrap(err)
		}
		el.SetData(nil)
		return nil
	})

	if err != nil {
		return nil
	}

	fw, err := os.Create(path.Join(templateDir, configFileName))
	if err != nil {
		return err
	}
	return jpb.Marshal(fw, theme)
}

// extract tar.gz from reader into dir
func (f *FileStore) extract(r io.Reader, dir string) (string, error) {
	templateDir := path.Join(f.baseDir, dir)
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
