package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

var routeDefaultPage = &models.Page{
	Uuid:        proto.String("page123"),
	PublishedAt: proto.Int64(1494650000),
	Theme:       proto.String("none"),
	Template:    proto.String("markdown"),
	Contents: []*models.Content{{
		Key:   proto.String("content"),
		Value: proto.String("string 1234"),
		Type: &models.Content_Short{
			Short: &models.ContentString{
				Title: proto.String("content"),
				Type:  models.ContentTextType_html.Enum(),
			},
		},
	}},
}

func TestListRoutes(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Routes["route1"] = &models.Route{
		Uuid:   proto.String("route1"),
		Path:   proto.String("/foo1"),
		Target: &models.Route_PageUuid{PageUuid: "page123"},
	}
	te.db.Routes["route2"] = &models.Route{
		Uuid:   proto.String("route2"),
		Path:   proto.String("/foo2"),
		Target: &models.Route_PageUuid{PageUuid: "page456"},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/routes", nil)
	err := te.module.ListRoutes(rw, req, []httprouter.Param{})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{
			"routes": [{
				"uuid": "route1",
				"path": "/foo1",
				"pageUuid": "page123"
			}, {
				"uuid": "route2",
				"path": "/foo2",
				"pageUuid": "page456"
			}]
		}`, rw.Body.String())
	}
}

func TestListRoutesByPage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Routes["route1"] = &models.Route{
		Uuid:   proto.String("route1"),
		Path:   proto.String("/foo1"),
		Target: &models.Route_PageUuid{PageUuid: "page123"},
	}
	te.db.Routes["route2"] = &models.Route{
		Uuid:   proto.String("route2"),
		Path:   proto.String("/foo2"),
		Target: &models.Route_PageUuid{PageUuid: "page456"},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/routes", nil)
	err := te.module.ListRoutesByPage(rw, req, []httprouter.Param{{
		Key: "uuid", Value: "page123",
	}})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{
			"routes": [{
				"uuid": "route1",
				"path": "/foo1",
				"pageUuid": "page123"
			}]
		}`, rw.Body.String())
	}
}

func TestFormatRoute(t *testing.T) {
	r := &models.Route{Path: proto.String(" !fOo@BaR& ")}
	formatRoute(r)
	assert.Equal(t, "/foo-bar", r.GetPath())

	r = &models.Route{Path: proto.String("../foo/!bar/../baz")}
	formatRoute(r)
	assert.Equal(t, "/foo/bar/baz", r.GetPath())
}

func TestUpdateRoute(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages[routeDefaultPage.GetUuid()] = routeDefaultPage

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/routes/",
		bytes.NewBufferString(`{
			"uuid": "route123",
			"path": "foobar",
			"pageUuid": "page123"
		}`),
	)
	err := te.module.UpdateRoute(rw, req, []httprouter.Param{})
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, &models.Route{
			Uuid:   proto.String("route123"),
			Path:   proto.String("/foobar"),
			Target: &models.Route_PageUuid{PageUuid: "page123"},
		}, te.db.Routes["route123"])
	}

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/foobar", nil)
	te.module.Content.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Contains(t, rw.Body.String(), "string 1234")
}

func TestUpdateRoutesByPage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages[routeDefaultPage.GetUuid()] = routeDefaultPage
	te.db.Routes["route1"] = &models.Route{
		Uuid:   proto.String("route1"),
		Path:   proto.String("/foobar"),
		Target: &models.Route_PageUuid{PageUuid: "page123"},
	}
	te.db.Routes["routex"] = &models.Route{
		Uuid:   proto.String("routex"),
		Path:   proto.String("/foobar"),
		Target: &models.Route_PageUuid{PageUuid: "page123"},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/pages/page123/routes",
		bytes.NewBufferString(`{
			"routes": [{
				"uuid": "route1",
				"path": "foo1",
				"pageUuid": "page123"
			}, {
				"uuid": "route2",
				"path": "foo2",
				"pageUuid": "page123"
			}, {
				"uuid": "route3",
				"path": "foo3",
				"pageUuid": "page456"
			}, {
				"path": "foo4",
				"pageUuid": "page123"
			}]
		}`),
	)
	err := te.module.UpdateRoutesByPage(rw, req, []httprouter.Param{{
		Key: "uuid", Value: "page123",
	}})
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, map[string]*models.Route{
			"route1": &models.Route{
				Uuid:   proto.String("route1"),
				Path:   proto.String("/foo1"),
				Target: &models.Route_PageUuid{PageUuid: "page123"},
			},
			"route2": &models.Route{
				Uuid:   proto.String("route2"),
				Path:   proto.String("/foo2"),
				Target: &models.Route_PageUuid{PageUuid: "page123"},
			},
			"route3": &models.Route{
				Uuid:   proto.String("route3"),
				Path:   proto.String("/foo3"),
				Target: &models.Route_PageUuid{PageUuid: "page123"},
			},
			"0": &models.Route{
				Uuid:   proto.String("0"),
				Path:   proto.String("/foo4"),
				Target: &models.Route_PageUuid{PageUuid: "page123"},
			},
		}, te.db.Routes)
	}

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/foo4", nil)
	te.module.Content.ServeHTTP(rw, req)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "string 1234")
	}
}

// DeleteRoute
// UpdateRoutesByPage
