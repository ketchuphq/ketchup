package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

const FILES_BUCKET = "files"

// GetFile fetches a file from the DB by key
func (m *Module) GetFile(uuid string) (*models.File, error) {
	file := &models.File{}
	err := m.Get(FILES_BUCKET, uuid, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetFileByName fetches a file from the DB by name
func (m *Module) GetFileByName(name string) (*models.File, error) {
	var file = &models.File{}
	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FILES_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			f := &models.File{}
			err := proto.Unmarshal(v, f)
			if err != nil {
				return err
			}
			if f.GetName() == name {
				file = f
				return nil
			}
		}
		return nil
	})
	return file, err
}

// UpdateFile updates (or creates if necessary) an existing file.
func (m *Module) UpdateFile(file *models.File) error {
	if file.GetUuid() == "" {
		file.Uuid = proto.String(uuid.NewV4().String())
	}
	return m.Update(FILES_BUCKET, file)
}

// DeleteFile deletes a file from the DB (but not from the store).
func (m *Module) DeleteFile(file *models.File) error {
	err := m.delete(FILES_BUCKET, file)
	if err != nil {
		return err
	}
	return nil
}

// ListFiles returns all files recorded in the DB (unsorted)
func (m *Module) ListFiles() ([]*models.File, error) {
	lst := []*models.File{}
	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FILES_BUCKET))
		return b.ForEach(func(_, v []byte) error {
			pb := &models.File{}
			err := proto.Unmarshal(v, pb)
			if err != nil {
				return err
			}
			lst = append(lst, pb)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return lst, nil
}
