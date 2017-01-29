package main

import (
	"github.com/octavore/naga/service"

	"github.com/octavore/ketchup/admin"
	"github.com/octavore/ketchup/db/bolt"
	"github.com/octavore/ketchup/plugins/pkg"
	"github.com/octavore/ketchup/server/api"
	"github.com/octavore/ketchup/server/content"
	"github.com/octavore/ketchup/server/tls"
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
