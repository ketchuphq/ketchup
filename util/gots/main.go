// package gots (or go-ts) is responsible for converting structs into
// typescript models. Has limited support for protobuf generated structs,
// specifically `oneof` and `enum` types

package main

import (
	"fmt"
	"os"

	router_api "github.com/octavore/nagax/proto/nagax/router/api"
	"github.com/octavore/pbts"

	"github.com/octavore/ketchup/proto/ketchup/api"
	"github.com/octavore/ketchup/proto/ketchup/models"
	"github.com/octavore/ketchup/proto/ketchup/packages"
)

func main() {
	out := os.Args[1]
	f, err := os.OpenFile(out, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("error opening file:", err)
	}

	g := pbts.NewGenerator(f)
	g.RegisterMany(
		models.Page{},
		models.Content{},
		models.Route{},
		models.Timestamp{},
		models.Theme{},
		models.ThemeTemplate{},
		models.ThemePlaceholder{},
		models.ThemeAsset{},
		models.ContentMultiple{},
		models.ContentText{},
		models.ContentString{},
		models.Author{},
		packages.Package{},
		packages.PackageRelease{},
		packages.Registry{},

		api.TLSSettingsReponse{},
		api.EnableTLSRequest{},
		api.ListPageRequest{},
		api.ListPageRequest_ListPageOptions{},
		api.ListPageResponse{},
		api.ListOptions{},

		router_api.Error{},
		router_api.ErrorResponse{},
	)
	g.Write()
}
