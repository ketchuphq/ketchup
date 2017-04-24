package main

import (
	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/admin"
	"github.com/ketchuphq/ketchup/db/bolt"
	"github.com/ketchuphq/ketchup/server/api"
	"github.com/ketchuphq/ketchup/server/backup"
	"github.com/ketchuphq/ketchup/server/content"
	"github.com/ketchuphq/ketchup/server/tls"
)

// set by goreleaser
var version = "dev"

type App struct {
	Content *content.Module
	API     *api.Module
	Admin   *admin.Module
	TLS     *tls.Module
	Backup  *backup.Module

	// configures backend module
	Bolt *bolt.Module
}

func (p *App) Init(c *service.Config) {}

func main() {
	api.KetchupVersion = version
	service.EnvVarName = "KETCHUP_ENV"
	service.Run(&App{})
}
