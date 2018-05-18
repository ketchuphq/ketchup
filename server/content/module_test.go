package content

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/db/dummy"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	dummytmpl "github.com/ketchuphq/ketchup/server/content/templates/dummystore"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

var page = &models.Page{
	Uuid:     proto.String("1234"),
	Theme:    proto.String("test-theme"),
	Template: proto.String("markdown"),
	Contents: []*models.Content{
		{
			Uuid:  proto.String(""),
			Key:   proto.String("content"),
			Value: proto.String("**hello world**"),
			Type: &models.Content_Text{
				Text: &models.ContentText{
					Title: proto.String("title"),
					Type:  models.ContentTextType_markdown.Enum(),
				},
			},
		},
	},
}

var pageRoute = &models.Route{
	Uuid:   proto.String("2222"),
	Path:   proto.String("/zee"),
	Target: &models.Route_PageUuid{PageUuid: "1234"},
}

func testModule() (*Module, func()) {
	m := &Module{}
	svc := service.New(m)
	m.DB.Register(dummy.New())
	stopper := svc.StartForTest()
	m.Templates.Stores = []store.ThemeStore{dummytmpl.New()}
	return m, stopper
}

func TestNotFound(t *testing.T) {
	m, stopper := testModule()
	defer stopper()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/zee", nil)
	m.ReloadRouter()
	m.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Fatal("unexpected response", rw.Code, rw.Body.String())
	}
}

func TestReloadNew(t *testing.T) {
	m, stopper := testModule()
	defer stopper()
	err := m.ReloadRouter()
	if err != nil {
		t.Fatal(err)
	}
	page := proto.Clone(page).(*models.Page)
	err = m.DB.UpdatePage(page)
	if err != nil {
		t.Fatal(err)
	}
	err = m.DB.UpdateRoute(pageRoute)
	if err != nil {
		t.Fatal(err)
	}

	// not found because not published
	err = m.ReloadRouter()
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/zee", nil)
	m.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Fatal("unexpected response", rw.Code, rw.Body.String())
	}

	// found because published
	page.PublishedAt = proto.Int64(time.Now().UnixNano() / 1e6)
	err = m.DB.UpdatePage(page)
	if err != nil {
		t.Fatal(err)
	}

	theme := &models.Theme{}
	json.Unmarshal([]byte(`{
    "name": "test-theme",
    "templates": {
      "markdown": {
        "name": "markdown",
        "engine": "html",
        "data": "<html>{{.Page.Content}}</html>",
        "placeholders": []
      }
    }
  }`), theme)
	err = m.Templates.Stores[0].Add(theme)
	if err != nil {
		t.Fatal(err)
	}

	// reload published route
	err = m.ReloadRouter()
	if err != nil {
		t.Fatal(err)
	}
	rw = httptest.NewRecorder()
	m.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("unexpected code", rw.Code, rw.Body.String())
	}
	if !strings.Contains(rw.Body.String(), "<strong>hello world</strong>") {
		t.Fatal("unexpected response", rw.Code, rw.Body.String())
	}
	page.PublishedAt = proto.Int64(time.Now().UnixNano() / 1e6)
	m.DB.UpdatePage(page)

	// updating route should work
	route := proto.Clone(pageRoute).(*models.Route)
	route.Path = proto.String("/two")
	err = m.DB.UpdateRoute(route)
	if err != nil {
		t.Fatal(err)
	}

	// updated route: old route does not work
	err = m.ReloadRouter()
	if err != nil {
		t.Fatal(err)
	}
	rw = httptest.NewRecorder()
	m.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Fatal("unexpected response", rw.Code, rw.Body.String())
	}

	// updated route: new route works
	rw = httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/two", nil)
	m.ServeHTTP(rw, req2)
	if rw.Code != http.StatusOK {
		t.Fatal("unexpected response", rw.Code, rw.Body.String())
	}
}
