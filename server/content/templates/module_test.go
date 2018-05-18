package templates

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/server/content/engines/html"
	"github.com/ketchuphq/ketchup/server/content/templates/defaultstore"
	"github.com/ketchuphq/ketchup/server/content/templates/dummystore"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

var testTemplate = &models.ThemeTemplate{
	Name:   proto.String("test-template"),
	Engine: proto.String(html.EngineTypeHTML),
	Data:   proto.String(`<div>{{.Page.Content}}</div>`),
	Placeholders: []*models.ThemePlaceholder{
		{
			Key: proto.String("content"),
			Type: &models.ThemePlaceholder_Text{
				Text: &models.ContentText{
					Type: models.ContentTextType_markdown.Enum(),
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
}

func setup(useMemStore bool, themes ...*models.Theme) (*Module, func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{}")
	}))

	configFile, err := ioutil.TempFile("", "ketchup-templates-test-")
	if err != nil {
		panic(err)
	}
	configFile.WriteString(fmt.Sprintf(`{"themes":{"registry_url":"%s"}}`, ts.URL))
	// config := fmt.Sprintf(`{"themes":{"registry_url":"%s"}}`, ts.URL)
	// m.ConfigModule.Byte = []byte(config)
	m := &Module{}
	svc := service.New(m)
	m.ConfigModule.TestConfigPath = configFile.Name()
	stop := svc.StartForTest()

	if useMemStore {
		m.themeStore = dummy.New()
		m.Stores = []store.ThemeStore{
			&defaultstore.DefaultStore{},
			m.themeStore,
			m.internalStore,
		}
	}

	for _, theme := range themes {
		err := m.themeStore.Add(theme)
		if err != nil {
			panic(err)
		}
	}
	return m, func() {
		stop()
		ts.Close()
		configFile.Close()
	}
}
