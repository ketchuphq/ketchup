package api

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestGetInfo(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Users["user123"] = &models.User{
		Uuid: proto.String("user123"),
	}

	// not logged in
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/user", nil)
	err := te.module.GetInfo(rw, req, nil)
	if assert.NoError(t, err) {
		assert.JSONEq(t, `{
			"registry_url": "http://themes.ketchuphq.com/registry.json",
			"version": ""
		}`, rw.Body.String())
	}
}
