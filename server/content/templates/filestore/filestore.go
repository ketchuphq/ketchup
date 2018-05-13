package filestore

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/plugins/pkg"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

const (
	configFileName       = "theme.json"
	fileStoreTemplateDir = "templates"
	fileStoreAssetsDir   = "assets"
)

var jpb = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	Indent:       "  ",
	OrigName:     false,
}

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

// updateThemeDirMap iterates over all folders in the base dir and reads all the
// theme configs found. also updates the mapping of folder name to theme name,
// if different.
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
		c, err := readConfig(themeConfigPath)
		if err != nil {
			// todo: log error
			return errors.Wrap(err)
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
func (f *FileStore) getTemplate(theme *models.Theme, templateName string) (*models.ThemeTemplate, error) {
	if theme == nil || theme.GetName() == "" {
		return nil, nil
	}
	themeDir := theme.GetName()
	if altDir := f.themeDirMap[theme.GetName()]; altDir != "" {
		themeDir = altDir
	}
	p := path.Join(f.baseDir, themeDir, fileStoreTemplateDir, templateName)
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}
	data := string(b)
	t := proto.Clone(theme.GetTemplates()[templateName]).(*models.ThemeTemplate)
	t.Theme = theme.Name
	t.Data = &data
	t.Name = &templateName

	if t.GetEngine() == "" {
		ext := strings.TrimLeft(path.Ext(templateName), ".")
		t.Engine = &ext
	}
	return t, nil
}

// GetAsset fetches an asset from the filesystem
func (f *FileStore) getAsset(theme *models.Theme, assetName string) (*models.ThemeAsset, error) {
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

// Get a theme from the file store. Template and asset data (i.e. the file
// contents) are NOT loaded.
func (f *FileStore) Get(themeName string) (store.Theme, error) {
	themeDir := themeName
	if altDir := f.themeDirMap[themeName]; altDir != "" {
		themeDir = altDir
	}

	themeConfigPath := path.Join(f.baseDir, themeDir, configFileName)
	t, err := readConfig(themeConfigPath)
	if err != nil || t == nil {
		return nil, nil
	}

	// get templates (todo: supported subdirs)
	baseTemplateDir := path.Clean(path.Join(f.baseDir, themeDir, fileStoreTemplateDir))
	err = filepath.Walk(baseTemplateDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		p = strings.TrimPrefix(path.Clean(p), baseTemplateDir)
		p = strings.TrimLeft(p, "/")
		if strings.HasPrefix(path.Base(p), ".") {
			return nil
		}
		if t.Templates[p] == nil {
			t.Templates[p] = &models.ThemeTemplate{}
		}
		t.Templates[p].Name = &p

		// infer engine if not specified
		if t.Templates[p].GetEngine() == "" {
			e := strings.TrimLeft(path.Ext(p), ".")
			t.Templates[p].Engine = &e
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	baseAssetDir := path.Clean(path.Join(f.baseDir, themeDir, fileStoreAssetsDir))
	err = filepath.Walk(baseAssetDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err)
		}
		if info.IsDir() {
			return nil
		}
		p = strings.TrimPrefix(path.Clean(p), baseAssetDir)
		p = strings.TrimLeft(p, "/")
		if strings.HasPrefix(path.Base(p), ".") {
			return nil
		}
		t.Assets[p] = &models.ThemeAsset{Name: &p}
		return nil
	})

	if err != nil {
		return nil, err
	}

	absThemeDir := path.Join(f.baseDir, themeDir)
	ok, err := isVCS(absThemeDir)
	if err != nil {
		return nil, errors.Wrap(err)
	} else if ok {
		currentRef, err := getCurrentRef(absThemeDir)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		return &Theme{Theme: t, store: f, ref: currentRef}, nil
	}

	return &Theme{Theme: t, store: f}, nil
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

// AddPackage adds a theme from a theme file by cloning it from
// the VCS location to the themeDir.
func (f *FileStore) AddPackage(p *packages.Package) error {
	themeDir := path.Join(f.baseDir, p.GetName())
	return pkg.CloneToDir(themeDir, p.GetVcsUrl())
}

// Add a theme directly to the themeDir.
func (f *FileStore) Add(theme *models.Theme) error {
	theme = proto.Clone(theme).(*models.Theme)
	templateDir := path.Join(f.baseDir, theme.GetName())
	perm := os.FileMode(0600)

	err := themeIterator(theme, func(fn string, el themeFile) error {
		p := path.Clean(path.Join(templateDir, fn))
		if strings.HasPrefix(p, "..") {
			return nil
		}

		err := os.MkdirAll(path.Dir(p), os.ModePerm)
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
		return err
	}

	fw, err := os.Create(path.Join(templateDir, configFileName))
	if err != nil {
		return err
	}
	return jpb.Marshal(fw, theme)
}

func (f *FileStore) UpdateThemeToRef(themeName, commitHash string) error {
	themeDir := themeName
	if altDir := f.themeDirMap[themeName]; altDir != "" {
		themeDir = altDir
	}
	repoDir := path.Join(f.baseDir, themeDir)
	return pkg.FetchDir(repoDir, commitHash)
}

func (f *FileStore) GetAsset(assetName string) (*models.ThemeAsset, error) {
	lst, err := ioutil.ReadDir(f.baseDir)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	for _, fi := range lst {
		if !fi.IsDir() {
			continue
		}
		p := path.Join(f.baseDir, fi.Name(), fileStoreAssetsDir, assetName)
		b, err := ioutil.ReadFile(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, errors.Wrap(err)
		}
		data := string(b)
		t := &models.ThemeAsset{
			// Theme: theme.Name,
			Data: &data,
			Name: &assetName,
		}
		return t, nil
	}

	return nil, nil
}
