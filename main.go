package main

import (
	"github.com/octavore/naga/service"

	"github.com/octavore/press/admin"
	"github.com/octavore/press/db/bolt"
	"github.com/octavore/press/plugins/pkg"
	"github.com/octavore/press/server/api"
	"github.com/octavore/press/server/content"
	"github.com/octavore/press/server/tls"
)

type App struct {
	Content *content.Module
	API     *api.Module
	Admin   *admin.Module
	TLS     *tls.Module
	Package *pkg.Module

	// configures backend module
	Bolt *bolt.Module
}

func (p *App) Init(c *service.Config) {}

func main() {
	service.Run(&App{})
}
