package api

import (
	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/server/content/templates"
	"github.com/ketchuphq/ketchup/server/tls"
	"github.com/stretchr/testify/mock"
	"github.com/xenolf/lego/acme"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/content/engines/html"
)

var testTemplate = &models.ThemeTemplate{
	Name:   proto.String("test-template"),
	Engine: proto.String(html.EngineTypeHTML),
	Data:   proto.String(`<div>{{.Page.Content}}</div>`),
	Placeholders: []*models.ThemePlaceholder{
		{
			Key: proto.String("bPlaceholder"),
			Type: &models.ThemePlaceholder_Text{
				Text: &models.ContentText{
					Title: proto.String("Template Placeholder"),
					Type:  models.ContentTextType_markdown.Enum(),
				},
			},
		},
	},
}

var testTheme = &models.Theme{
	Name: proto.String("test-theme"),
	Package: &packages.Package{
		VcsUrl: proto.String("https://localhost:8000/foo.git"),
	},
	Templates: map[string]*models.ThemeTemplate{
		"test-template": testTemplate,
	},
	Assets: map[string]*models.ThemeAsset{
		"app.js": {
			Name: proto.String("app.js"),
			Data: proto.String("var foo = 1;"),
		},
	},
	Placeholders: []*models.Data{
		{
			Key: proto.String("aPlaceholder"),
			Type: &models.Data_Short{
				Short: &models.ContentString{
					Title: proto.String("Theme Placeholder"),
					Type:  models.ContentTextType_markdown.Enum(),
				},
			},
		},
	},
}

type testTLSModule struct {
	mock.Mock
}

func (m *testTLSModule) GetAllRegisteredDomains() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *testTLSModule) GetRegistration(domain string, withPrivateKey bool) (*tls.Registration, error) {
	args := m.Called(domain, withPrivateKey)
	return args.Get(0).(*tls.Registration), args.Error(1)
}

func (m *testTLSModule) LoadCertResource(domain string) (*acme.CertificateResource, error) {
	args := m.Called(domain)
	return args.Get(0).(*acme.CertificateResource), args.Error(1)
}

func (m *testTLSModule) ObtainCert(email, domain string) error {
	args := m.Called(email, domain)
	return args.Error(0)
}

type testTemplateModule struct {
	mock.Mock
	*templates.Module
}

func (m *testTemplateModule) CheckThemeForUpdate(name string) (bool, string, string, error) {
	args := m.Called(name)
	return args.Bool(0), args.String(1), args.String(2), args.Error(3)
}

func (m *testTemplateModule) GetTemplate(theme, template string) (*models.ThemeTemplate, error) {
	return m.Module.GetTemplate(theme, template)
}

func (m *testTemplateModule) GetTheme(name string) (*models.Theme, string, error) {
	return m.Module.GetTheme(name)
}

func (m *testTemplateModule) InstallThemeFromPackage(p *packages.Package) error {
	m.Called(p)
	return nil
}

func (m *testTemplateModule) ListThemes() ([]*models.Theme, error) {
	m.Called()
	return m.Module.ListThemes()
}

func (m *testTemplateModule) Registry() (*packages.Registry, error) {
	m.Called()
	args := m.Called()
	return args.Get(0).(*packages.Registry), args.Error(1)
}

func (m *testTemplateModule) SearchRegistry(themeName string) (*packages.Package, error) {
	args := m.Called(themeName)
	return args.Get(0).(*packages.Package), args.Error(1)
}

func (m *testTemplateModule) UpdateTheme(name, ref string) error {
	m.Called(name, ref)
	return nil
}
