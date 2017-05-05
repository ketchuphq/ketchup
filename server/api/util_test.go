package api

import (
	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/db"
	"github.com/ketchuphq/ketchup/db/dummy"
	"github.com/ketchuphq/ketchup/util/testutil/memlogger"
)

type TestModule struct {
	*Module
	DB *db.Module
}

func (m *TestModule) Init(c *service.Config) {
	c.Setup = func() error {
		m.Logger.Logger = &memlogger.MemoryLogger{}
		return m.DB.Register(dummy.New())
	}
}

type testEnv struct {
	module *Module
	db     *dummy.DummyDB
	logger *memlogger.MemoryLogger
	stop   func()
}

func setup() testEnv {
	module := &TestModule{}
	stop := service.New(module).StartForTest()
	return testEnv{
		module: module.Module,
		db:     module.DB.Backend.(*dummy.DummyDB),
		logger: module.Logger.Logger.(*memlogger.MemoryLogger),
		stop:   stop,
	}
}
