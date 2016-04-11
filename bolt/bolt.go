package bolt

import (
	"fmt"
	"os"
	"time"
	"io"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/octavore/press/proto/press/models"
)

type Module struct {
	DB *bolt.DB
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() (err error) {
		m.DB, err = bolt.Open("default.db", os.ModePerm, &bolt.Options{
			Timeout: 30 * time.Second,
		})
		return
	}
}

type ErrNoKey string

func (e ErrNoKey) Error() string {
	return fmt.Sprintf("key not found: %s", e)
}

const (
	PAGE_BUCKET  = "pages"
	ROUTE_BUCKET = "routes"
)

func (m *Module) GetPage(uuid string) (*models.Page, error) {
	page := &models.Page{}
	err := m.Get(PAGE_BUCKET, uuid, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (m *Module) UpdatePage(page *models.Page) error {
	return m.Update(PAGE_BUCKET, page)
}

type AddressableProto interface {
	GetUuid() string
	proto.Message
}

func (m *Module) GetRoute(key string) (*models.Route, error) {
	route := &models.Route{}
	err := m.Get(ROUTE_BUCKET, key, route)
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (m *Module) Get(bucket, key string, pb proto.Message) error {
	return m.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get([]byte(key))
		if data == nil {
			return ErrNoKey(key)
		}
		return proto.Unmarshal(data, pb)
	})
}

func (m *Module) Update(bucket string, pb AddressableProto) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil
	}
	key := []byte(pb.GetUuid())
	if len(key) == 0 {
		return fmt.Errorf("no uuid for proto")
	}
	return m.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(key, data)
	})
}

func (m *Module) BackupToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return m.Backup(f)
}

func (m *Module) Backup(w io.Writer) error {
	return m.DB.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return err
	})
}
