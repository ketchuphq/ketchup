package templates

import (
	"log"

	"github.com/octavore/naga/service"
)

func (m *Module) registerCommands(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "themes:get <package name> <git url>",
		ShortUsage: "get the theme",
		Usage:      "get the theme",
		Run: func(ctx *service.CommandContext) {
			packageName := ctx.Args[0]
			packageURL := ctx.Args[1]
			p := m.ConfigModule.DataPath(m.config.Themes.Path, themeDir)
			log.Printf("Downloading %s to %s from %s...", packageName, p, packageURL)
			err := m.Pkg.Clone(packageName, packageURL, p)
			if err != nil {
				panic(err)
			}
		},
	})

	c.AddCommand(&service.Command{
		Keyword:    "themes:install <package name> <path to package tgz>",
		ShortUsage: "get the theme",
		Usage:      "get the theme",
		Run: func(ctx *service.CommandContext) {
			packageName := ctx.Args[0]
			packageURL := ctx.Args[1]
			p := m.ConfigModule.DataPath(m.config.Themes.Path, themeDir)
			log.Printf("Downloading %s to %s from %s...", packageName, p, packageURL)
			err := m.Pkg.Clone(packageName, packageURL, p)
			if err != nil {
				panic(err)
			}
		},
	})
}
