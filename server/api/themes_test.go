package api

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/nagax/router"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/content/templates/defaultstore"
	"github.com/ketchuphq/ketchup/server/content/templates/store"
)

func TestListThemes(t *testing.T) {
	te := setup()
	defer te.stop()

	tmplStore := te.module.Templates.Stores[1]
	err := tmplStore.Add(testTheme)
	if !assert.NoError(t, err) {
		t.Fail()
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/themes/", nil)
	err = te.module.ListThemes(rw, req, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.JSONEq(t, `{
		"themes": [
			{
				"name": "none",
				"templates": {
					"html": {
						"name": "html",
						"engine": "html",
						"placeholders": [
							{
								"key": "content",
								"text": {
									"type": "html"
								}
							}
						]
					},
					"markdown": {
						"name": "markdown",
						"engine": "html",
						"placeholders": [
							{
								"key": "content",
								"text": {
									"type": "markdown"
								}
							}
						]
					}
				}
			},
			{
				"name": "test-theme",
				"package": {
					"vcsUrl": "https://localhost:8000/foo.git"
				},
				"templates": {
					"test-template": {
						"name": "test-template",
						"engine": "html",
						"placeholders": [
							{
								"key": "bPlaceholder",
								"text": {
									"title": "Template Placeholder",
									"type": "markdown"
								}
							}
						]
					}
				},
				"assets": {
					"app.js": {
						"name": "app.js"
					}
				},
				"placeholders": [
					{
						"key": "aPlaceholder",
						"short": {
							"title": "Theme Placeholder",
							"type": "markdown"
						}
					}
				]
			}
		]
	}`, rw.Body.String())
}

func TestGetTheme(t *testing.T) {
	te := setup()
	defer te.stop()

	tmplStore := te.module.Templates.Stores[1]
	err := tmplStore.Add(testTheme)
	if !assert.NoError(t, err) {
		t.Fail()
	}

	// test theme
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/themes/"+testTheme.GetName(), nil)
	err = te.module.GetTheme(rw, req, []httprouter.Param{
		{Key: "name", Value: testTheme.GetName()},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	expected := &api.GetThemeResponse{
		Theme: store.StripData(testTheme),
	}
	output := &api.GetThemeResponse{}
	assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
	assert.Equal(t, expected, output)

	// test default theme
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/themes/none", nil)
	err = te.module.GetTheme(rw, req, []httprouter.Param{
		{Key: "name", Value: "none"},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	noneTheme, err := (&defaultstore.DefaultStore{}).Get("none")
	assert.Nil(t, err)
	expected = &api.GetThemeResponse{
		Theme: noneTheme.Proto(),
	}
	output = &api.GetThemeResponse{}
	assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
	assert.Equal(t, expected, output)
}

func TestGetTemplate(t *testing.T) {
	te := setup()
	defer te.stop()

	tmplStore := te.module.Templates.Stores[1]
	err := tmplStore.Add(testTheme)
	if !assert.NoError(t, err) {
		t.Fail()
	}

	// test template
	rw := httptest.NewRecorder()
	path := fmt.Sprintf("/api/v1/themes/%s/%s", testTheme.GetName(), testTemplate.GetName())
	req := httptest.NewRequest("GET", path, nil)
	err = te.module.GetTemplate(rw, req, []httprouter.Param{
		{Key: "name", Value: testTheme.GetName()},
		{Key: "template", Value: testTemplate.GetName()},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	output := &models.ThemeTemplate{}
	assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
	assert.Equal(t, testTheme.GetName(), output.GetTheme())
	output.Theme = nil // remove Theme because it's not in original template object
	assert.Equal(t, testTemplate, output)

	// test non existent template
	rw = httptest.NewRecorder()
	path = fmt.Sprintf("/api/v1/themes/%s/%s", testTheme.GetName(), "foo")
	req = httptest.NewRequest("GET", path, nil)
	err = te.module.GetTemplate(rw, req, []httprouter.Param{
		{Key: "name", Value: testTheme.GetName()},
		{Key: "template", Value: "foo"},
	})
	assert.Equal(t, router.ErrNotFound, err)

	// test default template
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/themes/none/markdown", nil)
	err = te.module.GetTemplate(rw, req, []httprouter.Param{
		{Key: "name", Value: "none"},
		{Key: "template", Value: "markdown"},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	noneTheme, err := (&defaultstore.DefaultStore{}).Get("none")
	assert.Nil(t, err)
	output = &models.ThemeTemplate{}
	assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
	noneMarkdownTemplate, err := noneTheme.GetTemplate("markdown")
	assert.NoError(t, err)
	assert.Equal(t, noneMarkdownTemplate, output)
}

func TestThemeRegistry(t *testing.T) {
	te := setup()
	defer te.stop()
	expectedRegistry := &packages.Registry{
		Packages: []*packages.Package{{Name: proto.String("my package")}},
	}
	mockTemplates := &testTemplateModule{Module: te.module.Templates}
	mockTemplates.On("Registry").Return(expectedRegistry, nil)
	te.module.templates = mockTemplates

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/theme-registry", nil)
	err := te.module.ThemeRegistry(rw, req, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	mockTemplates.AssertExpectations(t)
	assert.JSONEq(t, `{"packages": [{"name":"my package"}]}`, rw.Body.String())
	output := &packages.Registry{}
	assert.NoError(t, jsonpb.Unmarshal(rw.Body, output))
	assert.Equal(t, expectedRegistry, output)
}

func TestInstallTheme(t *testing.T) {
	te := setup()
	defer te.stop()

	// test successful install
	testPkg := &packages.Package{Name: proto.String("cow theme"), VcsUrl: proto.String("http://localhost:8000")}
	mockTemplates := &testTemplateModule{Module: te.module.Templates}
	mockTemplates.On("SearchRegistry", testPkg.GetName()).Return(testPkg, nil)
	mockTemplates.On("InstallThemeFromPackage", testPkg).Return()
	te.module.templates = mockTemplates
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/theme-install", bytes.NewBufferString(`{
		"name": "cow theme",
		"vcs_url": "http://localhost:8000"
	}`))
	err := te.module.InstallTheme(rw, req, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "", rw.Body.String())
	mockTemplates.AssertExpectations(t)

	// test no name
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/theme-install", bytes.NewBufferString(`{}`))
	err = te.module.InstallTheme(rw, req, nil)
	assert.EqualError(t, err, "Theme name is required.")

	// test vcs url mismatch
	mockTemplates = &testTemplateModule{Module: te.module.Templates}
	mockTemplates.On("SearchRegistry", testPkg.GetName()).Return(testPkg, nil)
	te.module.templates = mockTemplates
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/theme-install", bytes.NewBufferString(`{
		"name": "cow theme"
	}`))
	err = te.module.InstallTheme(rw, req, nil)
	assert.Equal(t, err, router.ErrNotFound)
}

func TestCheckThemeForUpdate(t *testing.T) {
	te := setup()
	defer te.stop()

	themeName := "myTheme"
	mockTemplates := &testTemplateModule{Module: te.module.Templates}
	mockTemplates.On("CheckThemeForUpdate", themeName).Return(true, "123", "456", nil)
	te.module.templates = mockTemplates

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/themes/myTheme/updates", nil)
	err := te.module.CheckThemeForUpdate(rw, req, httprouter.Params{
		{Key: "name", Value: themeName},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.JSONEq(t, `{"oldRef":"123","currentRef":"456"}`, rw.Body.String())
	mockTemplates.AssertExpectations(t)
}

func TestUpdateTheme(t *testing.T) {
	te := setup()
	defer te.stop()

	themeName := "myTheme"
	mockTemplates := &testTemplateModule{Module: te.module.Templates}
	mockTemplates.On("UpdateTheme", themeName, "789").Return(nil)
	te.module.templates = mockTemplates

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/themes/myTheme/update", bytes.NewBufferString(`{
		"name": "myTheme",
		"ref": "789"
	}`))
	err := te.module.UpdateTheme(rw, req, httprouter.Params{
		{Key: "name", Value: themeName},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	mockTemplates.AssertExpectations(t)

	// test name mismatch
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/themes/myTheme/update", bytes.NewBufferString(`{
		"name": "wrongName",
		"ref": "789"
	}`))
	err = te.module.UpdateTheme(rw, req, httprouter.Params{
		{Key: "name", Value: themeName},
	})
	assert.EqualError(t, err, "theme name mismatch")
}
