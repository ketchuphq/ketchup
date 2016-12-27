package bolt

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/octavore/press/db"
	"github.com/octavore/press/util/errors"
)

// Module bolt provides methods for writing and retrieving
// data from a bolt database.
type Module struct {
	Bolt    *bolt.DB
	Backend *db.Module
	Logger  *logger.Module
}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.Setup = func() (err error) {
		m.Backend.Backend = m
		if c.Env().IsTest() {
			return
		}
		m.Bolt, err = bolt.Open("default.db", os.ModePerm, &bolt.Options{
			Timeout: 5 * time.Second,
		})
		if err != nil {
			if err == bolt.ErrTimeout {
				m.Logger.Error("bolt: timeout while connecting; it may be that the database is already in use by another process.")
			}
			return errors.Wrap(err)
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
		m.Bolt, err = bolt.Open(testDB, os.ModePerm, &bolt.Options{
			Timeout: 30 * time.Second,
		})
		if err != nil {
			panic(err)
		}
		m.init()
		return
	}
	c.Stop = func() {
		err := m.Bolt.Close()
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
	return fmt.Sprintf("boltdb: key not found: %s", string(e))
}

func (m *Module) init() error {
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		buckets := []string{PAGE_BUCKET, ROUTE_BUCKET, USER_BUCKET}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return errors.Wrap(err)
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
	return m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get([]byte(key))
		if data == nil {
			return ErrNoKey(bucket + ":" + key)
		}
		return proto.Unmarshal(data, pb)
	})
}

func (m *Module) Update(bucket string, pb AddressableProto) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return errors.Wrap(err)
	}
	key := []byte(pb.GetUuid())
	if len(key) == 0 {
		return fmt.Errorf("no uuid for proto")
	}
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(key, data)
	})
}

func (m *Module) delete(bucket string, pb AddressableProto) error {
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(pb.GetUuid()))
	})
}

func (m *Module) BackupToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err)
	}
	defer f.Close()
	return m.Backup(f)
}

func (m *Module) Backup(w io.Writer) error {
	return m.Bolt.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return errors.Wrap(err)
	})
}

func (m *Module) Debug(w io.Writer) error {
	return m.Bolt.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			fmt.Fprintln(w, string(name)+":")
			return bucket.ForEach(func(key, value []byte) error {
				fmt.Fprintln(w, string(key))
				fmt.Fprintln(w, string(value))
				return nil
			})
		})
	})
}
