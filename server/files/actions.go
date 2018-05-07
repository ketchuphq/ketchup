package files

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/util/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/willnorris/imageproxy"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func (m *Module) Upload(filename string, wr io.Reader) (*models.File, error) {
	// files are assigned a random id when stored to discourage manual renaming of files
	filename = strings.TrimPrefix(filename, "/")
	file, err := m.DB.GetFileByName(filename)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if file == nil || file.GetUuid() == "" {
		file = &models.File{
			Uuid: proto.String(uuid.NewV4().String()),
			Name: proto.String(filename),
		}
	}

	err = m.store.Upload(file.GetUuid(), wr)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = m.DB.UpdateFile(file)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return file, nil
}

// Get returns a reader, and nil, nil if file is not found
func (m *Module) Get(filename string) (io.ReadCloser, error) {
	filename = strings.TrimPrefix(filename, "/")
	file, err := m.DB.GetFileByName(filename)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if file == nil {
		return nil, nil
	}
	return m.store.Get(file.GetUuid())
}

// Get returns a reader, and nil, nil if file is not found
func (m *Module) Delete(uuid string) error {
	file, err := m.DB.GetFile(uuid)
	if err != nil {
		return errors.Wrap(err)
	}
	if file == nil {
		return nil
	}
	err = m.DB.DeleteFile(file)
	if err != nil {
		return errors.Wrap(err)
	}
	return m.store.Delete(file.GetUuid())
}

// GetWithTransform attempts to transform the image
func (m *Module) GetWithTransform(filename string, optStr string) (io.ReadCloser, error) {
	filename = strings.TrimPrefix(filename, "/")
	r, err := m.Get(filename)
	if r == nil || err != nil || optStr == "" {
		return r, err
	}

	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	opts := imageproxy.ParseOptions(optStr)
	output, err := imageproxy.Transform(data, opts)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return ioutil.NopCloser(bytes.NewBuffer(output)), nil
}
