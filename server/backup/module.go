package backup

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/server/config"
)

type BackupConfig struct {
	Backup struct {
		Enabled bool `json:"enabled"`
	} `json:"backup"`
}

type Module struct {
	DB     *db.Module
	Logger *logger.Module
	Config *config.Module

	config BackupConfig
}

func (m *Module) Init(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "backup:create",
		ShortUsage: "create a backup",
		Usage:      "create a backup in {dataDir}/backup",
		Run: func(ctx *service.CommandContext) {
			err := m.WriteBackup()
			if err != nil {
				panic(err)
			}
		},
	})

	c.Setup = func() error {
		return m.Config.ReadConfig(&m.config)
	}

	c.Start = func() {
		if m.config.Backup.Enabled {
			go func() {
				err := m.WriteBackup()
				if err != nil {
					m.Logger.Error(err)
				}
			}()
		}
	}
}

var jpb = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	Indent:       "  ",
	OrigName:     false,
}

func (m *Module) WriteBackup() error {
	p := m.Config.DataPath("backups", "")
	err := os.MkdirAll(p, 0700)
	if err != nil {
		return err
	}

	exp, err := m.DB.Export()
	if err != nil {
		return nil
	}

	date := time.Now().Format("20060102-030405")
	backup := path.Join(p, date+".bak")

	buf := &bytes.Buffer{}
	err = jpb.Marshal(buf, exp)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(backup, buf.Bytes(), 0700)
	if err != nil {
		return err
	}
	return nil
}
