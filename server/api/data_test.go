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
	"github.com/ketchuphq/ketchup/server/router"
)

func TestGetData(t *testing.T) {
	te := setup()
	te.db.Data["test-key"] = &models.Data{
		Key:   proto.String("test-key"),
		Value: proto.String("test-value"),
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/data/test-key", nil)
	err := te.module.GetData(rw, req, []httprouter.Param{
		{Key: "key", Value: "test-key"},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{
			"key": "test-key",
			"value": "test-value"
		}`, rw.Body.String())
	}
}

func TestGetData_NotFound(t *testing.T) {
	te := setup()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/data/test-key", nil)
	err := te.module.GetData(rw, req, []httprouter.Param{
		{Key: "key", Value: "test-key"},
	})
	assert.EqualError(t, err, router.ErrNotFound.Error())
}

func TestSetData(t *testing.T) {
	// new data
	te := setup()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/data",
		bytes.NewBufferString(`{
			"data": [{
				"key": "new-key",
				"value": "1"
			}]
		}`))

	err := te.module.UpdateData(rw, req, httprouter.Params{})
	if !assert.NoError(t, err) {
		t.Fail()
	}
	assert.Equal(t, http.StatusOK, rw.Code)

	// update data
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/data",
		bytes.NewBufferString(`{
			"data": [{
				"key": "new-key",
				"value": "2"
			}]
		}`))
	err = te.module.UpdateData(rw, req, httprouter.Params{})
	if !assert.NoError(t, err) {
		t.Fail()
	}
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, &models.Data{
		Key:   proto.String("new-key"),
		Value: proto.String("2"),
		Type: &models.Data_Short{
			Short: &models.ContentString{
				Type: models.ContentTextType_text.Enum(),
			},
		},
	}, te.db.Data["new-key"])
}
