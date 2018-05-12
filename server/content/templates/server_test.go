package templates

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	m, stop := setup(false, testTheme)
	defer stop()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(rw, req)
	assert.Equal(t, 404, rw.Code)
	assert.Equal(t, "404 page not found\n", rw.Body.String())

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/app.js", nil)

	m.ServeHTTP(rw, req)
	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, testTheme.Assets["app.js"].GetData(), rw.Body.String())
}

var notFoundTheme = &models.Theme{
	Name: proto.String("not-found-theme"),
	Assets: map[string]*models.ThemeAsset{
		"404.html": {
			Name: proto.String("404.html"),
			Data: proto.String("double rainbow!"),
		},
	},
}

func TestNotFound(t *testing.T) {
	m, stop := setup(false, notFoundTheme)
	defer stop()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(rw, req)
	assert.Equal(t, 404, rw.Code)
	assert.Equal(t, "double rainbow!", rw.Body.String())

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/foo", nil)

	m.NotFound(rw, req)
	assert.Equal(t, 404, rw.Code)
	assert.Equal(t, "double rainbow!", rw.Body.String())
}
