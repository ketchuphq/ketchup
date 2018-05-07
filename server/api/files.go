package api

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/router"

	"github.com/ketchuphq/ketchup/db/bolt"
	"github.com/ketchuphq/ketchup/proto/ketchup/api"
)

const maxUploadMem = 16 << 20 // 16MiB

func (m *Module) GetFile(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	uuid := par.ByName("uuid")
	if uuid == "" {
		return router.ErrNotFound
	}
	file, err := m.DB.GetFile(uuid)
	if _, ok := err.(bolt.ErrNoKey); ok {
		return router.ErrNotFound
	}
	if err != nil {
		return err
	}
	file.Url = proto.String(m.Files.URLForFile(file))

	return router.ProtoOK(rw, &api.FileResponse{File: file})
}

func (m *Module) DeleteFile(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	uuid := par.ByName("uuid")
	if uuid == "" {
		return router.ErrNotFound
	}
	return m.Files.Delete(uuid)
}

func (m *Module) ListFiles(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	files, err := m.Files.DB.ListFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		file.Url = proto.String(m.Files.URLForFile(file))
	}

	return router.ProtoOK(rw, &api.ListFilesResponse{Files: files})
}

func (m *Module) UploadFile(rw http.ResponseWriter, req *http.Request, par router.Params) error {
	req.ParseMultipartForm(maxUploadMem)
	file, handler, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := m.Files.Upload(handler.Filename, file)
	if err != nil {
		return err
	}
	f.Url = proto.String(m.Files.URLForFile(f))

	return router.ProtoOK(rw, &api.FileResponse{File: f})
}
