package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/nagax/util/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
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

func (m *Module) GetUserByToken(token string) (*models.User, error) {
	if token == "" {
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
			if u.GetToken() == token {
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
	u2, err := m.GetUserByEmail(user.GetEmail())
	if err != nil {
		return err
	}
	if u2 != nil && u2.GetUuid() != user.GetUuid() {
		return errors.New("user already exists")
	}

	return m.Update(USER_BUCKET, user)
}

func (m *Module) ListUsers() ([]*models.User, error) {
	lst := []*models.User{}
	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USER_BUCKET))
		return b.ForEach(func(_, v []byte) error {
			pb := &models.User{}
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
