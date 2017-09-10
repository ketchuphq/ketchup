package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

const DATA_BUCKET = "data"

// UpdateData updates (or creates if necessary) an existing data.
// Requires key to be set.
func (m *Module) updateData(tx *bolt.Tx, data *models.Data) error {
	if data.GetKey() == "" {
		return errors.New("bolt: cannot update data without key")
	}
	m.updateTimestampedProto(data)
	key := []byte(data.GetKey())
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		return errors.Wrap(err)
	}
	b := tx.Bucket([]byte(DATA_BUCKET))
	return errors.Wrap(b.Put(key, dataBytes))
}

func (m *Module) UpdateData(data *models.Data) error {
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		return m.updateData(tx, data)
	})
}

func (m *Module) UpdateDataBatch(data []*models.Data) error {
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		for _, d := range data {
			err := m.updateData(tx, d)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetData fetches a data from the DB by key
func (m *Module) GetData(key string) (*models.Data, error) {
	data := &models.Data{}
	err := m.Get(DATA_BUCKET, key, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteData deletes a data from the DB (but not any related pages).
func (m *Module) DeleteData(data *models.Data) error {
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DATA_BUCKET))
		return b.Delete([]byte(data.GetKey()))
	})
}

// ListData returns all data stored in the DB (unsorted)
func (m *Module) ListData() ([]*models.Data, error) {
	lst := []*models.Data{}
	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DATA_BUCKET))
		return b.ForEach(func(_, v []byte) error {
			pb := &models.Data{}
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
