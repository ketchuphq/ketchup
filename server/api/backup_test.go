package api

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/ketchuphq/ketchup/db/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Routes = fixtures.Routes
	te.db.Pages = fixtures.Pages

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/data/test-key", nil)
	err := te.module.GetBackup(rw, req, nil)
	fmt.Println(rw.Body.String())
	assert.NoError(t, err)
	assert.JSONEq(t, fixtures.JSON, rw.Body.String())
}
