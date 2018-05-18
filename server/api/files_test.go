package api

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/nagax/router"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func TestGetFile(t *testing.T) {
	te := setup()
	defer te.stop()
	file, err := te.module.Files.Upload("app.js", bytes.NewBufferString(`let a = 'b';`))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/files/"+file.GetUuid(), nil)
	err = te.module.GetFile(rw, req, []httprouter.Param{
		{Key: "uuid", Value: file.GetUuid()},
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	expected := &api.FileResponse{File: file}
	output := &api.FileResponse{}
	assert.NoError(t, jsonpb.UnmarshalString(rw.Body.String(), output))
	assert.Equal(t, expected, output)
}

func TestGetFile__NotFound(t *testing.T) {
	te := setup()
	defer te.stop()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/files/1234", nil)
	err := te.module.GetFile(rw, req, []httprouter.Param{
		{Key: "uuid", Value: "1234"},
	})
	assert.EqualError(t, err, router.ErrNotFound.Error())

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/files/1234", nil)
	err = te.module.GetFile(rw, req, []httprouter.Param{
		{Key: "uuid", Value: ""},
	})
	assert.Equal(t, err, router.ErrNotFound)
}

func TestDeleteFile(t *testing.T) {
	te := setup()
	defer te.stop()
	file, err := te.module.Files.Upload("app.js", bytes.NewBufferString(`let a = 'b';`))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	// make delete request
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/api/v1/files/"+file.GetUuid(), nil)
	err = te.module.DeleteFile(rw, req, []httprouter.Param{
		{Key: "uuid", Value: file.GetUuid()},
	})
	assert.NoError(t, err)

	// assert file was deleted from db
	assert.Nil(t, te.db.Files["app.js"])

	// assert file was deleted from store
	f, err := te.module.Files.Get("app.js")
	assert.NoError(t, err)
	assert.Nil(t, f)

	// make delete request when already deleted
	rw = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "/api/v1/files/"+file.GetUuid(), nil)
	err = te.module.DeleteFile(rw, req, []httprouter.Param{
		{Key: "uuid", Value: file.GetUuid()},
	})
	assert.NoError(t, err)
}

func TestListFiles(t *testing.T) {
	te := setup()
	defer te.stop()

	files := map[string]string{
		"app.js":  "let a = true;",
		"app.css": "body { height: 100%; }",
	}
	for k, v := range files {
		_, err := te.module.Files.Upload(k, bytes.NewBufferString(v))
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/files/", nil)
	err := te.module.ListFiles(rw, req, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	output := &api.ListFilesResponse{}
	jsonpb.UnmarshalString(rw.Body.String(), output)
	expected := &api.ListFilesResponse{
		Files: []*models.File{{
			Uuid: output.Files[0].Uuid,
			Name: proto.String("app.css"),
			Url:  proto.String("/static/app.css"),
		}, {
			Uuid: output.Files[1].Uuid,
			Name: proto.String("app.js"),
			Url:  proto.String("/static/app.js"),
		}},
	}
	assert.Equal(t, expected, output)
}

func TestUploadFile(t *testing.T) {
	te := setup()
	defer te.stop()
	content := `let a = 1;`

	// create multipart object
	buf := &bytes.Buffer{}
	mp := multipart.NewWriter(buf)
	wr, err := mp.CreateFormFile("file", "app.js")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	_, err = wr.Write([]byte(content))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.NoError(t, mp.Close()) {
		t.FailNow()
	}

	// make upload request
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/files/", buf)
	req.Header.Set("content-type", "multipart/form-data; boundary="+mp.Boundary())
	err = te.module.UploadFile(rw, req, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	// check return code

	// check return value
	output := &api.FileResponse{}
	assert.NoError(t, jsonpb.UnmarshalString(rw.Body.String(), output))
	expected := &api.FileResponse{
		File: &models.File{
			Uuid: output.GetFile().Uuid,
			Name: proto.String("app.js"),
			Url:  proto.String("/static/app.js"),
		},
	}
	assert.Equal(t, expected, output)

	// check database updated
	uploadedFile := te.db.Files[output.GetFile().GetUuid()]
	assert.Equal(t, expected.File, uploadedFile)

	// check file contents uploaded
	r, err := te.module.Files.Get("app.js")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, []byte(content), data)
}
