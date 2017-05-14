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

func TestGetPage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages["page123"] = &models.Page{
		Uuid: proto.String("page123"),
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/pages/page123", nil)
	err := te.module.GetPage(rw, req, []httprouter.Param{
		{Key: "uuid", Value: "page123"},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{
			"uuid": "page123"
		}`, rw.Body.String())
	}
}

func TestGetRenderedPage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages["page123"] = &models.Page{
		Uuid: proto.String("page123"),
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/pages/page123", nil)
	err := te.module.GetRenderedPage(rw, req, []httprouter.Param{
		{Key: "uuid", Value: "page123"},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{
			"title": ""
		}`, rw.Body.String())
	}
}

func TestListPages(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages["page456"] = &models.Page{
		Uuid: proto.String("page456"),
		Timestamps: &models.Timestamp{
			UpdatedAt: proto.Int64(1494650100),
		},
	}
	te.db.Pages["page123"] = &models.Page{
		Uuid: proto.String("page123"),
		Timestamps: &models.Timestamp{
			UpdatedAt: proto.Int64(1494650000),
		},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/pages/", nil)
	err := te.module.ListPages(rw, req, []httprouter.Param{})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.JSONEq(t, `{
			"pages": [{
				"uuid": "page456",
				"timestamps": {
					"updatedAt": "1494650100"
				}
			}, {
				"uuid": "page123",
				"timestamps": {
					"updatedAt": "1494650000"
				}
			}]
		}`, rw.Body.String())
	}
}

func TestUpdatePage(t *testing.T) {
	te := setup()
	defer te.stop()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/pages/",
		bytes.NewBufferString(`{
			"uuid": "1234",
			"title": "hello world",
			"contents": [{
				"key": "content",
				"value": "hello world",
				"short": {
					"title": "content",
					"type": "html"
				}
			}, {
				"key": "ignored_content"
			}]
		}`),
	)

	err := te.module.UpdatePage(rw, req, []httprouter.Param{})
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Equal(t, &models.Page{
			Uuid:  proto.String("1234"),
			Title: proto.String("hello world"),
			Contents: []*models.Content{{
				Key:   proto.String("content"),
				Value: proto.String("hello world"),
				Type: &models.Content_Short{
					Short: &models.ContentString{
						Title: proto.String("content"),
						Type:  models.ContentTextType_html.Enum(),
					},
				},
			}},
		}, te.db.Pages["1234"])
	}
}

func TestPublishPage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages["page123"] = &models.Page{
		Uuid: proto.String("page123"),
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/pages/page123", nil)
	err := te.module.PublishPage(rw, req, []httprouter.Param{
		{Key: "uuid", Value: "page123"},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.NotNil(t, te.db.Pages["page123"].PublishedAt)
	}
}

func TestUnpublishPage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages["page123"] = &models.Page{
		Uuid:        proto.String("page123"),
		PublishedAt: proto.Int64(1494650000),
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/pages/page123", nil)
	err := te.module.UnpublishPage(rw, req, []httprouter.Param{
		{Key: "uuid", Value: "page123"},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Nil(t, te.db.Pages["page123"].PublishedAt)
	}
}

func TestDeletePage(t *testing.T) {
	te := setup()
	defer te.stop()
	te.db.Pages["page123"] = &models.Page{
		Uuid: proto.String("page123"),
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/v1/pages/test-key", nil)
	err := te.module.DeletePage(rw, req, []httprouter.Param{
		{Key: "uuid", Value: "page123"},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Nil(t, te.db.Pages["page123"])
	}
}
