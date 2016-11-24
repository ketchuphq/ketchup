package bolt

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"

	"github.com/octavore/press/proto/press/models"
)

const USER_BUCKET = "users"

func (m *Module) GetUser(uuid string) (*models.User, error) {
	user := &models.User{}
	err := m.Get(USER_BUCKET, uuid, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *Module) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, nil
	}
	var user *models.User
	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USER_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			u := &models.User{}
			err := proto.Unmarshal(v, u)
			if err != nil {
				return err
			}
			if u.GetEmail() == email {
				user = u
				return nil
			}
		}
		return nil
	})
	return user, err
}

func (m *Module) UpdateUser(user *models.User) error {
	if user.Uuid == nil {
		user.Uuid = proto.String(uuid.NewV4().String())
	}
	if user.GetEmail() == "" {
		return errors.New("user email is required")
	}
	return m.Update(USER_BUCKET, user)
}
