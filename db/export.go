package db

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/octavore/naga/service"

	"github.com/ketchuphq/ketchup/proto/ketchup/models/export"
)

var marshaler = &jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      "  ",
}

func (m *Module) registerExportCommand(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "db:export <output file>",
		ShortUsage: "export pages and routes to file (as json)",
		Run: func(ctx *service.CommandContext) {
			var path string
			if len(ctx.Args) > 0 {
				path = ctx.Args[0]
			}
			err := m.ExportToJSON(path)
			if err != nil {
				panic(err)
			}
		},
	})
}

func (m *Module) registerImportCommand(c *service.Config) {
	c.AddCommand(&service.Command{
		Keyword:    "db:import <input file>",
		ShortUsage: "import pages and routes from file",
		Run: func(ctx *service.CommandContext) {
			ctx.RequireExactlyNArgs(1)
			err := m.importFromJSON(ctx.Args[0])
			if err != nil {
				panic(err)
			}
		},
	})
}

func (m *Module) ExportToJSON(path string) error {
	var wr io.WriteCloser = os.Stdout
	if path != "" {
		_, err := os.Stat(path)
		if err == nil {
			return fmt.Errorf("file already eists")
		}
		if !os.IsNotExist(err) {
			return err
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		wr = f
	}
	defer wr.Close()
	e, err := m.Export()
	if err != nil {
		return err
	}
	err = marshaler.Marshal(wr, e)
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) importFromJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	e := &export.Export{}
	err = jsonpb.Unmarshal(f, e)
	if err != nil {
		return err
	}
	return m.Import(e)
}

// Export generates a dump of routes and routes
func (m *Module) Export() (*export.Export, error) {
	export := &export.Export{}
	pages, err := m.ListPages(nil)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		page, err := m.GetPage(page.GetUuid())
		if err != nil {
			return nil, err
		}
		export.Pages = append(export.Pages, page)
	}

	routes, err := m.ListRoutes(nil)
	if err != nil {
		return nil, err
	}
	export.Routes = routes
	return export, nil
}

// Import pages and routes from an Export object.
func (m *Module) Import(export *export.Export) error {
	for _, page := range export.GetPages() {
		err := m.UpdatePage(page)
		if err != nil {
			return err
		}
	}

	for _, route := range export.GetRoutes() {
		err := m.UpdateRoute(route)
		if err != nil {
			return err
		}
	}

	return nil
}
