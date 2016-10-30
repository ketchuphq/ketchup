package bolt

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
)

// Module bolt provides methods for writing and retrieving
// data from a bolt database.
type Module struct {
	DB *bolt.DB
}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.Setup = func() (err error) {
		if c.Env().IsTest() {
			return
		}
		m.DB, err = bolt.Open("default.db", os.ModePerm, &bolt.Options{
			Timeout: 30 * time.Second,
		})
		if err != nil {
			return err
		}
		return m.init()
	}
	var testDB string
	c.SetupTest = func() {
		suffix := time.Now().Format("20060102150405.999")
		testDB = fmt.Sprintf("test-%s.db", suffix)
		err := os.RemoveAll(testDB)
		if err != nil {
			panic(err)
		}
		m.DB, err = bolt.Open(testDB, os.ModePerm, &bolt.Options{
			Timeout: 30 * time.Second,
		})
		if err != nil {
			panic(err)
		}
		m.init()
		return
	}
	c.Stop = func() {
		err := m.DB.Close()
		if err != nil {
			panic(err)
		}
		if testDB != "" {
			err := os.RemoveAll(testDB)
			if err != nil {
				panic(err)
			}
		}
	}
}

type ErrNoKey string

func (e ErrNoKey) Error() string {
	return fmt.Sprintf("key not found: %s", string(e))
}

func (m *Module) init() error {
	return m.DB.Update(func(tx *bolt.Tx) error {
		buckets := []string{PAGE_BUCKET, ROUTE_BUCKET, USER_BUCKET}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type AddressableProto interface {
	GetUuid() string
	proto.Message
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
		return err
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
