package db

import (
	"fmt"

	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/util/errors"
)

type Module struct {
	Backend
}

func (m *Module) Init(c *service.Config) {
	m.registerExportCommand(c)
	m.registerImportCommand(c)

	c.Start = func() {
		if m.Backend == nil {
			panic("backend not configured")
		}
	}
}

func (m *Module) Register(b Backend) error {
	if m.Backend != nil {
		return errors.Wrap(fmt.Errorf("backend already configured"))
	}
	m.Backend = b
	return nil
}
