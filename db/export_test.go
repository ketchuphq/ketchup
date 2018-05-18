package db

import (
	"io/ioutil"
	"testing"

	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/db/dummy"
	"github.com/ketchuphq/ketchup/db/fixtures"
)

func TestImportFromJSON(t *testing.T) {
	m := &Module{}
	svc := service.New(m)
	backend := dummy.New()
	m.Register(backend)
	stop := svc.StartForTest()
	defer stop()

	f, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	defer f.Close()
	f.WriteString(fixtures.JSON)
	err = m.importFromJSON(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, fixtures.Pages, backend.Pages)
	assert.Equal(t, fixtures.Routes, backend.Routes)
}

func TestExportToJSON(t *testing.T) {
	m := &Module{}
	svc := service.New(m)
	backend := dummy.New()
	backend.Pages = fixtures.Pages
	backend.Routes = fixtures.Routes
	m.Register(backend)
	stop := svc.StartForTest()
	defer stop()

	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	outputPath := d + "data.json"
	err = m.ExportToJSON(outputPath)
	assert.NoError(t, err)

	output, err := ioutil.ReadFile(outputPath)
	assert.NoError(t, err)
	assert.JSONEq(t, fixtures.JSON, string(output))

	m.Backend = dummy.New()
	err = m.importFromJSON(outputPath)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.Pages, backend.Pages)
	assert.Equal(t, fixtures.Routes, backend.Routes)
}
