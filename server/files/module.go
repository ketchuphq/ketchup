package files

import (
	"fmt"
	"io"

	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/util/errors"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/server/config"
	"github.com/ketchuphq/ketchup/server/files/disk"
	"github.com/ketchuphq/ketchup/server/router"
)

const FileURLPrefix = "/static"

type Store interface {
	Upload(uuid string, r io.Reader) error
	Delete(uuid string) error
	Get(uuid string) (io.ReadCloser, error)
}

type Module struct {
	Config *config.Module
	DB     *db.Module
	Router *router.Module
	Logger *logger.Module

	store Store
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		u, err := disk.NewDiskStore(m.Config.DataPath("files", ""))
		if err != nil {
			return errors.Wrap(err)
		}
		m.store = u
		return nil
	}
}

func (m *Module) URLForFile(file *models.File) string {
	return fmt.Sprintf("%s/%s", FileURLPrefix, file.GetName())
}
