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
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/config"
)

type BoltConfig struct {
	Bolt struct {
		Path string `json:"path"`
	} `json:"bolt"`
}

// Module bolt provides methods for writing and retrieving
// data from a bolt database.
type Module struct {
	Bolt         *bolt.DB
	Backend      *db.Module
	ConfigModule *config.Module
	Logger       *logger.Module

	config BoltConfig
}

// Init implements service.Init
func (m *Module) Init(c *service.Config) {
	c.Setup = func() (err error) {
		m.Backend.Register(m)
		err = m.ConfigModule.ReadConfig(&m.config)
		if err != nil {
			return err
		}
		if c.Env().IsTest() {
			return
		}
		m.config.Bolt.Path = m.ConfigModule.DataPath(m.config.Bolt.Path, "default.db")
		return m.connectBolt(m.config.Bolt.Path)
	}
	var testDB string
	c.SetupTest = func() {
		suffix := time.Now().Format("20060102150405.999")
		testDB = fmt.Sprintf("test-%s.db", suffix)
		err := os.RemoveAll(testDB)
		if err != nil {
			panic(err)
		}
		err = m.connectBolt(testDB)
		if err != nil {
			panic(err)
		}
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

// ErrNoKey is a helper to create key not found errors.
type ErrNoKey string

func (e ErrNoKey) Error() string {
	return fmt.Sprintf("boltdb: key not found: %s", string(e))
}

func (m *Module) connectBolt(path string) error {
	var err error
	m.Bolt, err = bolt.Open(path, os.ModePerm, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		if err == bolt.ErrTimeout {
			m.Logger.Error("bolt: timeout while connecting; it may be that the database is already in use by another process.")
		}
		return errors.Wrap(err)
	}
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		buckets := []string{PAGE_BUCKET, ROUTE_BUCKET, USER_BUCKET, DATA_BUCKET, FILES_BUCKET}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	})
}

// Get is a generic method of getting a proto from the bolt database.
func (m *Module) Get(bucket, key string, pb proto.Message) error {
	return m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get([]byte(key))
		if data == nil {
			return errors.Wrap(ErrNoKey(bucket + ":" + key))
		}
		return proto.Unmarshal(data, pb)
	})
}

func (m *Module) updateTimestampedProto(tsp db.TimestampedProto) {
	ts := tsp.GetTimestamps()
	if ts == nil {
		ts = &models.Timestamp{}
	}
	// convert time to millis
	nowMillis := time.Now().UnixNano() / 1e6
	if ts.GetCreatedAt() == 0 {
		ts.CreatedAt = proto.Int64(nowMillis)
	}
	ts.UpdatedAt = proto.Int64(nowMillis)
	tsp.SetTimestamps(ts)
}

// Update is a generic way of updating AddressableProto data in the bolt database.
func (m *Module) Update(bucket string, pb db.AddressableProto) error {
	if tsp, ok := pb.(db.TimestampedProto); ok {
		m.updateTimestampedProto(tsp)
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		return errors.Wrap(err)
	}
	key := []byte(pb.GetUuid())
	if len(key) == 0 {
		return errors.New("no uuid for proto")
	}
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return errors.Wrap(b.Put(key, data))
	})
}

func (m *Module) delete(bucket string, pb db.AddressableProto) error {
	return m.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(pb.GetUuid()))
	})
}

// BackupToFile writes the entire database to the file at path.
func (m *Module) BackupToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err)
	}
	defer f.Close()
	return m.Backup(f)
}

// Backup writes the entire database to w.
func (m *Module) Backup(w io.Writer) error {
	return m.Bolt.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return errors.Wrap(err)
	})
}

// Debug prints out all data in the database. Does not deserialize saved protos.
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
